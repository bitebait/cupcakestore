package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type OrderService interface {
	FindAll(filter *models.ShoppingCartFilter) []models.ShoppingCart
}

type orderService struct {
	orderRepository repositories.OrderRepository
}

func NewOrderService(orderRepository repositories.OrderRepository) OrderService {
	return &orderService{
		orderRepository: orderRepository,
	}
}

func (o *orderService) FindAll(filter *models.ShoppingCartFilter) []models.ShoppingCart {
	return o.orderRepository.FindAll(filter)
}
