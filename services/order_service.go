package services

import (
	"fmt"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type OrderService interface {
	FindById(id uint) (models.Order, error)
	FindByCartId(cartID uint) (*models.Order, error)
	FindOrCreate(profileID, cartID uint) (*models.Order, error)
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
	return s.orderRepository.FindById(id)
}

func (s *orderService) FindByCartId(cartID uint) (*models.Order, error) {
	return s.orderRepository.FindByCartId(cartID)
}

func (s *orderService) FindOrCreate(profileID, cartID uint) (*models.Order, error) {
	return s.orderRepository.FindOrCreate(profileID, cartID)
}

func (s *orderService) FindAll(filter *models.OrderFilter) []models.Order {
	return s.orderRepository.FindAll(filter)
}

func (s *orderService) FindAllByUser(filter *models.OrderFilter) []models.Order {
	return s.orderRepository.FindAllByUser(filter)
}

func (s *orderService) Update(order *models.Order) error {
	return s.orderRepository.Update(order)
}

const (
	CashPaymentMethod = models.CashPaymentMethod
	PixPaymentMethod  = models.PixPaymentMethod
)

func (s *orderService) Payment(order *models.Order) error {
	order.Status = models.AwaitingPaymentStatus
	switch order.PaymentMethod {
	case CashPaymentMethod:
		order.Status = models.ProcessingStatus
	case PixPaymentMethod:
		if err := s.processPixPayment(order); err != nil {
			return err
		}
	}
	return s.orderRepository.Update(order)
}

func (s *orderService) processPixPayment(order *models.Order) error {
	storeConfig, err := s.storeConfigService.GetStoreConfig()
	if err != nil {
		return err
	}
	pixData := &models.PixPaymentData{
		Tipo:  string(storeConfig.PixKeyType),
		Chave: storeConfig.PixKey,
		Valor: fmt.Sprintf("%.2f", order.Total),
		Info:  fmt.Sprintf("CupCake Store R$ %v - ID#%v", order.Total, order.ID),
		Nome:  "Cupcake Store",
	}
	payment, err := models.GeneratePixPayment(pixData)
	if err != nil {
		return err
	}
	order.PixQR = payment.PixQR
	order.PixString = payment.PixString
	order.PixTransactionID = payment.PixTransactionID
	order.PixURL = payment.PixURL
	return nil
}

func (s *orderService) Cancel(id uint) error {
	return s.orderRepository.Cancel(id)
}
