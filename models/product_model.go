package models

import (
	"github.com/go-playground/validator/v10"
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
	Name         string  `gorm:"not null,type:varchar(60)"`
	Description  string  `gorm:"not null,type:varchar(200)"`
	Price        float64 `gorm:"not null"`
	Ingredients  string  `gorm:"not null,type:varchar(300)"`
	Image        string
	Thumbnail    string
	CurrentStock int
	IsActive     bool `gorm:"default:true"`
}

func (p *Product) Validate() error {
	v := validator.New()
	return v.Struct(p)
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return nil
}

func (p *Product) AfterDelete(tx *gorm.DB) (err error) {
	if err = tx.Where("product_id = ?", p.ID).Delete(&Stock{}).Error; err != nil {
		return err
	}
	return nil
}
