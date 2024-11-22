package repositories

import (
	"errors"
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type ShoppingCartRepository interface {
	FindAll(filter *models.ShoppingCartFilter) []models.ShoppingCart
	FindOrCreateByUserId(id uint) (models.ShoppingCart, error)
	FindOrCreateById(id uint) (models.ShoppingCart, error)
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
	if filter.ShoppingCart.ProfileID <= 0 || filter.Pagination.Page <= 0 || filter.Pagination.Limit <= 0 {
		log.Error("ShoppingCartRepository FindAll: invalid filter params")
		return nil
	}

	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit

	query := r.db.
		Model(&models.ShoppingCart{}).
		Where("profile_id = ?", filter.ShoppingCart.ProfileID).
		Preload("Profile").
		Preload("Items.Product")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Errorf("ShoppingCartRepository FindAll: %s", err.Error())
		return nil
	}
	filter.Pagination.Total = total

	var carts []models.ShoppingCart
	if err := query.Offset(offset).Limit(filter.Pagination.Limit).Find(&carts).Error; err != nil {
		log.Errorf("ShoppingCartRepository FindAll: %s", err.Error())
		return nil
	}

	return carts
}

func (r *shoppingCartRepository) FindOrCreateById(id uint) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	err := r.db.
		Preload("Profile").
		Preload("Items.Product").
		Where("id = ?", id).
		First(&cart).Error

	if err == nil {
		return cart, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("ShoppingCartRepository FindOrCreateById: %s", err.Error())
		return cart, err
	}

	cart.ID = id
	if err = r.db.Create(&cart).Error; err != nil {
		log.Errorf("ShoppingCartRepository FindOrCreateById: failed to create new cart: %s", err.Error())
		return cart, err
	}

	return cart, nil
}

func (r *shoppingCartRepository) FindOrCreateByUserId(userId uint) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	cart.ProfileID = userId

	err := r.db.
		Preload("Profile").
		Preload("Items.Product").
		Where("profile_id = ? AND order_id IS NULL", userId).
		FirstOrCreate(&cart).Error

	if err != nil {
		log.Errorf("ShoppingCartRepository FindOrCreateByUserId: %s", err.Error())
		return models.ShoppingCart{}, err
	}

	return cart, nil
}
