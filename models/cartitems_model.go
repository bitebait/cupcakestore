package models

import "gorm.io/gorm"

type ShoppingCartItem struct {
	gorm.Model
	ProductID      uint    `gorm:"not null" validate:"required"`
	Product        Product `validate:"-"`
	ItemPrice      float64
	ShoppingCartID uint `gorm:"not null" validate:"required"`
	Quantity       int
}