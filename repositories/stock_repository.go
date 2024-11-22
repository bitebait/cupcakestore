package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2/log"
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
	if err := r.db.Create(stock).Error; err != nil {
		log.Errorf("StockRepository Create: %s", err.Error())
		return err
	}

	return nil
}

func (r *stockRepository) SumProductStockQuantity(productID uint) (int, error) {
	var totalQuantity int64
	err := r.db.Model(&models.Stock{}).Where("product_id = ?", productID).Select("SUM(quantity)").Scan(&totalQuantity).Error

	if err != nil {
		log.Errorf("StockRepository SumProductStockQuantity: %s", err.Error())
	}

	return int(totalQuantity), err
}

func (r *stockRepository) FindByProductId(filter *models.StockFilter) []models.Stock {
	offset := calculateOffset(filter.Pagination)
	total := countTotalStocks(r.db, filter.Stock.ProductID)
	filter.Pagination.Total = total

	return fetchStocks(r.db, filter, offset)
}

func calculateOffset(pagination *models.Pagination) int {
	return (pagination.Page - 1) * pagination.Limit
}

func countTotalStocks(db *gorm.DB, productID uint) int64 {
	var total int64
	db.Model(&models.Stock{}).Where("product_id = ?", productID).Count(&total)

	return total
}

func fetchStocks(db *gorm.DB, filter *models.StockFilter, offset int) []models.Stock {
	var stocks []models.Stock

	query := db.Model(&models.Stock{}).
		Where("product_id = ?", filter.Stock.ProductID).
		Offset(offset).
		Limit(filter.Pagination.Limit).
		Order("created_at desc").
		Preload("Product").
		Preload("Profile.User")
	query.Find(&stocks)

	return stocks
}
