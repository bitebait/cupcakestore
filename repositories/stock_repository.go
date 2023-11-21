package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type StockRepository interface {
	Create(stock *models.Stock) error
	SumProductStockQuantity(productID uint) (int, error)
	FindByProductId(filter *models.StockFilter) []models.Stock
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
	var totalQuantity int64
	result := r.db.Model(&models.Stock{}).Where("product_id = ?", productID).Select("SUM(quantity)").Scan(&totalQuantity)
	return int(totalQuantity), result.Error
}

func (r *stockRepository) FindByProductId(filter *models.StockFilter) []models.Stock {
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit

	query := r.db.Model(&models.Stock{}).
		Where("product_id = ?", filter.Stock.ProductID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil
	}
	filter.Pagination.Total = total

	var stocks []models.Stock
	query.Offset(offset).
		Limit(filter.Pagination.Limit).
		Order("created_at desc").
		Preload("Product").
		Preload("Profile.User").
		Find(&stocks)
	return stocks
}
