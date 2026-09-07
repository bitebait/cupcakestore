package models

import (
	"errors"
	"gorm.io/gorm"
	"math"
	"strings"
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
	PixReceiverName          string `gorm:"type:varchar(100);default:''"`
}

func (s *StoreConfig) BeforeSave(tx *gorm.DB) error {
	s.PhysicalStoreAddress = strings.TrimSpace(s.PhysicalStoreAddress)
	s.PhysicalStoreCity = strings.TrimSpace(s.PhysicalStoreCity)
	s.PhysicalStoreState = strings.ToUpper(strings.TrimSpace(s.PhysicalStoreState))
	s.PixReceiverName = strings.TrimSpace(s.PixReceiverName)
	if err := s.Validate(); err != nil {
		return err
	}
	if s.PaymentPixIsActive {
		s.PixKey, _ = NormalizePixKey(string(s.PixKeyType), s.PixKey)
	}
	return nil
}

func (s StoreConfig) Validate() error {
	if s.DeliveryPrice < 0 || math.IsNaN(s.DeliveryPrice) || math.IsInf(s.DeliveryPrice, 0) {
		return errors.New("a taxa de entrega deve ser um valor não negativo")
	}
	if s.PaymentPixIsActive {
		_, err := NormalizePixKey(string(s.PixKeyType), s.PixKey)
		if err != nil {
			return err
		}
		if _, _, err := NormalizePixMerchant(s.PixReceiverName, s.PhysicalStoreCity); err != nil {
			return err
		}
	}
	if !s.PaymentPixIsActive && s.PixKeyType == "" {
		return nil
	}
	return s.validatePixType()
}

func (s StoreConfig) IsPickupAvailable() bool {
	return validAddress(s.PhysicalStoreAddress, s.PhysicalStoreCity, s.PhysicalStoreState)
}

func (s StoreConfig) IsPixAvailable() bool {
	if !s.PaymentPixIsActive {
		return false
	}
	if _, err := NormalizePixKey(string(s.PixKeyType), s.PixKey); err != nil {
		return false
	}
	_, _, err := NormalizePixMerchant(s.PixReceiverName, s.PhysicalStoreCity)
	return err == nil
}

func (s *StoreConfig) validatePixType() error {
	switch s.PixKeyType {
	case PixTypeEmail, PixTypePhone, PixTypeRandomKey, PixTypeCPF, PixTypeCNPJ:
		return nil
	default:
		return errors.New("invalid PixKeyType")
	}
}
