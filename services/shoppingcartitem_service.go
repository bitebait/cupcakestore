package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"math"
)

type ShoppingCartItemService interface {
	Create(cartID uint, productID uint, quantity int) error
	Update(item *models.ShoppingCartItem) error
	Delete(cartID uint, productID uint) error
}

type shoppingCartItemService struct {
	shoppingCartRepository repositories.ShoppingCartItemRepository
}

func NewShoppingCartItemService(shoppingCartRepository repositories.ShoppingCartItemRepository) ShoppingCartItemService {
	return &shoppingCartItemService{
		shoppingCartRepository: shoppingCartRepository,
	}
}

func (s shoppingCartItemService) Create(cartID uint, productID uint, quantity int) error {
	item := &models.ShoppingCartItem{
		ShoppingCartID: cartID,
		ProductID:      productID,
		Quantity:       int(math.Abs(float64(quantity))),
	}

	return s.shoppingCartRepository.Create(item)
}

func (s shoppingCartItemService) Update(item *models.ShoppingCartItem) error {
	return s.shoppingCartRepository.Update(item)
}

func (s shoppingCartItemService) Delete(cartID uint, productID uint) error {
	return s.shoppingCartRepository.Delete(cartID, productID)
}
