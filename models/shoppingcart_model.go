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

type PaymentMethod string
type ShoppingCartStatus string

const (
	CashPaymentMethod PaymentMethod = "Dinheiro"
	PixPaymentMethod  PaymentMethod = "Pix"

	ActiveStatus             ShoppingCartStatus = "Em Aberto"
	AwaitingPaymentStatus    ShoppingCartStatus = "Processando Pagamento"
	PaymentApprovedStatus    ShoppingCartStatus = "Pagamento Aprovado"
	ProcessingStatus         ShoppingCartStatus = "Preparando Pedido"
	DeliveredStatusAwaiting  ShoppingCartStatus = "Aguardando Envio"
	DeliveredStatusSent      ShoppingCartStatus = "Enviado"
	DeliveredStatusDelivered ShoppingCartStatus = "Entregue"
	CancelledStatus          ShoppingCartStatus = "Cancelado"
)

type ShoppingCart struct {
	gorm.Model
	ProfileID        uint               `gorm:"not null" validate:"required"`
	Profile          Profile            `validate:"-"`
	Items            []ShoppingCartItem `gorm:"foreignKey:ShoppingCartID"`
	Total            float64            `gorm:"default:0;trigger:false"`
	PixQR            string             `gorm:"default:''"`
	PixString        string             `gorm:"default:''"`
	PixTransactionID string             `gorm:"default:''"`
	PixURL           string             `gorm:"default:''"`
	Status           ShoppingCartStatus `gorm:"default:'Em Aberto'"`
	PaymentMethod    PaymentMethod      `gorm:"default:'Pix'"`
	IsDelivery       bool               `gorm:"not null;default:true"`
	DeliveryPrice    float64            `gorm:"default:0"`
}

func (c *ShoppingCart) BeforeSave(tx *gorm.DB) (err error) {
	if c.IsDelivery {
		var storeConfig StoreConfig
		if err = tx.First(&storeConfig).Error; err == nil {
			c.DeliveryPrice = storeConfig.DeliveryPrice
		}
	}
	return err
}

func (c *ShoppingCart) updateTotal() {
	var subtotal float64
	for _, item := range c.Items {
		subtotal += item.ItemPrice * float64(item.Quantity)
	}

	c.Total = subtotal
}
