package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ShoppingCartRepository interface {
	FindByUserId(id uint) (models.ShoppingCart, error)
	AddItemToCart(cartItem *models.ShoppingCartItem) error
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
	err := r.db.Where("profile_id = ? AND status = 'Em Aberto'", id).Preload("Profile").Preload("ShoppingCartItem").First(&cart).Error
	if err != nil {
		cart.ProfileID = id
		r.db.Create(&cart)
		return cart, nil
	}
	return cart, err
}

func (r shoppingCartRepository) AddItemToCart(cartItem *models.ShoppingCartItem) error {
	return r.db.FirstOrCreate(cartItem).Error
}
