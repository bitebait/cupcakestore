package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type ShoppingCartService interface {
	FindByUserId(id uint) (models.ShoppingCart, error)
	AddItemToCart(userID uint, productID uint, quantity int) error
}

type shoppingCartService struct {
	shoppingCartRepository repositories.ShoppingCartRepository
}

func NewShoppingCartService(shoppingCartRepository repositories.ShoppingCartRepository) ShoppingCartService {
	return &shoppingCartService{
		shoppingCartRepository: shoppingCartRepository,
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
			return s.shoppingCartRepository.UpdateCartItem(&item)
		}
	}

	cartItem := &models.ShoppingCartItem{
		ShoppingCartID: cart.ID,
		ProductID:      productID,
		Quantity:       quantity,
	}
	return s.shoppingCartRepository.AddItemToCart(cartItem)
}
