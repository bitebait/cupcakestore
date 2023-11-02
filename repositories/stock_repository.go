package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type StockRepository interface {
	Create(stock *models.Stock) error
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(database *gorm.DB) StockRepository {
	return &stockRepository{
		db: database,
	}
}

func (r *stockRepository) Create(stock *models.Stock) error {
	return r.db.Create(stock).Error
}
