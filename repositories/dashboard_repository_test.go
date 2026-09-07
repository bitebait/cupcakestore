package repositories

import (
	"testing"

	"github.com/bitebait/cupcakestore/models"
)

func TestDashboardCountsShippedOrdersAndOnlyCompletedFulfillment(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 1)
	order, err := NewOrderRepository(db).FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, status := range []models.ShoppingCartStatus{models.PaymentApprovedStatus, models.ProcessingStatus, models.DeliveredStatusAwaiting, models.DeliveredStatusSent, models.DeliveredStatusDelivered, models.CancelledStatus} {
		if err := db.Model(&order).UpdateColumn("status", status).Error; err != nil {
			t.Fatal(err)
		}
		info := NewDashboardRepository(db).GetInfo(30)
		wantOrders, wantCompleted := int64(1), int64(0)
		if status == models.CancelledStatus {
			wantOrders = 0
		}
		if status == models.DeliveredStatusDelivered {
			wantCompleted = 1
		}
		if info.NewOrders != wantOrders || info.CompletedOrders != wantCompleted {
			t.Fatalf("status %s: %+v", status, info)
		}
	}
}
