package models

import (
	"gorm.io/gorm"
)

type paymentMethod string
type shoppingCartStatus string

const (
	CashPaymentMethod paymentMethod = "Dinheiro"
	PixPaymentMethod  paymentMethod = "Pix"
)

const (
	ActiveStatus          shoppingCartStatus = "Em Aberto"
	AwaitingPaymentStatus shoppingCartStatus = "Aguardando Pagamento"
	PaymentApprovedStatus shoppingCartStatus = "Pagamento Aprovado"
	ProcessingStatus      shoppingCartStatus = "Em Processamento"
	ShippedStatus         shoppingCartStatus = "Enviado"
	DeliveredStatus       shoppingCartStatus = "Entregue"
	CancelledStatus       shoppingCartStatus = "Cancelado"
)

type ShoppingCart struct {
	gorm.Model
	ProfileID       uint               `gorm:"not null" validate:"required"`
	Profile         Profile            `validate:"-"`
	Items           []ShoppingCartItem `gorm:"foreignKey:ShoppingCartID"`
	TotalPrice      float64            `gorm:"default:0"`
	Status          shoppingCartStatus `gorm:"default:'Em Aberto'"`
	PaymentMethod   paymentMethod      `gorm:"default:'Pix'"`
	DeliveryAddress string             `gorm:"default:''"`
}
