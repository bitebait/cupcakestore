package models

import "gorm.io/gorm"

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
