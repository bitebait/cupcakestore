package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type StockRepository interface {
	Create(stock *models.Stock) error
	SumProductStockQuantity(productID uint) (int, error)
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

func (r *stockRepository) SumProductStockQuantity(productID uint) (int, error) {
	var count int64
	result := r.db.Model(&models.Stock{}).Where("product_id = ?", productID).Select("SUM(quantity)").Scan(&count)
	return int(count), result.Error
}
