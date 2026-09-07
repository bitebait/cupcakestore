package repositories

import (
	"errors"
	"testing"
	"time"

	"github.com/bitebait/cupcakestore/models"
)

func TestCheckoutDoesNotReserveStockWithoutFulfillmentOrPayment(t *testing.T) {
	for _, tc := range []struct {
		name   string
		fields map[string]any
		want   error
	}{
		{"no fulfillment", map[string]any{"delivery_is_active": false, "physical_store_address": " "}, models.ErrFulfillmentUnavailable},
		{"no payment", map[string]any{"payment_cash_is_active": false, "payment_pix_is_active": false}, models.ErrPaymentUnavailable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db, profile, product := commerceDB(t)
			cart := cartWithItem(t, db, profile.ID, product.ID, 2)
			if err := db.Model(&models.StoreConfig{}).Where("id > 0").UpdateColumns(tc.fields).Error; err != nil {
				t.Fatal(err)
			}
			if _, err := NewOrderRepository(db).FindOrCreate(profile.ID, cart.ID); !errors.Is(err, tc.want) {
				t.Fatalf("err=%v want=%v", err, tc.want)
			}
			assertStock(t, db, product.ID, 10)
			if err := db.First(&cart, cart.ID).Error; err != nil {
				t.Fatal(err)
			}
			if cart.OrderID != 0 {
				t.Fatal("failed checkout froze cart")
			}
		})
	}
}

func TestPaymentRejectsChangedStoreConfigurationBeforeSaving(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 1)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	version := paymentConfigVersion(t, db)
	if err := db.Model(&models.StoreConfig{}).Where("id > 0").UpdateColumns(map[string]any{"physical_store_address": "Rua diferente", "updated_at": version.Add(time.Second)}).Error; err != nil {
		t.Fatal(err)
	}
	order.PaymentMethod, order.Status = models.CashPaymentMethod, models.ProcessingStatus
	if err := repo.UpdatePayment(&order, version); err == nil {
		t.Fatal("payment accepted stale checkout terms")
	}
	stored, err := repo.FindById(order.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Status != models.ActiveStatus || stored.DeliveryDetail.StoreAddress == "Rua diferente" {
		t.Fatal("failed payment altered recorded terms")
	}
}

func TestPaymentCannotRewriteAlreadyIssuedPixInsideTransaction(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 1)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&order).UpdateColumns(map[string]any{"status": models.AwaitingPaymentStatus, "pix_string": "original"}).Error; err != nil {
		t.Fatal(err)
	}
	order.PaymentMethod, order.Status, order.PixString = models.CashPaymentMethod, models.ProcessingStatus, ""
	if err := repo.UpdatePayment(&order, paymentConfigVersion(t, db)); err == nil {
		t.Fatal("replaced issued Pix with another payment")
	}
}

func TestPickupDoesNotRequireCustomerAddressAndCannotLoseStoreAddress(t *testing.T) {
	db, profile, product := commerceDB(t)
	if err := db.Model(&profile).UpdateColumns(map[string]any{"address": "", "city": "", "state": "", "postal_code": ""}).Error; err != nil {
		t.Fatal(err)
	}
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	if order.IsDelivery || order.Total != order.ShoppingCart.Total {
		t.Fatalf("pickup default incorrect: %+v", order)
	}
	if err := db.Model(&models.StoreConfig{}).Where("id > 0").UpdateColumn("physical_store_address", "").Error; err != nil {
		t.Fatal(err)
	}
	order.PaymentMethod, order.Status = models.CashPaymentMethod, models.ProcessingStatus
	if err := repo.UpdatePayment(&order, paymentConfigVersion(t, db)); !errors.Is(err, models.ErrPickupUnavailable) {
		t.Fatalf("accepted pickup without address: %v", err)
	}
	stored, err := repo.FindById(order.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Status != models.ActiveStatus {
		t.Fatal("failed payment persisted")
	}
}

func TestPaymentValidatesAndSnapshotsTheCorrectedAddress(t *testing.T) {
	db, profile, product := commerceDB(t)
	if err := db.Model(&profile).UpdateColumns(map[string]any{"address": "", "city": "", "state": "", "postal_code": ""}).Error; err != nil {
		t.Fatal(err)
	}
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	order.IsDelivery, order.PaymentMethod, order.Status, order.Total = true, models.CashPaymentMethod, models.ProcessingStatus, order.ShoppingCart.Total+2.50
	if err := repo.UpdatePayment(&order, paymentConfigVersion(t, db)); !errors.Is(err, models.ErrDeliveryAddressIncomplete) {
		t.Fatalf("accepted incomplete delivery address: %v", err)
	}
	if err := db.Model(&profile).UpdateColumns(map[string]any{"address": "Rua corrigida, 42", "city": "São Paulo", "state": "sp", "postal_code": "01001-000"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := repo.UpdatePayment(&order, paymentConfigVersion(t, db)); err != nil {
		t.Fatal(err)
	}
	if order.DeliveryDetail.UserAddress != "Rua corrigida, 42" || order.DeliveryDetail.UserState != "SP" {
		t.Fatalf("payment retained draft address: %+v", order.DeliveryDetail)
	}
	if err := db.Model(&profile).UpdateColumn("address", "Outro endereço").Error; err != nil {
		t.Fatal(err)
	}
	order.Status = models.DeliveredStatusAwaiting
	if err := repo.Update(&order); err != nil {
		t.Fatal(err)
	}
	if order.DeliveryDetail.UserAddress != "Rua corrigida, 42" {
		t.Fatal("status update rewrote recorded delivery address")
	}
}
