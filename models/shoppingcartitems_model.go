package models

import (
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
	var product Product
	if err := tx.First(&product, item.ProductID).Error; err != nil {
		return err
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

	shoppingCart.updateTotal()

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
