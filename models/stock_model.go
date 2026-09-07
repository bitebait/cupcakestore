package models

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
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

func (s *Stock) Validate() error {
	v := validator.New()
	if err := v.Struct(s); err != nil {
		return err
	}
	if s.Quantity <= 0 {
		return errors.New("quantidade deve ser maior que zero")
	}
	if s.Type != StockEntrada && s.Type != StockSaida {
		return errors.New("tipo de movimentação de estoque inválido")
	}
	return nil
}

// Stock movements are immutable. Their signed quantity and product balance are
// written in the same transaction, using a conditional debit to prevent overselling.
func (s *Stock) BeforeCreate(tx *gorm.DB) error {
	if err := s.Validate(); err != nil {
		return err
	}
	if s.Type == StockSaida {
		s.Quantity = -s.Quantity
	}
	query := tx.Model(&Product{}).Where("id = ?", s.ProductID)
	if s.Type == StockEntrada {
		query = query.Unscoped()
	}
	if s.Quantity < 0 {
		query = query.Where("current_stock >= ?", -s.Quantity)
	}
	result := query.UpdateColumn("current_stock", gorm.Expr("current_stock + ?", s.Quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("produto não encontrado ou estoque insuficiente")
	}
	return nil
}

func (s *Stock) BeforeUpdate(tx *gorm.DB) error {
	return errors.New("movimentações de estoque não podem ser alteradas")
}
