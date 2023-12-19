package models

import (
	"gorm.io/gorm"
)

type ShoppingCartFilter struct {
	ShoppingCart *ShoppingCart
	Pagination   *Pagination
}

func NewShoppingCartFilter(profileID uint, page, limit int) *ShoppingCartFilter {
	shoppingCart := &ShoppingCart{
		ProfileID: profileID,
	}
	pagination := NewPagination(page, limit)
	return &ShoppingCartFilter{
		ShoppingCart: shoppingCart,
		Pagination:   pagination,
	}
}

type ShoppingCart struct {
	gorm.Model
	ProfileID uint               `gorm:"not null" validate:"required"`
	Profile   Profile            `validate:"-"`
	Items     []ShoppingCartItem `gorm:"foreignKey:ShoppingCartID;constraint:OnDelete:CASCADE"`
	Total     float64            `gorm:"default:0;trigger:false"`
	OrderID   uint               `gorm:"default:null"`
}

func (c *ShoppingCart) updateTotal() {
	var subtotal float64
	for _, item := range c.Items {
		subtotal += item.ItemPrice * float64(item.Quantity)
	}

	c.Total = subtotal
}

func (c *ShoppingCart) CountItems() int64 {
	var count int64
	for _, item := range c.Items {
		count += int64(item.Quantity)
	}

	return count
}
