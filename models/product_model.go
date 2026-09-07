package models

import (
	"errors"
	"math"
	"strings"

	"gorm.io/gorm"
)

type ProductFilter struct {
	Product    *Product
	Pagination *Pagination
}

func NewProductFilter(query string, page, limit int) *ProductFilter {
	product := &Product{
		Name: query,
	}
	pagination := NewPagination(page, limit)
	return &ProductFilter{
		Product:    product,
		Pagination: pagination,
	}
}

type Product struct {
	gorm.Model
	Name         string  `gorm:"not null;type:varchar(60)"`
	Description  string  `gorm:"not null;type:varchar(200)"`
	Price        float64 `gorm:"not null"`
	Ingredients  string  `gorm:"not null;type:varchar(300)"`
	Image        string
	Thumbnail    string
	CurrentStock int
	IsActive     bool `gorm:"default:true"`
}

func (p *Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("o nome do produto deve ser informado")
	}
	if p.Price <= 0 || math.IsNaN(p.Price) || math.IsInf(p.Price*100, 0) {
		return errors.New("o preço do produto deve ser maior que zero")
	}
	return nil
}

func (p *Product) BeforeSave(tx *gorm.DB) error {
	return p.Validate()
}
