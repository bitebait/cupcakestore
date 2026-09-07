package repositories

import (
	"testing"

	"github.com/bitebait/cupcakestore/models"
)

func TestOrderCompletionRespectsDeliveryChoice(t *testing.T) {
	for _, test := range []struct {
		name       string
		isDelivery bool
		status     models.ShoppingCartStatus
	}{
		{"pickup from preparing", false, models.ProcessingStatus},
		{"pickup from ready", false, models.DeliveredStatusAwaiting},
		{"delivery from preparing", true, models.ProcessingStatus},
		{"delivery from ready", true, models.DeliveredStatusAwaiting},
	} {
		t.Run(test.name, func(t *testing.T) {
			db, profile, product := commerceDB(t)
			cart := cartWithItem(t, db, profile.ID, product.ID, 1)
			repo := NewOrderRepository(db)
			order, err := repo.FindOrCreate(profile.ID, cart.ID)
			if err != nil {
				t.Fatal(err)
			}
			order.IsDelivery = test.isDelivery
			order.PaymentMethod = models.CashPaymentMethod
			order.Status = models.ProcessingStatus
			if !test.isDelivery {
				order.DeliveryPrice = 0
				order.Total = order.ShoppingCart.Total
			}
			if err := repo.Update(&order); err != nil {
				t.Fatal(err)
			}
			if test.status == models.DeliveredStatusAwaiting {
				order.Status = test.status
				if err := repo.Update(&order); err != nil {
					t.Fatal(err)
				}
			}
			order.Status = models.DeliveredStatusDelivered
			err = repo.Update(&order)
			if test.isDelivery && err == nil {
				t.Fatal("delivery completed without shipment")
			}
			if !test.isDelivery && err != nil {
				t.Fatalf("pickup completion rejected: %v", err)
			}
			persisted, err := repo.FindById(order.ID)
			if err != nil {
				t.Fatal(err)
			}
			wantStatus := models.DeliveredStatusDelivered
			if test.isDelivery {
				wantStatus = test.status
			}
			if persisted.Status != wantStatus || persisted.IsDelivery != test.isDelivery {
				t.Fatalf("unexpected completion state: status=%s delivery=%v", persisted.Status, persisted.IsDelivery)
			}
			assertStock(t, db, product.ID, 9)
		})
	}
}
