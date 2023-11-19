package models

import (
	"gorm.io/gorm"
)

type paymentMethod string
type paymentStatus string

const (
	CashPaymentMethod paymentMethod = "Dinheiro"
	PixPaymentMethod  paymentMethod = "Pix"
)

const (
	ActiveStatus          paymentStatus = "Em Aberto"
	AwaitingPaymentStatus paymentStatus = "Aguardando Pagamento"
	PaymentApprovedStatus paymentStatus = "Pagamento Aprovado"
	ProcessingStatus      paymentStatus = "Em Processamento"
	ShippedStatus         paymentStatus = "Enviado"
	DeliveredStatus       paymentStatus = "Entregue"
	CancelledStatus       paymentStatus = "Cancelado"
)

type ShoppingCart struct {
	gorm.Model
	ProfileID       uint               `gorm:"not null" validate:"required"`
	Profile         Profile            `validate:"-"`
	Items           []ShoppingCartItem `gorm:"foreignKey:ShoppingCartID"`
	TotalPrice      float64            `gorm:"default:0"`
	PaymentStatus   paymentStatus      `gorm:"default:'Em Aberto'"`
	PaymentMethod   paymentMethod      `gorm:"default:'Pix'"`
	DeliveryAddress string             `gorm:"default:''"`
}
