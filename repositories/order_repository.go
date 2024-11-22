package repositories

import (
	"errors"
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type OrderRepository interface {
	FindById(id uint) (models.Order, error)
	FindByCartId(cartID uint) (models.Order, error)
	FindOrCreate(profileID, cartID uint) (models.Order, error)
	FindAll(filter *models.OrderFilter) []models.Order
	FindAllByUser(filter *models.OrderFilter) []models.Order
	Update(order *models.Order) error
	Cancel(id uint) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(database *gorm.DB) OrderRepository {
	return &orderRepository{
		db: database,
	}
}

func (r *orderRepository) applyPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Profile").
		Preload("DeliveryDetail").
		Preload("ShoppingCart.Items.Product")
}

func (r *orderRepository) FindById(id uint) (models.Order, error) {
	var order models.Order
	err := r.applyPreloads(r.db).Preload("Profile.User").First(&order, id).Error

	if err != nil {
		log.Errorf("OrderRepository.FindOrCreateById: %s", err.Error())
		return order, err
	}

	return order, nil
}

func (r *orderRepository) FindByCartId(cartID uint) (models.Order, error) {
	var order models.Order
	err := r.applyPreloads(r.db).Where("shopping_cart_id = ?", cartID).First(&order).Error

	if err != nil {
		log.Errorf("OrderRepository.FindByCartId: %s", err.Error())
		return order, err
	}

	return order, nil
}

func (r *orderRepository) FindOrCreate(profileID, cartID uint) (models.Order, error) {
	foundOrder, err := r.FindByCartId(cartID)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		var cart models.ShoppingCart

		if err = r.db.First(&cart, cartID).Error; err != nil {
			log.Errorf("OrderRepository.FindOrCreate: %s", err.Error())
			return models.Order{}, err
		}

		if cart.ProfileID != profileID {
			log.Errorf("OrderRepository.FindOrCreate: perfil e carrinho não correspondem")
			return models.Order{}, errors.New("perfil e carrinho não correspondem")
		}

		order := models.Order{
			ProfileID:      profileID,
			ShoppingCart:   cart,
			ShoppingCartID: cart.ID,
			Status:         models.ActiveStatus,
			PaymentMethod:  models.PixPaymentMethod,
		}

		if err = r.db.Create(&order).Error; err != nil {
			log.Errorf("OrderRepository.FindOrCreate: %s", err.Error())
			return models.Order{}, err
		}

		cart.OrderID = order.ID
		if err = r.db.Save(&cart).Error; err != nil {
			log.Errorf("OrderRepository.FindOrCreate: %s", err.Error())
			return models.Order{}, err
		}

		order, err = r.FindById(order.ID)
		if err != nil {
			log.Errorf("OrderRepository.FindOrCreate: %s", err.Error())
			return models.Order{}, err
		}

		return order, nil
	} else if err != nil {
		log.Errorf("OrderRepository.FindOrCreate: %s", err.Error())
		return models.Order{}, err
	}

	return foundOrder, nil
}

func (r *orderRepository) FindAll(filter *models.OrderFilter) []models.Order {
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit
	query := r.applyPreloads(r.db).
		Model(&models.Order{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Errorf("OrderRepository.FindAll: %s", err.Error())
		return nil
	}
	filter.Pagination.Total = total

	var orders []models.Order
	if err := query.
		Offset(offset).
		Limit(filter.Pagination.Limit).
		Order("created_at desc,updated_at desc").
		Find(&orders).Error; err != nil {
		log.Errorf("OrderRepository.FindAll: %s", err.Error())
		return nil
	}

	return orders
}

func (r *orderRepository) FindAllByUser(filter *models.OrderFilter) []models.Order {
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit
	query := r.applyPreloads(r.db).
		Model(&models.Order{}).
		Where("profile_id = ?", filter.Order.ProfileID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Errorf("OrderRepository.FindAllByUser: %s", err.Error())
		return nil
	}
	filter.Pagination.Total = total

	var orders []models.Order
	if err := query.
		Offset(offset).
		Limit(filter.Pagination.Limit).
		Order("created_at desc,updated_at desc").
		Find(&orders).Error; err != nil {
		log.Errorf("OrderRepository.FindAllByUser: %s", err.Error())
		return nil
	}

	return orders
}

func (r *orderRepository) Update(order *models.Order) error {
	if err := r.db.Save(order).Error; err != nil {
		log.Errorf("OrderRepository.Update: %s", err.Error())
		return err
	}

	return nil
}

func (r *orderRepository) Cancel(id uint) error {
	foundOrder, err := r.FindById(id)

	if err != nil {
		log.Errorf("OrderRepository.Cancel: %s", err.Error())
		return err
	}

	foundOrder.Status = models.CancelledStatus
	if err := r.db.Save(&foundOrder).Error; err != nil {
		log.Errorf("OrderRepository.Cancel: %s", err.Error())
		return err
	}

	return nil
}
