package services

import (
	"errors"
	"fmt"
	"math"

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
	Payment(order *models.Order) error
	Cancel(id uint) error
}

type orderService struct {
	orderRepository    repositories.OrderRepository
	storeConfigService StoreConfigService
}

func NewOrderService(orderRepository repositories.OrderRepository, storeConfigService StoreConfigService) OrderService {
	return &orderService{
		orderRepository:    orderRepository,
		storeConfigService: storeConfigService,
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

func (s *orderService) Payment(order *models.Order) error {
	if order == nil {
		return errors.New("pedido inválido")
	}
	paymentMethod, isDelivery := order.PaymentMethod, order.IsDelivery
	current, err := s.orderRepository.FindById(order.ID)
	if err != nil || !current.CanProceedToPayment() {
		return errors.New("o pedido não pode seguir para pagamento")
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
	order.IsDelivery = isDelivery && storeConfig.DeliveryIsActive
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
		order.Status = models.ProcessingStatus
	case models.PixPaymentMethod:
		if !storeConfig.PaymentPixIsActive {
			return errors.New("pagamento com Pix indisponível")
		}
		if err := s.processPixPayment(order); err != nil {
			return errors.New("falha ao processar o pagamento com Pix")
		}
		order.Status = models.AwaitingPaymentStatus
	}
	if err := s.orderRepository.Update(order); err != nil {
		return errors.New("falha ao atualizar o status do pedido")
	}
	return nil
}

func (s *orderService) processPixPayment(order *models.Order) error {
	storeConfig, err := s.storeConfigService.GetStoreConfig()

	if err != nil {
		return err
	}

	orderTotalPrice := fmt.Sprintf("%.2f", order.Total)

	pixData := &models.PixPaymentData{
		Tipo:  string(storeConfig.PixKeyType),
		Chave: storeConfig.PixKey,
		Valor: orderTotalPrice,
		Info:  fmt.Sprintf("CupCake Store R$ %v - ID#%v", orderTotalPrice, order.ID),
		Nome:  "Cupcake Store",
	}

	payment, err := generatePixPayment(pixData)

	if err != nil {
		return err
	}

	order.PaymentMethod = models.PixPaymentMethod
	order.PixQR = payment.PixQR
	order.PixString = payment.PixString
	order.PixTransactionID = payment.PixTransactionID
	order.PixURL = payment.PixURL

	return nil
}

func (s *orderService) Cancel(id uint) error {
	if err := s.orderRepository.Cancel(id); err != nil {
		return errors.New("falha ao cancelar o pedido")
	}

	return nil
}
