package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2/log"
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
	if err := r.db.Create(item).Error; err != nil {
		log.Errorf("ShoppingCartItemRepository Create: %s", err.Error())
	}

	return nil
}

func (r *shoppingCartItemRepository) Update(item *models.ShoppingCartItem) error {
	if err := r.db.Save(item).Error; err != nil {
		log.Errorf("ShoppingCartItemRepository Update: %s", err.Error())
	}

	return nil
}

func (r *shoppingCartItemRepository) FindById(id uint) (models.ShoppingCartItem, error) {
	var cartItem models.ShoppingCartItem
	err := r.db.First(&cartItem, id).Error

	if err != nil {
		log.Errorf("ShoppingCartItemRepository FindOrCreateById: %s", err.Error())
	}

	return cartItem, err
}

func (r *shoppingCartItemRepository) Delete(cartID, productID uint) error {
	var item models.ShoppingCartItem

	if err := r.db.Where("shopping_cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error; err != nil {
		log.Errorf("ShoppingCartItemRepository Delete: %s", err.Error())
		return err
	}

	if err := r.db.Delete(&item).Error; err != nil {
		log.Errorf("ShoppingCartItemRepository Delete: %s", err.Error())
		return err

	}

	return nil
}
