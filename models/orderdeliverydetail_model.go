package models

import (
	"strings"

	"gorm.io/gorm"
)

type OrderDeliveryDetail struct {
	gorm.Model
	OrderID          uint
	UserFirstName    string `gorm:"type:varchar(100)"`
	UserLastName     string `gorm:"type:varchar(100)"`
	UserEmail        string `gorm:"type:varchar(100)"`
	UserAddress      string `gorm:"type:varchar(100)"`
	UserCity         string `gorm:"type:varchar(100)"`
	UserState        string `gorm:"type:varchar(100)"`
	UserPostalCode   string `gorm:"type:varchar(20)"`
	UserPhoneNumber  string `gorm:"type:varchar(20)"`
	StoreEmail       string `gorm:"type:varchar(100)"`
	StoreAddress     string `gorm:"type:varchar(200)"`
	StoreCity        string `gorm:"type:varchar(100)"`
	StoreState       string `gorm:"type:varchar(100)"`
	StorePostalCode  string `gorm:"type:varchar(20)"`
	StorePhoneNumber string `gorm:"type:varchar(20)"`
}

func NewOrderDeliveryDetail(orderID uint, profile Profile, store StoreConfig) OrderDeliveryDetail {
	return OrderDeliveryDetail{
		OrderID: orderID, UserFirstName: strings.TrimSpace(profile.FirstName), UserLastName: strings.TrimSpace(profile.LastName),
		UserEmail: profile.User.Email, UserAddress: strings.TrimSpace(profile.Address), UserCity: strings.TrimSpace(profile.City),
		UserState: strings.ToUpper(strings.TrimSpace(profile.State)), UserPostalCode: strings.TrimSpace(profile.PostalCode), UserPhoneNumber: strings.TrimSpace(profile.PhoneNumber),
		StoreEmail: store.PhysicalStoreEmail, StoreAddress: strings.TrimSpace(store.PhysicalStoreAddress),
		StoreCity: strings.TrimSpace(store.PhysicalStoreCity), StoreState: strings.ToUpper(strings.TrimSpace(store.PhysicalStoreState)),
		StorePostalCode: strings.TrimSpace(store.PhysicalStorePostalCode), StorePhoneNumber: strings.TrimSpace(store.PhysicalStorePhoneNumber),
	}
}
