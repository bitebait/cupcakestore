package models

import (
	"gorm.io/gorm"
)

type PaymentMethod string
type ShoppingCartStatus string

const (
	CashPaymentMethod PaymentMethod = "Dinheiro"
	PixPaymentMethod  PaymentMethod = "Pix"

	ActiveStatus          ShoppingCartStatus = "Em Aberto"
	AwaitingPaymentStatus ShoppingCartStatus = "Aguardando Pagamento"
	PaymentApprovedStatus ShoppingCartStatus = "Pagamento Aprovado"
	ProcessingStatus      ShoppingCartStatus = "Em Processamento"
	ShippedStatus         ShoppingCartStatus = "Enviado"
	DeliveredStatus       ShoppingCartStatus = "Entregue"
	CancelledStatus       ShoppingCartStatus = "Cancelado"
)

type ShoppingCart struct {
	gorm.Model
	ProfileID     uint               `gorm:"not null" validate:"required"`
	Profile       Profile            `validate:"-"`
	Items         []ShoppingCartItem `gorm:"foreignKey:ShoppingCartID"`
	Total         float64            `gorm:"default:0;trigger:false"`
	Status        ShoppingCartStatus `gorm:"default:'Em Aberto'"`
	PaymentMethod PaymentMethod      `gorm:"default:'Pix'"`
	IsDelivery    bool               `gorm:"not null;default:true"`
}

func (c *ShoppingCart) updateTotal() {
	var subtotal float64
	for _, item := range c.Items {
		subtotal += item.ItemPrice * float64(item.Quantity)
	}
	c.Total = subtotal
}
