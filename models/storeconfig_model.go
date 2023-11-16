package models

import (
	"gorm.io/gorm"
)

type StoreConfig struct {
	gorm.Model
	PaymentMethods           []PaymentMethod `gorm:"foreignKey:StoreConfigID"`
	PixInfo                  PixInformation  `gorm:"foreignKey:StoreConfigID"`
	DeliveryValue            float64
	DeliveryIsActive         bool   `gorm:"not null;default:true"`
	PhysicalStoreAddress     string `gorm:"type:varchar(200)"`
	PhysicalStoreCity        string `gorm:"type:varchar(100)"`
	PhysicalStoreState       string `gorm:"type:varchar(100)"`
	PhysicalStorePostalCode  string `gorm:"type:varchar(20)"`
	PhysicalStorePhoneNumber string `gorm:"type:varchar(20)"`
}

type PaymentMethod struct {
	gorm.Model
	StoreConfigID uint
	Name          string `gorm:"not null"`
	IsActive      bool   `gorm:"not null;default:true"`
}

type PixInformation struct {
	gorm.Model
	StoreConfigID uint
	Key           string `gorm:"not null"`
	Active        bool   `gorm:"not null;default:true"`
}
