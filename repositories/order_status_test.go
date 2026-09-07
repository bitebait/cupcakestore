package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"testing"
)

func TestStatusUpdatePreservesDeliveryTermsAfterStoreChanges(t *testing.T) {
	for _, change := range []struct {
		name   string
		values map[string]any
	}{
		{"new fee", map[string]any{"delivery_price": 9.0}},
		{"delivery disabled", map[string]any{"delivery_is_active": false}},
	} {
		t.Run(change.name, func(t *testing.T) {
			db, profile, product := commerceDB(t)
			cart := cartWithItem(t, db, profile.ID, product.ID, 1)
			repo := NewOrderRepository(db)
			order, err := repo.FindOrCreate(profile.ID, cart.ID)
			if err != nil {
				t.Fatal(err)
			}
			// Existing Pix charge created before the store settings changed.
			if err := db.Model(&models.Order{}).Where("id = ?", order.ID).Updates(map[string]any{"status": models.PaymentApprovedStatus, "pix_url": "/charge/original", "pix_string": "original-code"}).Error; err != nil {
				t.Fatal(err)
			}
			order, err = repo.FindById(order.ID)
			if err != nil {
				t.Fatal(err)
			}
			if err := db.Model(&models.StoreConfig{}).Where("id > 0").UpdateColumns(change.values).Error; err != nil {
				t.Fatal(err)
			}
			order.Status = models.ProcessingStatus
			if err := repo.Update(&order); err != nil {
				t.Fatalf("status-only update failed after store changed: %v", err)
			}
			if order.Total != 7.75 || order.DeliveryPrice != 2.5 || !order.IsDelivery || order.PixURL != "/charge/original" || order.PixString != "original-code" {
				t.Fatalf("recorded payment or delivery terms changed: %+v", order)
			}
			if order.Status != models.ProcessingStatus {
				t.Fatalf("status = %v", order.Status)
			}
			assertStock(t, db, product.ID, 9)
		})
	}
}

func TestStatusUpdateCannotRewritePaymentOrDelivery(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 1)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&models.Order{}).Where("id = ?", order.ID).Update("status", models.PaymentApprovedStatus).Error; err != nil {
		t.Fatal(err)
	}
	order, err = repo.FindById(order.ID)
	if err != nil {
		t.Fatal(err)
	}
	order.Total = 0.01
	order.DeliveryPrice = 0
	order.IsDelivery = false
	order.PaymentMethod = models.CashPaymentMethod
	order.PixURL = "/untrusted"
	order.PixString = "overwritten"
	order.Status = models.ProcessingStatus
	if err := repo.Update(&order); err != nil {
		t.Fatal(err)
	}
	if order.Total != 7.75 || order.DeliveryPrice != 2.5 || !order.IsDelivery || order.PaymentMethod != models.PixPaymentMethod || order.PixURL != "" || order.PixString != "" {
		t.Fatalf("status operation rewrote financial terms: %+v", order)
	}
}

func TestAdministrativeUpdateCannotSkipCheckoutOrPaymentConfirmation(t *testing.T) {
	for _, initial := range []models.ShoppingCartStatus{models.ActiveStatus, models.AwaitingPaymentStatus} {
		t.Run(string(initial), func(t *testing.T) {
			db, profile, product := commerceDB(t)
			cart := cartWithItem(t, db, profile.ID, product.ID, 1)
			repo := NewOrderRepository(db)
			order, err := repo.FindOrCreate(profile.ID, cart.ID)
			if err != nil {
				t.Fatal(err)
			}
			if err := db.Model(&models.Order{}).Where("id = ?", order.ID).Update("status", initial).Error; err != nil {
				t.Fatal(err)
			}
			for _, target := range []models.ShoppingCartStatus{models.PaymentApprovedStatus, models.ProcessingStatus} {
				order, err = repo.FindById(order.ID)
				if err != nil {
					t.Fatal(err)
				}
				order.Status = target
				if err := repo.Update(&order); err == nil {
					t.Fatalf("admin skipped checkout/payment: %s -> %s", initial, target)
				}
			}
		})
	}
}

func TestPaymentUpdateStillValidatesCurrentDeliveryTerms(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 1)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&models.StoreConfig{}).Where("id > 0").UpdateColumn("delivery_price", 4).Error; err != nil {
		t.Fatal(err)
	}
	order.Status = models.ProcessingStatus
	order.PaymentMethod = models.CashPaymentMethod
	if err := repo.UpdatePayment(&order, paymentConfigVersion(t, db)); err == nil {
		t.Fatal("accepted outdated fee for a new payment choice")
	}
	order.Total = 9.25
	if err := repo.UpdatePayment(&order, paymentConfigVersion(t, db)); err != nil {
		t.Fatal(err)
	}
	if order.Total != 9.25 || order.DeliveryPrice != 4 {
		t.Fatalf("payment choice ignored current delivery fee: %+v", order)
	}
	order.IsDelivery = false
	order.Total = 5.25
	if err := repo.UpdatePayment(&order, paymentConfigVersion(t, db)); err == nil {
		t.Fatal("rewrote financial terms after processing began")
	}
}
