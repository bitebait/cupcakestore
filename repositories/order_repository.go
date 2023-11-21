package repositories

import (
	"errors"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
	"log"
)

type OrderRepository interface {
	FindByCartId(cartID uint) (*models.Order, error)
	FindOrCreate(profileID, cartID uint) (*models.Order, error)
	FindAll(filter *models.OrderFilter) []models.Order
	Update(order *models.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(database *gorm.DB) OrderRepository {
	return &orderRepository{
		db: database,
	}
}

func (r *orderRepository) FindByCartId(cartID uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("shopping_cart_id = ?", cartID).
		Preload("Profile").
		Preload("ShoppingCart.Items.Product").
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindOrCreate(profileID, cartID uint) (*models.Order, error) {
	foundOrder, err := r.FindByCartId(cartID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		var cart models.ShoppingCart
		if err := r.db.Where("id = ?", cartID).First(&cart).Error; err != nil {
			return nil, err
		}

		order := models.Order{
			ProfileID:      profileID,
			ShoppingCart:   cart,
			ShoppingCartID: cart.ID,
		}

		if err := r.db.Create(&order).Error; err != nil {
			return nil, err
		}

		cart.OrderID = order.ID
		if err := r.db.Save(&cart).Error; err != nil {
			return nil, err
		}

		return &order, nil
	} else if err != nil {
		return nil, err
	}

	log.Println(foundOrder)
	return foundOrder, nil
}

func (r *orderRepository) FindAll(filter *models.OrderFilter) []models.Order {
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit

	query := r.db.Model(&models.Order{}).
		Where("profile_id = ?", filter.Order.ProfileID).
		Preload("Profile").
		Preload("ShoppingCart.Items.Product")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil
	}
	filter.Pagination.Total = total

	var orders []models.Order
	if err := query.Offset(offset).Limit(filter.Pagination.Limit).Find(&orders).Error; err != nil {
		return nil
	}
	return orders
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}
