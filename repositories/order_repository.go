package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindAll(filter *models.ShoppingCartFilter) []models.ShoppingCart
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(database *gorm.DB) OrderRepository {
	return &orderRepository{
		db: database,
	}
}

func (r *orderRepository) FindAll(filter *models.ShoppingCartFilter) []models.ShoppingCart {
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit

	query := r.db.Model(&models.ShoppingCart{})

	query = query.Where("profile_id = ? ", filter.ShoppingCart.ProfileID).Preload("Profile").Preload("Items.Product")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil
	}
	filter.Pagination.Total = total

	var orders []models.ShoppingCart
	query.Offset(offset).Limit(filter.Pagination.Limit).Find(&orders)
	return orders
}
