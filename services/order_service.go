package services

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type OrderService interface {
	FindById(id uint) (models.Order, error)
	FindByCartId(cartID uint) (models.Order, error)
	FindOrCreate(profileID, cartID uint) (models.Order, error)
	FindAll(filter *models.OrderFilter) []models.Order
	FindAllByUser(filter *models.OrderFilter) []models.Order
	Update(order *models.Order) error
	Payment(order *models.Order, expectedConfigVersion time.Time) error
	Cancel(id uint) error
	CancelForCustomer(id, profileID uint) error
	ConfirmPayment(id, staffProfileID uint) error
}

type orderService struct {
	orderRepository    repositories.OrderRepository
	storeConfigService StoreConfigService
	generatePix        func(*models.PixPaymentData) (*models.PixInfo, error)
}

func NewOrderService(orderRepository repositories.OrderRepository, storeConfigService StoreConfigService) OrderService {
	return &orderService{
		orderRepository:    orderRepository,
		storeConfigService: storeConfigService,
		generatePix:        generateLocalPixPayment,
	}
}

func (s *orderService) FindById(id uint) (models.Order, error) {
	order, err := s.orderRepository.FindById(id)

	if err != nil {
		err = errors.New("falha ao encontrar o pedido pelo id")
	}

	return order, err
}

func (s *orderService) FindByCartId(cartID uint) (models.Order, error) {
	order, err := s.orderRepository.FindByCartId(cartID)

	if err != nil {
		err = errors.New("falha ao encontrar o pedido pelo id do carrinho")
	}

	return order, err
}

func (s *orderService) FindOrCreate(profileID, cartID uint) (models.Order, error) {
	order, err := s.orderRepository.FindOrCreate(profileID, cartID)

	if err != nil {
		if errors.Is(err, models.ErrContactIncomplete) || errors.Is(err, models.ErrDeliveryAddressIncomplete) || errors.Is(err, models.ErrFulfillmentUnavailable) || errors.Is(err, models.ErrPaymentUnavailable) {
			return order, err
		}
		err = errors.New("falha ao criar ou encontrar o pedido")
	}

	return order, err
}

func (s *orderService) FindAll(filter *models.OrderFilter) []models.Order {
	return s.orderRepository.FindAll(filter)
}

func (s *orderService) FindAllByUser(filter *models.OrderFilter) []models.Order {
	return s.orderRepository.FindAllByUser(filter)
}

func (s *orderService) Update(order *models.Order) error {
	if err := s.orderRepository.Update(order); err != nil {
		return errors.New("falha ao atualizar o pedido")
	}

	return nil
}

func (s *orderService) Payment(order *models.Order, expectedConfigVersion time.Time) error {
	if order == nil {
		return errors.New("pedido inválido")
	}
	paymentMethod, isDelivery := order.PaymentMethod, order.IsDelivery
	current, err := s.orderRepository.FindById(order.ID)
	if err != nil || !current.CanProceedToPayment() {
		return errors.New("o pedido não pode seguir para pagamento")
	}
	// An already issued Pix belongs to the recorded delivery terms. Reopening
	// it must not issue another invoice or apply a subsequently changed fee.
	if current.CanRedirectToPixPayment() {
		if current.PaymentMethod != paymentMethod || current.IsDelivery != isDelivery {
			return errors.New("este pedido já possui um Pix emitido; entre em contato com a loja para alterar o pagamento ou a entrega")
		}
		if err := ValidateOrderPix(&current); err != nil {
			return errors.New("os dados do Pix deste pedido estão inválidos; entre em contato com a loja")
		}
		*order = current
		return nil
	}
	current.PaymentMethod = paymentMethod
	*order = current
	if err := order.Validate(); err != nil {
		return err
	}
	storeConfig, err := s.storeConfigService.GetStoreConfig()
	if err != nil {
		return errors.New("falha ao carregar as formas de pagamento")
	}
	if expectedConfigVersion.IsZero() || !storeConfig.UpdatedAt.Equal(expectedConfigVersion) {
		return errors.New("as condições da loja mudaram; confira o valor e a forma de recebimento antes de continuar")
	}
	if err := storeConfig.ValidateFulfillment(current.Profile, isDelivery); err != nil {
		return err
	}
	order.IsDelivery = isDelivery
	order.DeliveryPrice = 0
	if order.IsDelivery {
		order.DeliveryPrice = storeConfig.DeliveryPrice
	}
	if order.DeliveryPrice < 0 || math.IsNaN(order.DeliveryPrice) || math.IsInf(order.DeliveryPrice, 0) {
		return errors.New("taxa de entrega inválida")
	}
	order.Total = math.Round((order.ShoppingCart.Total+order.DeliveryPrice)*100) / 100
	if math.IsNaN(order.Total) || math.IsInf(order.Total, 0) || order.Total <= 0 {
		return errors.New("valor total do pedido inválido")
	}
	switch order.PaymentMethod {
	case models.CashPaymentMethod:
		if !storeConfig.PaymentCashIsActive {
			return errors.New("pagamento em dinheiro indisponível")
		}
		order.PixQR, order.PixString, order.PixTransactionID, order.PixURL = "", "", "", ""
		order.Status = models.ProcessingStatus
	case models.PixPaymentMethod:
		if !storeConfig.IsPixAvailable() {
			return errors.New("pagamento com Pix indisponível")
		}
		if err := s.processPixPayment(order, storeConfig); err != nil {
			return errors.New("falha ao processar o pagamento com Pix")
		}
		order.Status = models.AwaitingPaymentStatus
	}
	if err := s.orderRepository.UpdatePayment(order, storeConfig.UpdatedAt); err != nil {
		return errors.New("falha ao atualizar o status do pedido")
	}
	return nil
}

func (s *orderService) processPixPayment(order *models.Order, storeConfig models.StoreConfig) error {
	orderTotalPrice := fmt.Sprintf("%.2f", order.Total)
	transactionHash := sha256.Sum256([]byte(fmt.Sprintf("%d:%s", order.ID, order.CreatedAt.UTC().Format("20060102T150405.000000000"))))

	pixData := &models.PixPaymentData{
		Tipo:  string(storeConfig.PixKeyType),
		Chave: storeConfig.PixKey,
		Valor: orderTotalPrice,
		Nome:  storeConfig.PixReceiverName,
		City:  storeConfig.PhysicalStoreCity,
		Txid:  fmt.Sprintf("C%X", transactionHash[:12]),
	}

	payment, err := s.generatePix(pixData)

	if err != nil {
		return err
	}
	if payment == nil || payment.PixString == "" {
		return errors.New("dados do Pix ausentes")
	}

	order.PaymentMethod = models.PixPaymentMethod
	order.PixQR = payment.PixQR
	order.PixString = payment.PixString
	order.PixTransactionID = payment.PixTransactionID
	order.PixURL = payment.PixURL
	return ValidateOrderPix(order)
}

func (s *orderService) Cancel(id uint) error {
	if err := s.orderRepository.Cancel(id); err != nil {
		return errors.New("falha ao cancelar o pedido")
	}

	return nil
}

func (s *orderService) CancelForCustomer(id, profileID uint) error {
	return s.orderRepository.CancelForCustomer(id, profileID)
}

func (s *orderService) ConfirmPayment(id, staffProfileID uint) error {
	order, err := s.orderRepository.FindById(id)
	if err != nil {
		return errors.New("pedido não encontrado")
	}
	if err := ValidateOrderPix(&order); err != nil {
		return errors.New("o Pix deste pedido é inválido; confira os dados antes de confirmar")
	}
	if err := s.orderRepository.ConfirmPayment(id, staffProfileID); err != nil {
		return errors.New("não foi possível confirmar o recebimento; recarregue o pedido")
	}
	return nil
}
