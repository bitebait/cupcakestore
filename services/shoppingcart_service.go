package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type ShoppingCartService interface {
	FindByUserId(id uint) (models.ShoppingCart, error)
	AddItemToCart(userID uint, productID uint, quantity int) error
	RemoveFromCart(userID uint, productID uint) error
}

type shoppingCartService struct {
	shoppingCartRepository  repositories.ShoppingCartRepository
	shoppingCartItemService ShoppingCartItemService
}

func NewShoppingCartService(shoppingCartRepository repositories.ShoppingCartRepository, shoppingCartItemService ShoppingCartItemService) ShoppingCartService {
	return &shoppingCartService{
		shoppingCartRepository:  shoppingCartRepository,
		shoppingCartItemService: shoppingCartItemService,
	}
}

func (s shoppingCartService) FindByUserId(id uint) (models.ShoppingCart, error) {
	return s.shoppingCartRepository.FindByUserId(id)
}

func (s shoppingCartService) AddItemToCart(userID uint, productID uint, quantity int) error {
	cart, err := s.FindByUserId(userID)
	if err != nil {
		return err
	}

	for _, item := range cart.Items {
		if item.ProductID == productID {
			item.Quantity += quantity
			return s.shoppingCartItemService.Update(&item)
		}
	}

	return s.shoppingCartItemService.Create(cart.ID, productID, quantity)
}

func (s shoppingCartService) RemoveFromCart(userID uint, productID uint) error {
	cart, err := s.FindByUserId(userID)
	if err != nil {
		return err
	}

	return s.shoppingCartItemService.Delete(cart.ID, productID)
}
