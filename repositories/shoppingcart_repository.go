package repositories

import (
	"errors"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ShoppingCartRepository interface {
	FindAll(filter *models.ShoppingCartFilter) []models.ShoppingCart
	FindByUserId(id uint) (models.ShoppingCart, error)
	FindById(id uint) (models.ShoppingCart, error)
}

type shoppingCartRepository struct {
	db *gorm.DB
}

func NewShoppingCartRepository(database *gorm.DB) ShoppingCartRepository {
	return &shoppingCartRepository{
		db: database,
	}
}

func (r *shoppingCartRepository) FindAll(filter *models.ShoppingCartFilter) []models.ShoppingCart {
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit

	query := r.db.Model(&models.ShoppingCart{}).
		Where("profile_id = ? ", filter.ShoppingCart.ProfileID).
		Preload("Profile").
		Preload("Items.Product")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil
	}
	filter.Pagination.Total = total

	var orders []models.ShoppingCart
	query.Offset(offset).Limit(filter.Pagination.Limit).Find(&orders)
	return orders
}

func (r *shoppingCartRepository) FindById(id uint) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	err := r.db.
		Where("id = ?", id).
		Preload("Profile").
		Preload("Items.Product").
		First(&cart).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		cart.ProfileID = id
		err = r.db.Create(&cart).Error
	}
	return cart, err
}

func (r *shoppingCartRepository) FindByUserId(id uint) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	cart.ProfileID = id
	err := r.db.
		Where("profile_id = ? AND order_id IS NULL", id).
		Preload("Profile").
		Preload("Items.Product").
		FirstOrCreate(&cart).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = r.db.Create(&cart).Error
	}
	return cart, err
}
