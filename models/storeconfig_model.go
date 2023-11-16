package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type StoreConfig struct {
	gorm.Model
	PaymentMethods       []PaymentMethod `gorm:"foreignKey:StoreConfigID"`
	PixInfo              PixInformation  `gorm:"foreignKey:StoreConfigID"`
	DeliveryValue        float64         `gorm:"not null" validate:"required,gte=0"`
	PhysicalStoreAddress string          `gorm:"not null" validate:"required"`
}

type PaymentMethod struct {
	gorm.Model
	StoreConfigID uint
	Name          string `gorm:"not null" validate:"required"`
	IsActive      bool   `gorm:"not null;default:true"`
}

type PixInformation struct {
	gorm.Model
	StoreConfigID uint
	Key           string `gorm:"not null" validate:"required"`
	Active        bool   `gorm:"not null;default:true"`
}

func (s *StoreConfig) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

func (s *StoreConfig) BeforeCreate(tx *gorm.DB) error {
	return s.Validate()
}

func (s *StoreConfig) BeforeSave(tx *gorm.DB) error {
	return s.Validate()
}
