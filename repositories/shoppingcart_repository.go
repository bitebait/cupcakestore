package repositories

import (
	"errors"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ShoppingCartRepository interface {
	FindByUserId(id uint) (models.ShoppingCart, error)
	AddItemToCart(cartItem *models.ShoppingCartItem) error
	UpdateCartItem(item *models.ShoppingCartItem) error
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
	err := r.db.Where("profile_id = ? AND status = ?", id, models.ActiveStatus).Preload("Profile").Preload("Items.Product").First(&cart).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		cart.ProfileID = id
		err = r.db.Create(&cart).Error
		return cart, err
	}
	return cart, err
}

func (r shoppingCartRepository) AddItemToCart(cartItem *models.ShoppingCartItem) error {
	return r.db.Create(cartItem).Error
}

func (r shoppingCartRepository) UpdateCartItem(item *models.ShoppingCartItem) error {
	return r.db.Save(item).Error
}
