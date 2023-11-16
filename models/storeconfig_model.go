package models

import (
	"gorm.io/gorm"
)

type StoreConfig struct {
	gorm.Model
	PaymentMethods       []PaymentMethod `gorm:"foreignKey:StoreConfigID"`
	PixInfo              PixInformation  `gorm:"foreignKey:StoreConfigID"`
	DeliveryValue        float64         `gorm:"not null"`
	PhysicalStoreAddress string          `gorm:"not null"`
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
