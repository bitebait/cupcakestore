package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(database *gorm.DB) ProductRepository {
	return &productRepository{
		db: database,
	}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}
