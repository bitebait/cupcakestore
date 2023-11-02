package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Stock struct {
	gorm.Model
	ProductID uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
}

func (s *Stock) Validate() error {
	v := validator.New()
	return v.Struct(s)
}

func (s *Stock) BeforeCreate(tx *gorm.DB) error {
	if err := s.Validate(); err != nil {
		return err
	}
	return nil
}
