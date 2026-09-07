package services

import (
	"testing"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"gorm.io/gorm"
)

type paymentOrderRepository struct {
	repositories.OrderRepository
	current models.Order
	updated bool
}

func (r *paymentOrderRepository) FindById(uint) (models.Order, error) { return r.current, nil }
func (r *paymentOrderRepository) Update(order *models.Order) error {
	r.current = *order
	r.updated = true
	return nil
}

type paymentStoreConfig struct {
	StoreConfigService
	config models.StoreConfig
}

func (s paymentStoreConfig) GetStoreConfig() (models.StoreConfig, error) { return s.config, nil }

func TestCashPaymentUsesStoredSubtotalAndDeliveryChoice(t *testing.T) {
	for _, delivery := range []bool{false, true} {
		t.Run(map[bool]string{false: "pickup", true: "delivery"}[delivery], func(t *testing.T) {
			repo := &paymentOrderRepository{current: models.Order{Model: gorm.Model{ID: 1}, Status: models.ActiveStatus, ShoppingCart: models.ShoppingCart{Total: 10}}}
			store := paymentStoreConfig{config: models.StoreConfig{DeliveryIsActive: true, DeliveryPrice: 2.50, PaymentCashIsActive: true}}
			order := models.Order{Model: gorm.Model{ID: 1}, PaymentMethod: models.CashPaymentMethod, IsDelivery: delivery, Total: 0.01, ShoppingCart: models.ShoppingCart{Total: 0.01}}
			if err := NewOrderService(repo, store).Payment(&order); err != nil {
				t.Fatal(err)
			}
			want := 10.0
			if delivery {
				want += 2.50
			}
			if !repo.updated || order.Total != want || order.Status != models.ProcessingStatus || order.IsDelivery != delivery {
				t.Fatalf("incorrect payment: %+v", order)
			}
		})
	}
}

func TestPaymentRejectsDisabledMethodsAndClosedOrders(t *testing.T) {
	for _, status := range []models.ShoppingCartStatus{models.ActiveStatus, models.CancelledStatus, models.ProcessingStatus} {
		repo := &paymentOrderRepository{current: models.Order{Model: gorm.Model{ID: 1}, Status: status, ShoppingCart: models.ShoppingCart{Total: 10}}}
		order := models.Order{Model: gorm.Model{ID: 1}, PaymentMethod: models.CashPaymentMethod}
		if err := NewOrderService(repo, paymentStoreConfig{}).Payment(&order); err == nil || repo.updated {
			t.Fatalf("invalid payment accepted for %s", status)
		}
	}
}
