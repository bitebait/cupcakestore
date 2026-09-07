package repositories

import (
	"errors"
	"math"

	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	var order models.Order
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var cart models.ShoppingCart
		if profileID == 0 || cartID == 0 {
			return gorm.ErrRecordNotFound
		}
		if err := tx.Where("id = ? AND profile_id = ?", cartID, profileID).First(&cart).Error; err != nil {
			return err
		}
		err := r.applyPreloads(tx).Where("shopping_cart_id = ? AND profile_id = ?", cartID, profileID).First(&order).Error
		if err == nil {
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := lockOpenCart(tx, cartID); err != nil {
			return err
		}
		if err := tx.Preload("Items.Product").First(&cart, cartID).Error; err != nil {
			return err
		}
		if len(cart.Items) == 0 {
			return errors.New("o carrinho está vazio")
		}
		var total float64
		for i := range cart.Items {
			item := &cart.Items[i]
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item.Product, item.ProductID).Error; err != nil {
				return err
			}
			if item.Quantity <= 0 || !item.Product.IsActive || item.Product.Price <= 0 || math.IsNaN(item.Product.Price) || math.IsInf(item.Product.Price, 0) {
				return errors.New("produto ou quantidade inválida no carrinho")
			}
			item.ItemPrice = item.Product.Price
			if err := tx.Omit("Product").Save(item).Error; err != nil {
				return err
			}
			total += item.ItemPrice * float64(item.Quantity)
			movement := models.Stock{ProfileID: profileID, ProductID: item.ProductID, Quantity: item.Quantity, Type: models.StockSaida}
			if err := tx.Create(&movement).Error; err != nil {
				return err
			}
		}
		cart.Total = math.Round(total*100) / 100
		if err := tx.Model(&cart).UpdateColumn("total", cart.Total).Error; err != nil {
			return err
		}
		var store models.StoreConfig
		if err := tx.First(&store).Error; err != nil {
			return err
		}
		order = models.Order{ProfileID: profileID, ShoppingCartID: cartID, Status: models.ActiveStatus, PaymentMethod: models.PixPaymentMethod, IsDelivery: store.DeliveryIsActive, Total: cart.Total, StockReserved: true}
		if order.IsDelivery {
			order.DeliveryPrice = store.DeliveryPrice
		}
		if order.DeliveryPrice < 0 || math.IsNaN(order.DeliveryPrice) || math.IsInf(order.DeliveryPrice, 0) {
			return errors.New("taxa de entrega inválida")
		}
		order.Total = math.Round((cart.Total+order.DeliveryPrice)*100) / 100
		if math.IsNaN(order.Total) || math.IsInf(order.Total, 0) {
			return errors.New("valor total do pedido inválido")
		}
		if err := tx.Omit(clause.Associations).Create(&order).Error; err != nil {
			return err
		}
		if err := tx.Model(&cart).UpdateColumn("order_id", order.ID).Error; err != nil {
			return err
		}
		// Preserve the delivery information as it was when the order was placed.
		var profile models.Profile
		if err := tx.Preload("User").First(&profile, profileID).Error; err != nil {
			return err
		}
		detail := models.OrderDeliveryDetail{
			OrderID: order.ID, UserFirstName: profile.FirstName, UserLastName: profile.LastName,
			UserEmail: profile.User.Email, UserAddress: profile.Address, UserCity: profile.City,
			UserState: profile.State, UserPostalCode: profile.PostalCode, UserPhoneNumber: profile.PhoneNumber,
			StoreEmail: store.PhysicalStoreEmail, StoreAddress: store.PhysicalStoreAddress,
			StoreCity: store.PhysicalStoreCity, StoreState: store.PhysicalStoreState,
			StorePostalCode: store.PhysicalStorePostalCode, StorePhoneNumber: store.PhysicalStorePhoneNumber,
		}
		if err := tx.Create(&detail).Error; err != nil {
			return err
		}
		if err := tx.Model(&order).UpdateColumn("delivery_detail_d", detail.ID).Error; err != nil {
			return err
		}
		return r.applyPreloads(tx).First(&order, order.ID).Error
	})
	return order, err
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
	if order == nil || order.ID == 0 {
		return errors.New("pedido inválido")
	}
	if err := order.Validate(); err != nil {
		return err
	}
	if order.Status == models.CancelledStatus {
		return r.Cancel(order.ID)
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&existing, order.ID).Error; err != nil {
			return err
		}
		if existing.ProfileID != order.ProfileID || existing.ShoppingCartID != order.ShoppingCartID {
			return errors.New("o perfil e o carrinho do pedido não podem ser alterados")
		}
		if existing.Status == models.CancelledStatus || existing.Status == models.DeliveredStatusDelivered {
			return errors.New("o pedido já foi finalizado")
		}
		if !existing.CanTransitionTo(order.Status) {
			return errors.New("transição de status do pedido inválida")
		}
		// Item prices and ownership are immutable; the delivery choice may change
		// until payment, with its amount calculated from the store configuration.
		updates := map[string]interface{}{
			"status": order.Status, "payment_method": order.PaymentMethod,
			"pix_qr": order.PixQR, "pix_string": order.PixString,
			"pix_transaction_id": order.PixTransactionID, "pix_url": order.PixURL,
		}
		if existing.IsActiveOrAwaitingPayment() && (order.Status == models.AwaitingPaymentStatus || order.Status == models.ProcessingStatus) {
			var store models.StoreConfig
			if err := tx.First(&store).Error; err != nil {
				return err
			}
			var cart models.ShoppingCart
			if err := tx.First(&cart, existing.ShoppingCartID).Error; err != nil {
				return err
			}
			isDelivery := order.IsDelivery && store.DeliveryIsActive
			deliveryPrice := 0.0
			if isDelivery {
				deliveryPrice = store.DeliveryPrice
			}
			total := math.Round((cart.Total+deliveryPrice)*100) / 100
			if deliveryPrice < 0 || math.IsNaN(total) || math.IsInf(total, 0) || total != order.Total {
				return errors.New("o valor do pedido foi alterado; tente novamente")
			}
			updates["is_delivery"], updates["delivery_price"], updates["total"] = isDelivery, deliveryPrice, total
		}
		result := tx.Model(&models.Order{}).Where("id = ? AND status = ?", existing.ID, existing.Status).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errors.New("o pedido foi alterado; tente novamente")
		}
		return r.applyPreloads(tx).First(order, existing.ID).Error
	})
}

func (r *orderRepository) Cancel(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, id).Error; err != nil {
			return err
		}
		if order.Status == models.CancelledStatus {
			return nil
		}
		if order.Status == models.DeliveredStatusSent || order.Status == models.DeliveredStatusDelivered {
			return errors.New("não é possível cancelar um pedido enviado ou entregue")
		}
		result := tx.Model(&models.Order{}).Where("id = ? AND status = ?", id, order.Status).UpdateColumns(map[string]interface{}{"status": models.CancelledStatus, "stock_reserved": false})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errors.New("o pedido foi alterado; tente novamente")
		}
		if !order.StockReserved {
			return nil
		}
		var items []models.ShoppingCartItem
		if err := tx.Where("shopping_cart_id = ?", order.ShoppingCartID).Find(&items).Error; err != nil {
			return err
		}
		for _, item := range items {
			movement := models.Stock{ProfileID: order.ProfileID, ProductID: item.ProductID, Quantity: item.Quantity, Type: models.StockEntrada}
			if err := tx.Create(&movement).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
