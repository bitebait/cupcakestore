package repositories

import (
	"errors"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ShoppingCartRepository interface {
	FindByUserId(id uint) (models.ShoppingCart, error)
	FindById(id uint) (models.ShoppingCart, error)
	Update(cart *models.ShoppingCart) error
}

type shoppingCartRepository struct {
	db *gorm.DB
}

func NewShoppingCartRepository(database *gorm.DB) ShoppingCartRepository {
	return &shoppingCartRepository{
		db: database,
	}
}

func (r shoppingCartRepository) FindById(id uint) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	err := r.db.Where("id = ?", id, models.ActiveStatus).Preload("Profile").Preload("Items.Product").First(&cart).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		cart.ProfileID = id
		err = r.db.Create(&cart).Error
		return cart, err
	}
	return cart, err
}

func (r shoppingCartRepository) FindByUserId(id uint) (models.ShoppingCart, error) {
	var cart models.ShoppingCart
	cart.ProfileID = id
	err := r.db.Where("profile_id = ? AND status = ?", id, models.ActiveStatus).Preload("Profile").Preload("Items.Product").FirstOrCreate(&cart).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = r.db.Create(&cart).Error
		return cart, err
	}
	return cart, err
}

func (r shoppingCartRepository) Update(cart *models.ShoppingCart) error {
	return r.db.Save(cart).Error
}
