package repositories

import (
	"errors"
	"time"

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

// lockOpenCart serializes cart writes and checkout in both SQLite and Postgres.
// The conditional write also rejects changes to a finalized cart.
func lockOpenCart(tx *gorm.DB, cartID uint) error {
	result := tx.Model(&models.ShoppingCart{}).
		Where("id = ? AND order_id IS NULL", cartID).
		UpdateColumn("updated_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("carrinho não encontrado ou já finalizado")
	}
	return nil
}

func (r *shoppingCartItemRepository) Create(item *models.ShoppingCartItem) error {
	if item == nil || item.Quantity <= 0 {
		return errors.New("quantidade deve ser maior que zero")
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := lockOpenCart(tx, item.ShoppingCartID); err != nil {
			return err
		}
		var existing models.ShoppingCartItem
		err := tx.Where("shopping_cart_id = ? AND product_id = ?", item.ShoppingCartID, item.ProductID).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Omit("Product").Create(item).Error
		}
		if err != nil {
			return err
		}
		maxInt := int(^uint(0) >> 1)
		if existing.Quantity < 0 || item.Quantity > maxInt-existing.Quantity {
			return errors.New("quantidade inválida")
		}
		existing.Quantity += item.Quantity
		if err := tx.Omit("Product").Save(&existing).Error; err != nil {
			return err
		}
		*item = existing
		return nil
	})
}

func (r *shoppingCartItemRepository) Update(item *models.ShoppingCartItem) error {
	if item == nil || item.ID == 0 || item.Quantity <= 0 {
		return errors.New("item e quantidade válida devem ser informados")
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := lockOpenCart(tx, item.ShoppingCartID); err != nil {
			return err
		}
		var existing models.ShoppingCartItem
		if err := tx.Where("id = ? AND shopping_cart_id = ? AND product_id = ?", item.ID, item.ShoppingCartID, item.ProductID).First(&existing).Error; err != nil {
			return err
		}
		existing.Quantity = item.Quantity
		if err := tx.Omit("Product").Save(&existing).Error; err != nil {
			return err
		}
		*item = existing
		return nil
	})
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
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := lockOpenCart(tx, cartID); err != nil {
			return err
		}
		var item models.ShoppingCartItem
		if err := tx.Where("shopping_cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error; err != nil {
			return err
		}
		return tx.Delete(&item).Error
	})
}
