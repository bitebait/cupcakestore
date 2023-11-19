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

func (item *ShoppingCartItem) BeforeCreate(tx *gorm.DB) error {
	var product Product
	if err := tx.First(&product, item.ProductID).Error; err != nil {
		return err
	}
	item.ItemPrice = product.Price
	return nil
}

func (item *ShoppingCartItem) AfterCreate(tx *gorm.DB) error {
	shoppingCart := &ShoppingCart{}
	tx.First(shoppingCart, item.ShoppingCartID)
	shoppingCart.TotalPrice += item.ItemPrice * float64(item.Quantity)
	tx.Save(shoppingCart)
	return nil
}

func (item *ShoppingCartItem) AfterDelete(tx *gorm.DB) error {
	shoppingCart := &ShoppingCart{}
	tx.First(shoppingCart, item.ShoppingCartID)
	shoppingCart.TotalPrice -= item.ItemPrice * float64(item.Quantity)
	tx.Save(shoppingCart)
	return nil
}
