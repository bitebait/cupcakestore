package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type Stock struct {
	gorm.Model
	ProductID uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
}

func (s *Stock) CountStock(tx *gorm.DB) int {
	var count int64
	result := tx.Model(&Stock{}).Where("product_id = ?", s.ProductID).Select("SUM(quantity)").Scan(&count)
	if result.Error != nil {
		log.Println(result.Error)
	}
	return int(count)
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

func (s *Stock) AfterSave(tx *gorm.DB) (err error) {
	tx.Model(&Product{}).Where("id = ?", s.ProductID).Update("CurrentStock", s.CountStock(tx))
	return
}
