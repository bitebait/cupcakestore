package services

import (
	"errors"
	"testing"
	"time"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"gorm.io/gorm"
)

type paymentOrderRepository struct {
	repositories.OrderRepository
	current     models.Order
	updated     bool
	confirmedBy uint
}

func (r *paymentOrderRepository) ConfirmPayment(_ uint, staffID uint) error {
	r.confirmedBy = staffID
	return nil
}

func (r *paymentOrderRepository) FindById(uint) (models.Order, error) { return r.current, nil }
func (r *paymentOrderRepository) UpdatePayment(order *models.Order, _ time.Time) error {
	r.current = *order
	r.updated = true
	return nil
}

type countedPaymentStoreConfig struct {
	paymentStoreConfig
	calls int
}

func (s *countedPaymentStoreConfig) GetStoreConfig() (models.StoreConfig, error) {
	s.calls++
	config := s.config
	if config.UpdatedAt.IsZero() {
		config.UpdatedAt = paymentTestVersion
	}
	return config, nil
}

func TestPixPaymentUsesOneConfigurationSnapshot(t *testing.T) {
	repo := &paymentOrderRepository{current: models.Order{Profile: paymentProfile(), Model: gorm.Model{ID: 12}, Status: models.ActiveStatus, ShoppingCart: models.ShoppingCart{Total: 10}}}
	store := &countedPaymentStoreConfig{paymentStoreConfig: paymentStoreConfig{config: models.StoreConfig{PhysicalStoreAddress: "Rua da Loja, 12", PhysicalStoreCity: "São Paulo", PhysicalStoreState: "SP", PixReceiverName: "Cupcake Store", DeliveryIsActive: true, DeliveryPrice: 2.50, PaymentPixIsActive: true, PixKeyType: models.PixTypeEmail, PixKey: "shop@example.com"}}}
	service := NewOrderService(repo, store).(*orderService)
	service.generatePix = func(data *models.PixPaymentData) (*models.PixInfo, error) {
		if data.Valor != "12.50" || data.Chave != "shop@example.com" || data.Tipo != "email" {
			t.Fatalf("incorrect invoice data: %+v", data)
		}
		if data.Nome != "Cupcake Store" || data.City != "São Paulo" || len(data.Txid) != 25 {
			t.Fatalf("invalid merchant or transaction: %+v", data)
		}
		return generateLocalPixPayment(data)
	}
	order := models.Order{Model: gorm.Model{ID: 12}, PaymentMethod: models.PixPaymentMethod, IsDelivery: true, Total: 0.01}
	if err := service.Payment(&order, paymentTestVersion); err != nil {
		t.Fatal(err)
	}
	if store.calls != 1 || !repo.updated || order.Status != models.AwaitingPaymentStatus || order.Total != 12.50 || order.PixURL != "" || order.PixString == "" {
		t.Fatalf("unexpected payment: config calls=%d order=%+v", store.calls, order)
	}
}

func TestPixRetryPreservesInvoiceAndRecordedDeliveryTerms(t *testing.T) {
	repo := &paymentOrderRepository{current: models.Order{Profile: paymentProfile(), Model: gorm.Model{ID: 12}, Status: models.AwaitingPaymentStatus, PaymentMethod: models.PixPaymentMethod, IsDelivery: true, DeliveryPrice: 2.50, Total: 12.50, ShoppingCart: models.ShoppingCart{Total: 10}, PixURL: "/invoice/12", PixString: paymentPix(t).PixString, PixTransactionID: paymentPix(t).PixTransactionID}}
	// Nil configuration and generator prove that retrying does not depend on
	// current availability/pricing or regenerate the payment instructions.
	service := &orderService{orderRepository: repo}
	order := models.Order{Model: gorm.Model{ID: 12}, PaymentMethod: models.PixPaymentMethod, IsDelivery: true, Total: 99}
	if err := service.Payment(&order, paymentTestVersion); err != nil {
		t.Fatal(err)
	}
	if repo.updated || order.Total != 12.50 || order.DeliveryPrice != 2.50 || order.PixString != paymentPix(t).PixString {
		t.Fatalf("retry changed invoice: %+v", order)
	}
}

func TestPendingPixCannotChangePaymentOrDelivery(t *testing.T) {
	for _, tc := range []struct {
		method   models.PaymentMethod
		delivery bool
	}{
		{models.CashPaymentMethod, true},
		{models.PixPaymentMethod, false},
	} {
		repo := &paymentOrderRepository{current: models.Order{Profile: paymentProfile(), Model: gorm.Model{ID: 12}, Status: models.AwaitingPaymentStatus, PaymentMethod: models.PixPaymentMethod, IsDelivery: true, DeliveryPrice: 2.50, Total: 12.50, ShoppingCart: models.ShoppingCart{Total: 10}, PixURL: "/invoice/12", PixString: paymentPix(t).PixString, PixTransactionID: paymentPix(t).PixTransactionID}}
		service := &orderService{orderRepository: repo}
		order := models.Order{Model: gorm.Model{ID: 12}, PaymentMethod: tc.method, IsDelivery: tc.delivery}
		if err := service.Payment(&order, paymentTestVersion); err == nil || repo.updated {
			t.Fatalf("changed pending invoice: method=%s delivery=%v err=%v", tc.method, tc.delivery, err)
		}
	}
}

func TestPaymentRejectsUnavailableDeliveryWithoutChangingToPickup(t *testing.T) {
	repo := &paymentOrderRepository{current: models.Order{Profile: paymentProfile(), Model: gorm.Model{ID: 1}, Status: models.ActiveStatus, ShoppingCart: models.ShoppingCart{Total: 10}}}
	order := models.Order{Model: gorm.Model{ID: 1}, PaymentMethod: models.CashPaymentMethod, IsDelivery: true}
	err := NewOrderService(repo, paymentStoreConfig{config: models.StoreConfig{PhysicalStoreAddress: "Rua da Loja, 12", PhysicalStoreCity: "São Paulo", PhysicalStoreState: "SP", PixReceiverName: "Cupcake Store", PaymentCashIsActive: true}}).Payment(&order, paymentTestVersion)
	if err == nil || repo.updated {
		t.Fatalf("unavailable delivery silently accepted: %+v err=%v", order, err)
	}
}

func TestCashPaymentClearsUnsubmittedPixDetails(t *testing.T) {
	repo := &paymentOrderRepository{current: models.Order{Profile: paymentProfile(), Model: gorm.Model{ID: 1}, Status: models.ActiveStatus, PaymentMethod: models.PixPaymentMethod, ShoppingCart: models.ShoppingCart{Total: 10}, PixURL: "/invoice/old", PixString: "old-code", PixQR: "old-image", PixTransactionID: "old-id"}}
	order := models.Order{Model: gorm.Model{ID: 1}, PaymentMethod: models.CashPaymentMethod}
	if err := NewOrderService(repo, paymentStoreConfig{config: models.StoreConfig{PhysicalStoreAddress: "Rua da Loja, 12", PhysicalStoreCity: "São Paulo", PhysicalStoreState: "SP", PixReceiverName: "Cupcake Store", PaymentCashIsActive: true}}).Payment(&order, paymentTestVersion); err != nil {
		t.Fatal(err)
	}
	if !repo.updated || order.PixURL != "" || order.PixString != "" || order.PixQR != "" || order.PixTransactionID != "" {
		t.Fatalf("cash payment retained stale Pix details: %+v", order)
	}
}

func TestPixFailureDoesNotPersistPayment(t *testing.T) {
	for _, tc := range []struct {
		name    string
		payment *models.PixInfo
		err     error
	}{
		{name: "generator failure", err: errors.New("unavailable")},
		{name: "missing response"},
		{name: "missing code", payment: &models.PixInfo{PixURL: "/invoice/12"}},
		{name: "invalid payload", payment: &models.PixInfo{PixURL: ".evil.example/12", PixString: "code"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &paymentOrderRepository{current: models.Order{Profile: paymentProfile(), Model: gorm.Model{ID: 12}, Status: models.ActiveStatus, ShoppingCart: models.ShoppingCart{Total: 10}}}
			service := NewOrderService(repo, paymentStoreConfig{config: models.StoreConfig{PhysicalStoreAddress: "Rua da Loja, 12", PhysicalStoreCity: "São Paulo", PhysicalStoreState: "SP", PixReceiverName: "Cupcake Store", PaymentPixIsActive: true, PixKeyType: models.PixTypeEmail, PixKey: "shop@example.com"}}).(*orderService)
			service.generatePix = func(*models.PixPaymentData) (*models.PixInfo, error) { return tc.payment, tc.err }
			order := models.Order{Model: gorm.Model{ID: 12}, PaymentMethod: models.PixPaymentMethod}
			if err := service.Payment(&order, paymentTestVersion); err == nil || repo.updated {
				t.Fatalf("invalid Pix persisted: err=%v updated=%v", err, repo.updated)
			}
		})
	}
}

type paymentStoreConfig struct {
	StoreConfigService
	config models.StoreConfig
}

func (s paymentStoreConfig) GetStoreConfig() (models.StoreConfig, error) {
	config := s.config
	if config.UpdatedAt.IsZero() {
		config.UpdatedAt = paymentTestVersion
	}
	return config, nil
}

func TestCashPaymentUsesStoredSubtotalAndDeliveryChoice(t *testing.T) {
	for _, delivery := range []bool{false, true} {
		t.Run(map[bool]string{false: "pickup", true: "delivery"}[delivery], func(t *testing.T) {
			repo := &paymentOrderRepository{current: models.Order{Profile: paymentProfile(), Model: gorm.Model{ID: 1}, Status: models.ActiveStatus, ShoppingCart: models.ShoppingCart{Total: 10}}}
			store := paymentStoreConfig{config: models.StoreConfig{PhysicalStoreAddress: "Rua da Loja, 12", PhysicalStoreCity: "São Paulo", PhysicalStoreState: "SP", PixReceiverName: "Cupcake Store", DeliveryIsActive: true, DeliveryPrice: 2.50, PaymentCashIsActive: true}}
			order := models.Order{Model: gorm.Model{ID: 1}, PaymentMethod: models.CashPaymentMethod, IsDelivery: delivery, Total: 0.01, ShoppingCart: models.ShoppingCart{Total: 0.01}}
			if err := NewOrderService(repo, store).Payment(&order, paymentTestVersion); err != nil {
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
		repo := &paymentOrderRepository{current: models.Order{Profile: paymentProfile(), Model: gorm.Model{ID: 1}, Status: status, ShoppingCart: models.ShoppingCart{Total: 10}}}
		order := models.Order{Model: gorm.Model{ID: 1}, PaymentMethod: models.CashPaymentMethod}
		if err := NewOrderService(repo, paymentStoreConfig{}).Payment(&order, paymentTestVersion); err == nil || repo.updated {
			t.Fatalf("invalid payment accepted for %s", status)
		}
	}
}

func paymentProfile() models.Profile {
	return models.Profile{FirstName: "Cliente", LastName: "Teste", PhoneNumber: "11999999999", Address: "Rua Cliente, 5", City: "São Paulo", State: "SP", PostalCode: "01001000"}
}

func paymentPix(t *testing.T) *models.PixInfo {
	t.Helper()
	payment, err := generateLocalPixPayment(&models.PixPaymentData{Tipo: "email", Chave: "shop@example.com", Valor: "12.50", Nome: "Cupcake Store", City: "São Paulo", Txid: "C123"})
	if err != nil {
		t.Fatal(err)
	}
	return payment
}

func TestManualConfirmationRejectsAlteredPaymentInstructions(t *testing.T) {
	pix := paymentPix(t)
	for _, tc := range []struct {
		name    string
		total   float64
		txid    string
		wantErr bool
	}{
		{"valid", 12.50, pix.PixTransactionID, false},
		{"wrong amount", 0.01, pix.PixTransactionID, true},
		{"wrong identifier", 12.50, "OTHER", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &paymentOrderRepository{current: models.Order{Model: gorm.Model{ID: 12}, Status: models.AwaitingPaymentStatus, PaymentMethod: models.PixPaymentMethod, Total: tc.total, PixString: pix.PixString, PixTransactionID: tc.txid}}
			err := NewOrderService(repo, nil).ConfirmPayment(12, 8)
			if (err != nil) != tc.wantErr || (repo.confirmedBy == 0) != tc.wantErr {
				t.Fatalf("err=%v confirmedBy=%d", err, repo.confirmedBy)
			}
		})
	}
}

var paymentTestVersion = time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

func TestPaymentRequiresTheConditionsShownAtCheckout(t *testing.T) {
	for _, method := range []models.PaymentMethod{models.CashPaymentMethod, models.PixPaymentMethod} {
		for _, expectedVersion := range []time.Time{{}, paymentTestVersion.Add(-time.Second)} {
			repo := &paymentOrderRepository{current: models.Order{Model: gorm.Model{ID: 12}, Profile: paymentProfile(), Status: models.ActiveStatus, ShoppingCart: models.ShoppingCart{Total: 10}}}
			service := NewOrderService(repo, paymentStoreConfig{config: models.StoreConfig{Model: gorm.Model{UpdatedAt: paymentTestVersion}, PaymentCashIsActive: true, PaymentPixIsActive: true, PixKeyType: models.PixTypeEmail, PixKey: "shop@example.com", PixReceiverName: "Cupcake Store", PhysicalStoreAddress: "Rua 12", PhysicalStoreCity: "São Paulo", PhysicalStoreState: "SP"}}).(*orderService)
			service.generatePix = func(*models.PixPaymentData) (*models.PixInfo, error) {
				t.Fatal("generated payment before customer reviewed changed conditions")
				return nil, nil
			}
			order := models.Order{Model: gorm.Model{ID: 12}, PaymentMethod: method}
			if err := service.Payment(&order, expectedVersion); err == nil || repo.updated {
				t.Fatalf("accepted stale/missing conditions: method=%s err=%v", method, err)
			}
		}
	}
}
