package models

import (
	"errors"
	"time"

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
	ProfileID            uint               `gorm:"not null" validate:"required"`
	Profile              Profile            `validate:"-"`
	ShoppingCartID       uint               `gorm:"not null" validate:"required"`
	ShoppingCart         ShoppingCart       `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Status               ShoppingCartStatus `gorm:"default:'Em Aberto'"`
	PaymentMethod        PaymentMethod      `gorm:"default:'Pix'"`
	PixQR                string             `gorm:"default:''"`
	PixString            string             `gorm:"default:''"`
	PixTransactionID     string             `gorm:"default:''"`
	PixURL               string             `gorm:"default:''"`
	PaymentConfirmedAt   *time.Time
	PaymentConfirmedByID *uint
	IsDelivery           bool                `gorm:"not null"`
	DeliveryPrice        float64             `gorm:"default:0"`
	DeliveryDetailD      uint                `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	DeliveryDetail       OrderDeliveryDetail `validate:"-"`
	StockReserved        bool                `gorm:"not null;default:false"`
	Total                float64             `gorm:"default:0"`
}

func (o *Order) IsCurrentUserOrder(profileID uint) bool {
	return o.ProfileID == profileID
}

func (o *Order) CanRedirectToPixPayment() bool {
	return o.Status == AwaitingPaymentStatus && o.PaymentMethod == PixPaymentMethod
}

func (o *Order) CanProceedToPayment() bool {
	return o.ShoppingCart.Total > 0 && o.IsActiveOrAwaitingPayment()
}

func (o *Order) CanProceedToCheckout() bool {
	return o.ShoppingCart.Total > 0 && o.IsActiveOrAwaitingPayment()
}

// Validate checks the persisted state on both creation and updates.
func (o *Order) Validate() error {
	if err := o.validateStatus(); err != nil {
		return err
	}
	return o.validatePaymentMethod()
}

func (o *Order) IsActiveOrAwaitingPayment() bool {
	return o.Status == ActiveStatus || o.Status == AwaitingPaymentStatus
}

func (o *Order) validateStatus() error {
	switch o.Status {
	case ActiveStatus, AwaitingPaymentStatus, PaymentApprovedStatus, ProcessingStatus, DeliveredStatusAwaiting, DeliveredStatusSent, DeliveredStatusDelivered, CancelledStatus:
		return nil
	default:
		return errors.New("invalid shopping cart status")
	}
}

func (o *Order) validatePaymentMethod() error {
	switch o.PaymentMethod {
	case CashPaymentMethod, PixPaymentMethod:
		return nil
	default:
		return errors.New("invalid payment method")
	}
}

// CanTransitionTo prevents stale requests from reopening or moving orders backwards.
func (o *Order) CanTransitionTo(status ShoppingCartStatus) bool {
	if o.Status == status {
		return true
	}
	switch o.Status {
	case ActiveStatus:
		return status == AwaitingPaymentStatus || status == ProcessingStatus || status == CancelledStatus
	case AwaitingPaymentStatus:
		return status == PaymentApprovedStatus || status == ProcessingStatus || status == CancelledStatus
	case PaymentApprovedStatus:
		return status == ProcessingStatus || status == CancelledStatus
	case ProcessingStatus:
		return status == DeliveredStatusAwaiting || status == CancelledStatus || (!o.IsDelivery && status == DeliveredStatusDelivered)
	case DeliveredStatusAwaiting:
		return (o.IsDelivery && status == DeliveredStatusSent) || status == CancelledStatus || (!o.IsDelivery && status == DeliveredStatusDelivered)
	case DeliveredStatusSent:
		return status == DeliveredStatusDelivered
	default:
		return false
	}
}

// CanUpdateStatus excludes checkout and bank confirmation from routine fulfillment.
func (o *Order) CanUpdateStatus(status ShoppingCartStatus) bool {
	if o.Status == status {
		return true
	}
	if o.Status == ActiveStatus || o.Status == AwaitingPaymentStatus {
		return status == CancelledStatus
	}
	return o.CanTransitionTo(status)
}

func (o Order) CanCustomerCancel() bool {
	if o.PaymentConfirmedAt != nil || (o.PaymentMethod == PixPaymentMethod && o.Status != ActiveStatus && o.Status != AwaitingPaymentStatus) {
		return false
	}
	switch o.Status {
	case ActiveStatus, AwaitingPaymentStatus, ProcessingStatus, DeliveredStatusAwaiting:
		return true
	}
	return false
}

// AvailableTransitions exposes the same rules used by persistence to order management.
// A value receiver also makes the choices available to server-rendered templates.
func (o Order) AvailableTransitions() []ShoppingCartStatus {
	var statuses []ShoppingCartStatus
	for _, status := range []ShoppingCartStatus{
		AwaitingPaymentStatus, PaymentApprovedStatus, ProcessingStatus,
		DeliveredStatusAwaiting, DeliveredStatusSent, DeliveredStatusDelivered, CancelledStatus,
	} {
		if status != o.Status && o.CanUpdateStatus(status) {
			statuses = append(statuses, status)
		}
	}
	return statuses
}
