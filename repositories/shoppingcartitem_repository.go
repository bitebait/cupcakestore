package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ShoppingCartItemRepository interface {
	Create(item *models.ShoppingCartItem) error
	Update(item *models.ShoppingCartItem) error
	FindById(id uint) (models.ShoppingCartItem, error)
	Delete(cartID, productID uint) error
}

type shoppingCartItemRepository struct {
	db *gorm.DB
}

func NewShoppingCartItemRepository(database *gorm.DB) ShoppingCartItemRepository {
	return &shoppingCartItemRepository{
		db: database,
	}
}

func (r *shoppingCartItemRepository) Create(item *models.ShoppingCartItem) error {
	return r.db.Create(item).Error
}

func (r *shoppingCartItemRepository) Update(item *models.ShoppingCartItem) error {
	return r.db.Save(item).Error
}

func (r *shoppingCartItemRepository) FindById(id uint) (models.ShoppingCartItem, error) {
	var cartItem models.ShoppingCartItem
	err := r.db.First(&cartItem, id).Error
	return cartItem, err
}

func (r *shoppingCartItemRepository) Delete(cartID, productID uint) error {
	item := &models.ShoppingCartItem{}
	err := r.db.
		Where("shopping_cart_id = ? AND product_id = ?", cartID, productID).
		First(item).Error
	if err != nil {
		return err
	}

	err = r.db.Delete(item).Error
	return err
}
