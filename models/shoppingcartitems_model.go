package models

import (
	"errors"
	"math"

	"gorm.io/gorm"
)

type ShoppingCartItem struct {
	gorm.Model
	ProductID      uint    `gorm:"not null" validate:"required"`
	Product        Product `validate:"-"`
	ItemPrice      float64 `gorm:"default:0"`
	Quantity       int     `gorm:"default:1"`
	ShoppingCartID uint    `gorm:"not null"`
}

func (item *ShoppingCartItem) BeforeSave(tx *gorm.DB) error {
	if item.Quantity <= 0 || item.ShoppingCartID == 0 || item.ProductID == 0 {
		return errors.New("produto, carrinho e quantidade válida devem ser informados")
	}
	if err := item.ensureOpenCart(tx); err != nil {
		return err
	}
	var product Product
	if err := tx.First(&product, item.ProductID).Error; err != nil {
		return err
	}
	if !product.IsActive || product.CurrentStock < item.Quantity {
		return errors.New("produto indisponível ou estoque insuficiente")
	}
	if product.Price <= 0 || math.IsNaN(product.Price) || math.IsInf(product.Price, 0) {
		return errors.New("preço do produto inválido")
	}
	item.ItemPrice = product.Price
	return nil
}

func (item *ShoppingCartItem) updateShoppingCartTotal(tx *gorm.DB) error {
	shoppingCart := &ShoppingCart{}
	if err := tx.Preload("Items").First(shoppingCart, item.ShoppingCartID).Error; err != nil {
		return err
	}
	previousTotal := shoppingCart.Total

	if err := shoppingCart.updateTotal(); err != nil {
		return err
	}

	if shoppingCart.Total != previousTotal {
		if len(shoppingCart.Items) <= 0 {
			shoppingCart.Total = 0.0
		}
		if err := tx.Model(shoppingCart).Select("Total").Updates(shoppingCart).Error; err != nil {
			return err
		}
	}

	return nil
}

func (item *ShoppingCartItem) AfterSave(tx *gorm.DB) error {
	return item.updateShoppingCartTotal(tx)
}

func (item *ShoppingCartItem) AfterDelete(tx *gorm.DB) error {
	return item.updateShoppingCartTotal(tx)
}

func (item *ShoppingCartItem) ensureOpenCart(tx *gorm.DB) error {
	var cart ShoppingCart
	if err := tx.First(&cart, item.ShoppingCartID).Error; err != nil {
		return err
	}
	if cart.OrderID != 0 {
		return errors.New("o carrinho já foi finalizado")
	}
	return nil
}

func (item *ShoppingCartItem) BeforeDelete(tx *gorm.DB) error {
	return item.ensureOpenCart(tx)
}
