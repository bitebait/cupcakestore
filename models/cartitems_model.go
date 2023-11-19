package models

import "gorm.io/gorm"

type ShoppingCartItem struct {
	gorm.Model
	ProductID      uint    `gorm:"not null" validate:"required"`
	Product        Product `validate:"-"`
	ItemPrice      float64 `gorm:"default:0"`
	Quantity       int     `gorm:"default:1"`
	ShoppingCartID uint    `gorm:"not null"`
}

func (p *ShoppingCartItem) BeforeSave(tx *gorm.DB) error {
	var product Product
	if err := tx.First(&product, p.ProductID).Error; err != nil {
		return err
	}
	p.ItemPrice = product.Price
	return nil
}

func (p *ShoppingCartItem) AfterSave(tx *gorm.DB) error {
	var shoppingCart ShoppingCart
	if err := tx.First(&shoppingCart, p.ShoppingCartID).Error; err != nil {
		return err
	}
	err := tx.Model(&shoppingCart).Association("Items").Find(&shoppingCart.Items)
	if err != nil {
		return err
	}
	shoppingCart.TotalPrice = calculateTotalPrice(shoppingCart.Items)
	tx.Model(&shoppingCart).UpdateColumn("TotalPrice", shoppingCart.TotalPrice)
	return nil
}

func calculateTotalPrice(items []ShoppingCartItem) float64 {
	var totalPrice float64
	for _, cartItem := range items {
		totalPrice += cartItem.ItemPrice * float64(cartItem.Quantity)
	}
	return totalPrice
}
