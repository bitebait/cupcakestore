package services

import (
	"errors"
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

	if err := s.shoppingCartRepository.Create(item); err != nil {
		return errors.New("falha ao criar o item do carrinho")
	}

	return nil
}

func (s shoppingCartItemService) Update(item *models.ShoppingCartItem) error {
	if err := s.shoppingCartRepository.Update(item); err != nil {
		return errors.New("falha ao atualizar o item do carrinho")
	}

	return nil
}

func (s shoppingCartItemService) Delete(cartID uint, productID uint) error {
	if err := s.shoppingCartRepository.Delete(cartID, productID); err != nil {
		return errors.New("falha ao deletar o item do carrinho")
	}

	return nil
}
