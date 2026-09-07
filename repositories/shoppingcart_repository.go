package repositories

import (
	"errors"

	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
	"log/slog"
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
		slog.Error("ShoppingCartRepository FindAll: invalid filter params")
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
		slog.Error("ShoppingCartRepository FindAll", "error", err)
		return nil
	}
	filter.Pagination.Total = total

	var carts []models.ShoppingCart
	if err := query.Offset(offset).Limit(filter.Pagination.Limit).Find(&carts).Error; err != nil {
		slog.Error("ShoppingCartRepository FindAll", "error", err)
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

	return cart, err
}

func (r *shoppingCartRepository) FindOrCreateByUserId(userID uint) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	if userID == 0 {
		return cart, gorm.ErrRecordNotFound
	}
	find := func(tx *gorm.DB) error {
		return tx.Preload("Profile").Preload("Items.Product").
			Where("profile_id = ? AND order_id IS NULL", userID).First(&cart).Error
	}
	if err := find(r.db); err == nil {
		return cart, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return cart, err
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Lock the owner before the second lookup so concurrent requests create
		// only one open cart. This works with SQLite and PostgreSQL.
		result := tx.Model(&models.Profile{}).Where("id = ?", userID).
			UpdateColumn("updated_at", gorm.Expr("updated_at"))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		if err := find(tx); err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		cart = models.ShoppingCart{ProfileID: userID}
		return tx.Create(&cart).Error
	})
	return cart, err
}
