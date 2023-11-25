package models

import (
	"errors"
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
	ProfileID        uint                `gorm:"not null" validate:"required"`
	Profile          Profile             `validate:"-"`
	ShoppingCartID   uint                `gorm:"not null" validate:"required"`
	ShoppingCart     ShoppingCart        `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Status           ShoppingCartStatus  `gorm:"default:'Em Aberto'"`
	PaymentMethod    PaymentMethod       `gorm:"default:'Pix'"`
	PixQR            string              `gorm:"default:''"`
	PixString        string              `gorm:"default:''"`
	PixTransactionID string              `gorm:"default:''"`
	PixURL           string              `gorm:"default:''"`
	IsDelivery       bool                `gorm:"not null;default:true"`
	DeliveryPrice    float64             `gorm:"default:0"`
	DeliveryDetailD  uint                `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	DeliveryDetail   OrderDeliveryDetail `validate:"-"`
	Total            float64             `gorm:"default:0"`
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

func (o *Order) AfterCreate(tx *gorm.DB) (err error) {
	if err = tx.Preload("Profile.User").Preload("ShoppingCart.Items").First(&o).Error; err != nil {
		return err
	}
	for _, item := range o.ShoppingCart.Items {
		stock := &Stock{
			ProfileID: o.ProfileID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Type:      StockSaida,
		}
		if err = tx.Model(&Stock{}).Create(stock).Error; err != nil {
			return err
		}
	}
	return nil
}

func (o *Order) AfterUpdate(tx *gorm.DB) (err error) {
	if err = tx.Preload("Profile.User").Preload("ShoppingCart.Items").First(&o).Error; err != nil {
		return err
	}
	if o.Status == CancelledStatus {
		for _, item := range o.ShoppingCart.Items {
			stock := &Stock{
				ProfileID: o.ProfileID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Type:      StockEntrada,
			}
			if err = tx.Model(&Stock{}).Create(stock).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *Order) AfterSave(tx *gorm.DB) (err error) {
	storeConfig := &StoreConfig{}
	if err = tx.First(storeConfig).Error; err != nil {
		return err
	}

	if o.Status == ActiveStatus && o.ShoppingCartID > 0 {
		orderDeliveryDetail := &OrderDeliveryDetail{}
		if err = tx.Where("order_id", o.ID).FirstOrInit(&orderDeliveryDetail).Error; err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			if err = tx.Preload("Profile.User").Preload("ShoppingCart.Items").First(&o).Error; err != nil {
				return err
			}
			orderDeliveryDetail.OrderID = o.ID
			orderDeliveryDetail.UserFirstName = o.Profile.FirstName
			orderDeliveryDetail.UserLastName = o.Profile.LastName
			orderDeliveryDetail.UserEmail = o.Profile.User.Email
			orderDeliveryDetail.UserAddress = o.Profile.Address
			orderDeliveryDetail.UserCity = o.Profile.City
			orderDeliveryDetail.UserState = o.Profile.State
			orderDeliveryDetail.UserPostalCode = o.Profile.PostalCode
			orderDeliveryDetail.UserPhoneNumber = o.Profile.PhoneNumber
			orderDeliveryDetail.StoreEmail = storeConfig.PhysicalStoreEmail
			orderDeliveryDetail.StoreAddress = storeConfig.PhysicalStoreAddress
			orderDeliveryDetail.StoreCity = storeConfig.PhysicalStoreCity
			orderDeliveryDetail.StoreState = storeConfig.PhysicalStoreState
			orderDeliveryDetail.StorePostalCode = storeConfig.PhysicalStorePostalCode
			orderDeliveryDetail.StorePhoneNumber = storeConfig.PhysicalStorePhoneNumber
			tx.Save(orderDeliveryDetail)
			o.DeliveryDetailD = orderDeliveryDetail.ID
		}
	}
	return err
}

func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
	storeConfig := &StoreConfig{}

	if o.IsDelivery {
		if err = tx.First(storeConfig).Error; err == nil {
			o.DeliveryPrice = storeConfig.DeliveryPrice
		}
	} else {
		o.DeliveryPrice = 0
	}

	o.Total = o.ShoppingCart.Total + o.DeliveryPrice
	return err
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if err = o.validateStatus(); err != nil {
		return err
	}

	if err = o.validatePaymentMethod(); err != nil {
		return err
	}
	return nil
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
