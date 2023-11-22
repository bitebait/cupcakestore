package models

import (
	"gorm.io/gorm"
)

type OrderFilter struct {
	Order      *Order
	Pagination *Pagination
}

func NewOrderFilter(profileID uint, page, limit int) *OrderFilter {
	order := &Order{
		ProfileID: profileID,
	}
	pagination := NewPagination(page, limit)
	return &OrderFilter{
		Order:      order,
		Pagination: pagination,
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

type Order struct {
	gorm.Model
	ProfileID        uint               `gorm:"not null" validate:"required"`
	Profile          Profile            `validate:"-"`
	ShoppingCart     ShoppingCart       `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	ShoppingCartID   uint               `gorm:"not null" validate:"required"`
	Status           ShoppingCartStatus `gorm:"default:'Em Aberto'"`
	PaymentMethod    PaymentMethod      `gorm:"default:'Pix'"`
	PixQR            string             `gorm:"default:''"`
	PixString        string             `gorm:"default:''"`
	PixTransactionID string             `gorm:"default:''"`
	PixURL           string             `gorm:"default:''"`
	IsDelivery       bool               `gorm:"not null;default:true"`
	DeliveryPrice    float64            `gorm:"default:0"`
	Total            float64            `gorm:"default:0;trigger:false"`
}

func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
	if o.IsDelivery {
		var storeConfig StoreConfig
		if err = tx.First(&storeConfig).Error; err == nil {
			o.DeliveryPrice = storeConfig.DeliveryPrice
		}
	}
	o.Total = o.ShoppingCart.Total + o.DeliveryPrice
	return err
}

func (o *Order) IsActiveOrAwaitingPayment() bool {
	return o.Status == ActiveStatus || o.Status == AwaitingPaymentStatus
}
