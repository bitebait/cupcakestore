package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"math"
)

type ShoppingCartService interface {
	FindById(id uint) (models.ShoppingCart, error)
	FindByUserId(userID uint) (models.ShoppingCart, error)
	AddItemToCart(userID, productID uint, quantity int) error
	RemoveFromCart(userID, productID uint) error
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

func (s *shoppingCartService) FindById(id uint) (models.ShoppingCart, error) {
	return s.shoppingCartRepository.FindById(id)
}

func (s *shoppingCartService) FindByUserId(userID uint) (models.ShoppingCart, error) {
	return s.shoppingCartRepository.FindByUserId(userID)
}

func (s *shoppingCartService) AddItemToCart(userID, productID uint, quantity int) error {
	cart, err := s.FindByUserId(userID)
	if err != nil {
		return err
	}

	for _, item := range cart.Items {
		if item.ProductID == productID {
			item.Quantity += int(math.Abs(float64(quantity)))
			return s.shoppingCartItemService.Update(&item)
		}
	}

	return s.shoppingCartItemService.Create(cart.ID, productID, int(math.Abs(float64(quantity))))
}

func (s *shoppingCartService) RemoveFromCart(userID, productID uint) error {
	cart, err := s.FindByUserId(userID)
	if err != nil {
		return err
	}

	return s.shoppingCartItemService.Delete(cart.ID, productID)
}
