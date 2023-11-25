package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type StockFilter struct {
	Stock      *Stock
	Pagination *Pagination
}

func NewStockFilter(productID uint, page, limit int) *StockFilter {
	return &StockFilter{
		Stock:      &Stock{ProductID: productID},
		Pagination: NewPagination(page, limit),
	}
}

type stockType string

const (
	StockEntrada stockType = "entrada"
	StockSaida   stockType = "saída"
)

type Stock struct {
	gorm.Model
	ProfileID uint      `gorm:"not null" validate:"required"`
	Profile   Profile   `validate:"-"`
	ProductID uint      `gorm:"not null" validate:"required"`
	Product   Product   `validate:"-"`
	Quantity  int       `gorm:"not null" validate:"required"`
	Type      stockType `validate:"required"`
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

func (s *Stock) BeforeSave(tx *gorm.DB) error {
	if err := s.Validate(); err != nil {
		return err
	}

	if s.Type == StockSaida {
		s.Quantity = -s.Quantity
	}

	return nil
}

func (s *Stock) BeforeCreate(tx *gorm.DB) error {
	return s.Validate()
}

func (s *Stock) AfterSave(tx *gorm.DB) error {
	tx.Model(&Product{}).Where("id = ?", s.ProductID).Update("CurrentStock", s.CountStock(tx))
	return nil
}
