package repositories

import (
	"errors"
	"math"
	"time"

	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log/slog"
)

type OrderRepository interface {
	FindById(id uint) (models.Order, error)
	FindByCartId(cartID uint) (models.Order, error)
	FindOrCreate(profileID, cartID uint) (models.Order, error)
	FindAll(filter *models.OrderFilter) []models.Order
	FindAllByUser(filter *models.OrderFilter) []models.Order
	Update(order *models.Order) error
	UpdatePayment(order *models.Order, configVersion time.Time) error
	Cancel(id uint) error
	CancelForCustomer(id, profileID uint) error
	ConfirmPayment(id, staffProfileID uint) error
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
		slog.Error("OrderRepository.FindOrCreateById", "error", err)
		return order, err
	}

	return order, nil
}

func (r *orderRepository) FindByCartId(cartID uint) (models.Order, error) {
	var order models.Order
	err := r.applyPreloads(r.db).Where("shopping_cart_id = ?", cartID).First(&order).Error

	if err != nil {
		slog.Error("OrderRepository.FindByCartId", "error", err)
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
		var store models.StoreConfig
		if err := tx.First(&store).Error; err != nil {
			return err
		}
		var profile models.Profile
		if err := tx.Preload("User").First(&profile, profileID).Error; err != nil {
			return err
		}
		if err := store.ValidateCheckout(profile); err != nil {
			return err
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
		order = models.Order{ProfileID: profileID, ShoppingCartID: cartID, Status: models.ActiveStatus, PaymentMethod: models.PixPaymentMethod, IsDelivery: store.DeliveryIsActive && profile.HasDeliveryAddress(), Total: cart.Total, StockReserved: true}
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
		// Checkout begins a draft snapshot. Payment validates and records the
		// final contact/address so the customer can correct details beforehand.
		detail := models.NewOrderDeliveryDetail(order.ID, profile, store)
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
		slog.Error("OrderRepository.FindAll", "error", err)
		return nil
	}
	filter.Pagination.Total = total

	var orders []models.Order
	if err := query.
		Offset(offset).
		Limit(filter.Pagination.Limit).
		Order("created_at desc,updated_at desc").
		Find(&orders).Error; err != nil {
		slog.Error("OrderRepository.FindAll", "error", err)
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
		slog.Error("OrderRepository.FindAllByUser", "error", err)
		return nil
	}
	filter.Pagination.Total = total

	var orders []models.Order
	if err := query.
		Offset(offset).
		Limit(filter.Pagination.Limit).
		Order("created_at desc,updated_at desc").
		Find(&orders).Error; err != nil {
		slog.Error("OrderRepository.FindAllByUser", "error", err)
		return nil
	}

	return orders
}

// Update changes only fulfillment status, preserving the agreed payment and delivery terms.
func (r *orderRepository) Update(order *models.Order) error {
	return r.update(order, false, time.Time{})
}

// UpdatePayment validates a payment choice against current store settings before saving it.
func (r *orderRepository) UpdatePayment(order *models.Order, configVersion time.Time) error {
	return r.update(order, true, configVersion)
}

func (r *orderRepository) update(order *models.Order, paymentChoice bool, configVersion time.Time) error {
	if order == nil || order.ID == 0 {
		return errors.New("pedido inválido")
	}
	if err := order.Validate(); err != nil {
		return err
	}
	if order.Status == models.CancelledStatus && !paymentChoice {
		return r.Cancel(order.ID)
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&existing, order.ID).Error; err != nil {
			return err
		}
		if !existing.UpdatedAt.Equal(order.UpdatedAt) {
			return errors.New("o pedido foi alterado; recarregue a página e tente novamente")
		}
		if existing.ProfileID != order.ProfileID || existing.ShoppingCartID != order.ShoppingCartID {
			return errors.New("o perfil e o carrinho do pedido não podem ser alterados")
		}
		if existing.Status == models.CancelledStatus || existing.Status == models.DeliveredStatusDelivered {
			return errors.New("o pedido já foi finalizado")
		}
		if (paymentChoice && !existing.CanTransitionTo(order.Status)) || (!paymentChoice && !existing.CanUpdateStatus(order.Status)) {
			return errors.New("transição de status do pedido inválida")
		}
		updates := map[string]any{"status": order.Status}
		if paymentChoice {
			if existing.Status != models.ActiveStatus ||
				(order.PaymentMethod == models.CashPaymentMethod && order.Status != models.ProcessingStatus) ||
				(order.PaymentMethod == models.PixPaymentMethod && order.Status != models.AwaitingPaymentStatus) {
				return errors.New("o pedido não aceita uma nova escolha de pagamento")
			}
			updates["payment_method"] = order.PaymentMethod
			updates["pix_qr"], updates["pix_string"] = order.PixQR, order.PixString
			updates["pix_transaction_id"], updates["pix_url"] = order.PixTransactionID, order.PixURL
			var store models.StoreConfig
			if err := tx.First(&store).Error; err != nil {
				return err
			}
			if !store.UpdatedAt.Equal(configVersion) {
				return errors.New("a configuração da loja foi alterada; recarregue o pagamento")
			}
			if (order.PaymentMethod == models.CashPaymentMethod && !store.PaymentCashIsActive) || (order.PaymentMethod == models.PixPaymentMethod && !store.IsPixAvailable()) {
				return models.ErrPaymentUnavailable
			}
			var cart models.ShoppingCart
			if err := tx.First(&cart, existing.ShoppingCartID).Error; err != nil {
				return err
			}
			var profile models.Profile
			if err := tx.Preload("User").First(&profile, existing.ProfileID).Error; err != nil {
				return err
			}
			if err := store.ValidateFulfillment(profile, order.IsDelivery); err != nil {
				return err
			}
			isDelivery := order.IsDelivery
			deliveryPrice := 0.0
			if isDelivery {
				deliveryPrice = store.DeliveryPrice
			}
			total := math.Round((cart.Total+deliveryPrice)*100) / 100
			if deliveryPrice < 0 || math.IsNaN(total) || math.IsInf(total, 0) || total != order.Total {
				return errors.New("o valor do pedido foi alterado; tente novamente")
			}
			updates["is_delivery"], updates["delivery_price"], updates["total"] = isDelivery, deliveryPrice, total
			var previousDetail models.OrderDeliveryDetail
			if err := tx.Where("order_id = ?", existing.ID).First(&previousDetail).Error; err != nil {
				return err
			}
			detail := models.NewOrderDeliveryDetail(existing.ID, profile, store)
			detail.Model = previousDetail.Model
			if err := tx.Select("*").Omit("id", "created_at", "deleted_at").Updates(&detail).Error; err != nil {
				return err
			}
		}
		result := tx.Model(&models.Order{}).Where("id = ? AND status = ? AND updated_at = ?", existing.ID, existing.Status, order.UpdatedAt).Updates(updates)
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
	return r.cancel(id, 0)
}

func (r *orderRepository) CancelForCustomer(id, profileID uint) error {
	if profileID == 0 {
		return errors.New("cliente inválido")
	}
	return r.cancel(id, profileID)
}

func (r *orderRepository) cancel(id, customerProfileID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, id).Error; err != nil {
			return err
		}
		if customerProfileID != 0 {
			if order.ProfileID != customerProfileID {
				return errors.New("pedido não pertence ao cliente")
			}
			if order.Status != models.CancelledStatus && !order.CanCustomerCancel() {
				return errors.New("entre em contato com a loja para solicitar o cancelamento deste pedido")
			}
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

func (r *orderRepository) ConfirmPayment(id, staffProfileID uint) error {
	if id == 0 || staffProfileID == 0 {
		return errors.New("pedido e responsável devem ser informados")
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		var staff models.Profile
		if err := tx.Preload("User").First(&staff, staffProfileID).Error; err != nil {
			return err
		}
		if !staff.User.IsActive || !staff.User.IsStaff {
			return errors.New("somente a equipe ativa pode confirmar pagamentos")
		}
		var order models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, id).Error; err != nil {
			return err
		}
		if order.PaymentMethod != models.PixPaymentMethod {
			return errors.New("este pedido não possui um pagamento Pix")
		}
		if order.PaymentConfirmedAt != nil {
			return nil
		}
		if order.Status != models.AwaitingPaymentStatus || order.PixString == "" {
			return errors.New("este pedido não aguarda confirmação de Pix")
		}
		now := time.Now()
		result := tx.Model(&models.Order{}).Where("id = ? AND status = ? AND payment_confirmed_at IS NULL", id, models.AwaitingPaymentStatus).
			UpdateColumns(map[string]any{"status": models.PaymentApprovedStatus, "payment_confirmed_at": now, "payment_confirmed_by_id": staffProfileID, "updated_at": now})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errors.New("o pedido foi alterado; recarregue a página")
		}
		return nil
	})
}
