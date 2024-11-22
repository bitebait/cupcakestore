package services

import (
	"errors"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"math"
)

type ShoppingCartService interface {
	FindOrCreateById(id uint) (models.ShoppingCart, error)
	FindOrCreateByUserId(userID uint) (models.ShoppingCart, error)
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

func (s *shoppingCartService) FindOrCreateById(id uint) (models.ShoppingCart, error) {
	cart, err := s.shoppingCartRepository.FindOrCreateById(id)

	if err != nil {
		err = errors.New("falha ao criar ou encontrar o carrinho de compras")
	}

	return cart, err
}

func (s *shoppingCartService) FindOrCreateByUserId(userID uint) (models.ShoppingCart, error) {
	cart, err := s.shoppingCartRepository.FindOrCreateByUserId(userID)

	if err != nil {
		err = errors.New("falha ao criar ou encontrar o carrinho de compras do usuário")
	}

	return cart, err
}

func (s *shoppingCartService) AddItemToCart(userID, productID uint, quantity int) error {
	cart, err := s.FindOrCreateByUserId(userID)

	if err != nil {
		return errors.New("falha ao encontrar o carrinho de compras do usuário")
	}

	for _, item := range cart.Items {
		if item.ProductID == productID {
			item.Quantity += int(math.Abs(float64(quantity)))
			err := s.shoppingCartItemService.Update(&item)
			if err != nil {
				return errors.New("falha ao atualizar o item do carrinho")
			}
			return nil
		}
	}

	if err := s.shoppingCartItemService.Create(cart.ID, productID, quantity); err != nil {
		return errors.New("falha ao adicionar o item ao carrinho")
	}

	return nil
}

func (s *shoppingCartService) RemoveFromCart(userID, productID uint) error {
	cart, err := s.FindOrCreateByUserId(userID)

	if err != nil {
		return errors.New("falha ao encontrar o carrinho de compras do usuário")
	}

	if err := s.shoppingCartItemService.Delete(cart.ID, productID); err != nil {
		return errors.New("falha ao remover o item do carrinho")
	}

	return nil
}
