package repositories

import (
	"sync"
	"testing"

	"github.com/bitebait/cupcakestore/models"
)

func TestPaymentConfirmationRequiresStaffAndPreservesFirstConfirmation(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	// BR Code validation belongs to the service; this test isolates the
	// persistence contract and the race with customer cancellation.
	if err := db.Model(&order).UpdateColumns(map[string]any{"status": models.AwaitingPaymentStatus, "pix_string": "recorded-payload"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := repo.ConfirmPayment(order.ID, profile.ID); err == nil {
		t.Fatal("customer confirmed own payment")
	}
	if err := db.Model(&models.User{}).Where("id = ?", profile.UserID).UpdateColumn("is_staff", true).Error; err != nil {
		t.Fatal(err)
	}
	if err := repo.ConfirmPayment(order.ID, profile.ID); err != nil {
		t.Fatal(err)
	}
	confirmed, err := repo.FindById(order.ID)
	if err != nil {
		t.Fatal(err)
	}
	if confirmed.Status != models.PaymentApprovedStatus || confirmed.PaymentConfirmedAt == nil || confirmed.PaymentConfirmedByID == nil || *confirmed.PaymentConfirmedByID != profile.ID {
		t.Fatalf("missing confirmation record: %+v", confirmed)
	}
	if err := repo.ConfirmPayment(order.ID, profile.ID); err != nil {
		t.Fatal(err)
	}
	again, err := repo.FindById(order.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !again.PaymentConfirmedAt.Equal(*confirmed.PaymentConfirmedAt) || !again.UpdatedAt.Equal(confirmed.UpdatedAt) {
		t.Fatal("repeat confirmation rewrote original evidence")
	}
	if err := repo.CancelForCustomer(order.ID, profile.ID); err == nil {
		t.Fatal("customer cancelled confirmed payment")
	}
	assertStock(t, db, product.ID, 8)
	if err := repo.Cancel(order.ID); err != nil {
		t.Fatal(err)
	}
	assertStock(t, db, product.ID, 10)
}

func TestPaymentConfirmationAndCustomerCancellationHaveOneWinner(t *testing.T) {
	db, profile, product := commerceDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	// Match the single-connection SQLite configuration used by the app.
	sqlDB.SetMaxOpenConns(1)
	if err := db.Model(&models.User{}).Where("id = ?", profile.UserID).UpdateColumn("is_staff", true).Error; err != nil {
		t.Fatal(err)
	}
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&order).UpdateColumns(map[string]any{"status": models.AwaitingPaymentStatus, "pix_string": "recorded-payload"}).Error; err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	var confirmationErr, cancellationErr error
	go func() { defer wg.Done(); <-start; confirmationErr = repo.ConfirmPayment(order.ID, profile.ID) }()
	go func() { defer wg.Done(); <-start; cancellationErr = repo.CancelForCustomer(order.ID, profile.ID) }()
	close(start)
	wg.Wait()
	if (confirmationErr == nil) == (cancellationErr == nil) {
		t.Fatalf("expected one successful action: confirm=%v cancel=%v", confirmationErr, cancellationErr)
	}
	stored, err := repo.FindById(order.ID)
	if err != nil {
		t.Fatal(err)
	}
	if confirmationErr == nil {
		if stored.Status != models.PaymentApprovedStatus || stored.PaymentConfirmedAt == nil {
			t.Fatal("customer cancelled after confirmation")
		}
		assertStock(t, db, product.ID, 8)
	} else {
		if stored.Status != models.CancelledStatus || stored.PaymentConfirmedAt != nil {
			t.Fatal("confirmed a cancelled order")
		}
		assertStock(t, db, product.ID, 10)
	}
}

func TestCancelledOrCashOrdersCannotBeConfirmedAsPix(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  models.ShoppingCartStatus
		method  models.PaymentMethod
		payload string
	}{
		{"draft", models.ActiveStatus, models.PixPaymentMethod, "payload"},
		{"cash", models.ProcessingStatus, models.CashPaymentMethod, ""},
		{"cancelled", models.CancelledStatus, models.PixPaymentMethod, "payload"},
		{"missing Pix", models.AwaitingPaymentStatus, models.PixPaymentMethod, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db, profile, product := commerceDB(t)
			if err := db.Model(&models.User{}).Where("id = ?", profile.UserID).UpdateColumn("is_staff", true).Error; err != nil {
				t.Fatal(err)
			}
			cart := cartWithItem(t, db, profile.ID, product.ID, 1)
			repo := NewOrderRepository(db)
			order, err := repo.FindOrCreate(profile.ID, cart.ID)
			if err != nil {
				t.Fatal(err)
			}
			if err := db.Model(&order).UpdateColumns(map[string]any{"status": tc.status, "payment_method": tc.method, "pix_string": tc.payload}).Error; err != nil {
				t.Fatal(err)
			}
			if err := repo.ConfirmPayment(order.ID, profile.ID); err == nil {
				t.Fatal("invalid confirmation accepted")
			}
		})
	}
}

func TestCustomerCancellationChecksOwnershipInsideTransaction(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.CancelForCustomer(order.ID, profile.ID+1); err == nil {
		t.Fatal("cancelled another customer's order")
	}
	assertStock(t, db, product.ID, 8)
	if err := repo.CancelForCustomer(order.ID, profile.ID); err != nil {
		t.Fatal(err)
	}
	if err := repo.CancelForCustomer(order.ID, profile.ID); err != nil {
		t.Fatal(err)
	}
	assertStock(t, db, product.ID, 10)
}
