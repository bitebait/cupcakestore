package models

import (
	"errors"
	"gorm.io/gorm"
	"math"
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
	DeliveryPrice            float64 `gorm:"default:0"`
	DeliveryIsActive         bool    `gorm:"not null;default:true"`
	PhysicalStoreEmail       string  `gorm:"type:varchar(100);default:''"`
	PhysicalStoreAddress     string  `gorm:"type:varchar(200);default:''"`
	PhysicalStoreCity        string  `gorm:"type:varchar(100);default:''"`
	PhysicalStoreState       string  `gorm:"type:varchar(100);default:''"`
	PhysicalStorePostalCode  string  `gorm:"type:varchar(20);default:''"`
	PhysicalStorePhoneNumber string  `gorm:"type:varchar(20);default:''"`
	PaymentCashIsActive      bool    `gorm:"not null;default:true"`
	PaymentPixIsActive       bool    `gorm:"not null;default:false"`
	PixKey                   string  `gorm:"default:''"`
	PixKeyType               pixType
}

func (s *StoreConfig) BeforeSave(tx *gorm.DB) error {
	if s.DeliveryPrice < 0 || math.IsNaN(s.DeliveryPrice) || math.IsInf(s.DeliveryPrice, 0) {
		return errors.New("a taxa de entrega deve ser um valor não negativo")
	}
	if s.PaymentPixIsActive && s.PixKey == "" {
		return errors.New("informe a chave Pix antes de ativar o pagamento")
	}
	if !s.PaymentPixIsActive && s.PixKeyType == "" {
		return nil
	}
	return s.validatePixType()
}

func (s *StoreConfig) validatePixType() error {
	switch s.PixKeyType {
	case PixTypeEmail, PixTypePhone, PixTypeRandomKey, PixTypeCPF, PixTypeCNPJ:
		return nil
	default:
		return errors.New("invalid PixKeyType")
	}
}
