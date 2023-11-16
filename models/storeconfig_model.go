package models

import (
	"gorm.io/gorm"
)

type pixType string

const (
	PixTypeEmail     pixType = "email"
	PixTypePhone     pixType = "celular"
	PixTypeRandomKey pixType = "aleatoria"
	PixTypeCPF       pixType = "cpf"
	PixTypeCNPJ      pixType = "cnpj"
)

type StoreConfig struct {
	gorm.Model
	DeliveryValue            float64 `gorm:"default:0"`
	DeliveryIsActive         bool    `gorm:"not null;default:true"`
	PhysicalStoreAddress     string  `gorm:"type:varchar(200);default:''"`
	PhysicalStoreCity        string  `gorm:"type:varchar(100);default:''"`
	PhysicalStoreState       string  `gorm:"type:varchar(100);default:''"`
	PhysicalStorePostalCode  string  `gorm:"type:varchar(20);default:''"`
	PhysicalStorePhoneNumber string  `gorm:"type:varchar(20);default:''"`
	PaymentCashIsActive      bool    `gorm:"not null;default:true"`
	PaymentPixIsActive       bool    `gorm:"not null;default:true"`
	PixKey                   string  `gorm:"default:''"`
	PixType                  pixType
}
