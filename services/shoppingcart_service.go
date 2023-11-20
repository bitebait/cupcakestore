package services

import (
	"fmt"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type ShoppingCartService interface {
	FindById(id uint) (models.ShoppingCart, error)
	FindByUserId(userID uint) (models.ShoppingCart, error)
	AddItemToCart(userID, productID uint, quantity int) error
	RemoveFromCart(userID, productID uint) error
	Update(cart *models.ShoppingCart) error
	Payment(cart *models.ShoppingCart) error
}

type shoppingCartServiceImpl struct {
	shoppingCartRepository  repositories.ShoppingCartRepository
	shoppingCartItemService ShoppingCartItemService
	storeConfigService      StoreConfigService
}

func NewShoppingCartService(shoppingCartRepository repositories.ShoppingCartRepository, shoppingCartItemService ShoppingCartItemService, storeConfigService StoreConfigService) ShoppingCartService {
	return &shoppingCartServiceImpl{
		shoppingCartRepository:  shoppingCartRepository,
		shoppingCartItemService: shoppingCartItemService,
		storeConfigService:      storeConfigService,
	}
}

func (s *shoppingCartServiceImpl) FindById(id uint) (models.ShoppingCart, error) {
	return s.shoppingCartRepository.FindById(id)
}

func (s *shoppingCartServiceImpl) FindByUserId(userID uint) (models.ShoppingCart, error) {
	return s.shoppingCartRepository.FindByUserId(userID)
}

func (s *shoppingCartServiceImpl) AddItemToCart(userID, productID uint, quantity int) error {
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

func (s *shoppingCartServiceImpl) RemoveFromCart(userID, productID uint) error {
	cart, err := s.FindByUserId(userID)
	if err != nil {
		return err
	}

	return s.shoppingCartItemService.Delete(cart.ID, productID)
}

func (s *shoppingCartServiceImpl) Update(cart *models.ShoppingCart) error {
	return s.shoppingCartRepository.Update(cart)
}

func (s *shoppingCartServiceImpl) Payment(cart *models.ShoppingCart) error {
	var err error
	cart.Status = models.AwaitingPaymentStatus

	if cart.PaymentMethod == models.CashPaymentMethod {
		cart.Status = models.ProcessingStatus
	} else if cart.PaymentMethod == models.PixPaymentMethod {
		err = s.processPixPayment(cart)
	}

	if err != nil {
		return err
	}

	return s.shoppingCartRepository.Update(cart)
}

func (s *shoppingCartServiceImpl) processPixPayment(cart *models.ShoppingCart) error {
	storeConfig, err := s.storeConfigService.GetStoreConfig()
	if err != nil {
		return err
	}

	pixData := &models.PixPaymentData{
		Tipo:  string(storeConfig.PixKeyType),
		Chave: storeConfig.PixKey,
		Valor: fmt.Sprintf("%.2f", cart.Total),
		Info:  fmt.Sprintf("CupCake Store R$ %v - ID#%v", cart.Total, cart.ID),
		Nome:  "Cupcake Store",
	}

	payment, err := models.GeneratePixPayment(pixData)
	if err != nil {
		return err
	}

	cart.PixQR = payment.PixQR
	cart.PixString = payment.PixString
	cart.PixTransactionID = payment.PixTransactionID
	cart.PixURL = payment.PixURL

	return nil
}
