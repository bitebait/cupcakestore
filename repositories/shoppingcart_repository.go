package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ShoppingCartRepository interface {
	FindByUserId(id uint) (models.ShoppingCart, error)
}

type shoppingCartRepository struct {
	db *gorm.DB
}

func NewShoppingCartRepository(database *gorm.DB) ShoppingCartRepository {
	return &shoppingCartRepository{
		db: database,
	}
}

func (r shoppingCartRepository) FindByUserId(id uint) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	err := r.db.Where("user_id = ?", id).Preload("Profile").Preload("ShoppingCartItem").First(&cart).Error
	return cart, err
}
