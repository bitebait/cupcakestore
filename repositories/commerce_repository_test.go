package repositories

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/bitebait/cupcakestore/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func commerceDB(t *testing.T) (*gorm.DB, models.Profile, models.Product) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "commerce.db")+"?_foreign_keys=on&_busy_timeout=5000"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Product{}, &models.Stock{}, &models.ShoppingCart{}, &models.ShoppingCartItem{}, &models.Order{}, &models.OrderDeliveryDetail{}, &models.StoreConfig{}); err != nil {
		t.Fatal(err)
	}
	user := models.User{Email: "customer@example.com", Password: "test-only"}
	if err := db.Session(&gorm.Session{SkipHooks: true}).Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	profile := models.Profile{UserID: user.ID, FirstName: "Cliente", Address: "Rua original"}
	if err := db.Create(&profile).Error; err != nil {
		t.Fatal(err)
	}
	product := models.Product{Name: "Cupcake", Price: 5.25, IsActive: true}
	if err := db.Create(&product).Error; err != nil {
		t.Fatal(err)
	}
	stock := models.Stock{ProfileID: profile.ID, ProductID: product.ID, Quantity: 10, Type: models.StockEntrada}
	if err := db.Create(&stock).Error; err != nil {
		t.Fatal(err)
	}
	store := models.StoreConfig{DeliveryPrice: 2.50, DeliveryIsActive: true, PixKeyType: models.PixTypeEmail}
	if err := db.Create(&store).Error; err != nil {
		t.Fatal(err)
	}
	return db, profile, product
}

func cartWithItem(t *testing.T, db *gorm.DB, profileID, productID uint, quantity int) models.ShoppingCart {
	t.Helper()
	cart, err := NewShoppingCartRepository(db).FindOrCreateByUserId(profileID)
	if err != nil {
		t.Fatal(err)
	}
	if err := NewShoppingCartItemRepository(db).Create(&models.ShoppingCartItem{ShoppingCartID: cart.ID, ProductID: productID, Quantity: quantity}); err != nil {
		t.Fatal(err)
	}
	return cart
}

func assertStock(t *testing.T, db *gorm.DB, productID uint, want int) {
	t.Helper()
	var product models.Product
	if err := db.Unscoped().First(&product, productID).Error; err != nil {
		t.Fatal(err)
	}
	if product.CurrentStock != want {
		t.Fatalf("stock = %d, want %d", product.CurrentStock, want)
	}
	sum, err := NewStockRepository(db).SumProductStockQuantity(productID)
	if err != nil || sum != want {
		t.Fatalf("ledger = %d, %v; want %d", sum, err, want)
	}
}

func TestCheckoutReservesStockAndCancellationIsIdempotent(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	assertStock(t, db, product.ID, 8)
	if order.Total != 13 || !order.StockReserved || order.DeliveryDetail.UserAddress != profile.Address {
		t.Fatalf("incorrect order snapshot: %+v", order)
	}
	repeated, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil || repeated.ID != order.ID {
		t.Fatalf("repeated checkout = %d, %v", repeated.ID, err)
	}
	assertStock(t, db, product.ID, 8)
	if _, err := repo.FindOrCreate(profile.ID+1, cart.ID); err == nil {
		t.Fatal("another profile accessed an existing order")
	}
	for i := 0; i < 2; i++ {
		if err := repo.Cancel(order.ID); err != nil {
			t.Fatal(err)
		}
		assertStock(t, db, product.ID, 10)
	}
	order.Status = models.ProcessingStatus
	if err := repo.Update(&order); err == nil {
		t.Fatal("cancelled order was reopened")
	}
}

func TestCheckoutRollsBackWhenStockRunsOut(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 8)
	movement := models.Stock{ProfileID: profile.ID, ProductID: product.ID, Quantity: 5, Type: models.StockSaida}
	if err := db.Create(&movement).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := NewOrderRepository(db).FindOrCreate(profile.ID, cart.ID); err == nil {
		t.Fatal("checkout succeeded without stock")
	}
	assertStock(t, db, product.ID, 5)
	var count int64
	db.Model(&models.Order{}).Count(&count)
	if count != 0 {
		t.Fatalf("failed checkout left %d orders", count)
	}
	if err := db.First(&cart, cart.ID).Error; err != nil || cart.OrderID != 0 {
		t.Fatalf("failed checkout finalized cart: %+v, %v", cart, err)
	}
}

func TestCheckoutRollsBackEarlierStockMovements(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	// A failure after reservation must undo the product balance and ledger.
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.StoreConfig{}).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := NewOrderRepository(db).FindOrCreate(profile.ID, cart.ID); err == nil {
		t.Fatal("checkout succeeded without store configuration")
	}
	assertStock(t, db, product.ID, 10)
}

func TestCartRejectsInvalidItemsAndFreezesAfterCheckout(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	repo := NewShoppingCartItemRepository(db)
	for _, quantity := range []int{-1, 0, 11, int(^uint(0) >> 1)} {
		if err := repo.Create(&models.ShoppingCartItem{ShoppingCartID: cart.ID, ProductID: product.ID, Quantity: quantity}); err == nil {
			t.Fatalf("accepted invalid quantity %d", quantity)
		}
	}
	if err := repo.Create(&models.ShoppingCartItem{ShoppingCartID: cart.ID, ProductID: product.ID + 1, Quantity: 1}); err == nil {
		t.Fatal("accepted missing product")
	}
	if err := repo.Create(&models.ShoppingCartItem{ShoppingCartID: cart.ID, ProductID: product.ID, Quantity: 1}); err != nil {
		t.Fatal(err)
	}
	if err := db.Preload("Items").First(&cart, cart.ID).Error; err != nil {
		t.Fatal(err)
	}
	if len(cart.Items) != 1 || cart.Items[0].Quantity != 3 || cart.Total != 15.75 {
		t.Fatalf("incorrect cart aggregation: %+v", cart)
	}
	if _, err := NewOrderRepository(db).FindOrCreate(profile.ID, cart.ID); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(&models.ShoppingCartItem{ShoppingCartID: cart.ID, ProductID: product.ID, Quantity: 1}); err == nil {
		t.Fatal("added item to finalized cart")
	}
	if err := repo.Delete(cart.ID, product.ID); err == nil {
		t.Fatal("removed item from finalized cart")
	}
}

func TestStockRejectsOverdrawAndInvalidMovement(t *testing.T) {
	db, profile, product := commerceDB(t)
	for _, movement := range []models.Stock{
		{ProfileID: profile.ID, ProductID: product.ID, Quantity: 11, Type: models.StockSaida},
		{ProfileID: profile.ID, ProductID: product.ID, Quantity: -1, Type: models.StockEntrada},
		{ProfileID: profile.ID, ProductID: product.ID, Quantity: 1, Type: "invalid"},
		{ProfileID: profile.ID, ProductID: product.ID + 1, Quantity: 1, Type: models.StockEntrada},
	} {
		if err := NewStockRepository(db).Create(&movement); err == nil {
			t.Fatalf("accepted invalid movement: %+v", movement)
		}
		assertStock(t, db, product.ID, 10)
	}
}

func TestProductUpdatePreservesStockAndValidatesPrice(t *testing.T) {
	db, _, product := commerceDB(t)
	repo := NewProductRepository(db)
	product.CurrentStock = 999
	product.Price = 6
	if err := repo.Update(&product); err != nil {
		t.Fatal(err)
	}
	assertStock(t, db, product.ID, 10)
	for _, price := range []float64{-1, 0, math.NaN(), math.Inf(1)} {
		product.Price = price
		if err := repo.Update(&product); err == nil {
			t.Fatalf("accepted invalid price %v", price)
		}
	}
}

func TestOrderUpdatePreservesPriceAndDeliverySnapshot(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 1)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	order.Total = 0.01
	order.ShoppingCart.Total = 0.01
	order.DeliveryPrice = 0
	order.Status = models.ProcessingStatus
	if err := repo.Update(&order); err == nil {
		t.Fatal("accepted a manipulated order total")
	}
	order, err = repo.FindById(order.ID)
	if err != nil {
		t.Fatal(err)
	}
	order.Status = models.ProcessingStatus
	if err := repo.Update(&order); err != nil {
		t.Fatal(err)
	}
	if order.Total != 7.75 || order.DeliveryPrice != 2.5 {
		t.Fatalf("order totals changed: %+v", order)
	}
	assertStock(t, db, product.ID, 9)
	order.Status = "invalid"
	if err := repo.Update(&order); err == nil {
		t.Fatal("accepted invalid order status")
	}
}

func TestPaymentCanChoosePickup(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 1)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	order.IsDelivery = false
	order.DeliveryPrice = 0
	order.Total = 5.25
	order.Status = models.ProcessingStatus
	order.PaymentMethod = models.CashPaymentMethod
	if err := repo.Update(&order); err != nil {
		t.Fatal(err)
	}
	if order.IsDelivery || order.DeliveryPrice != 0 || order.Total != 5.25 {
		t.Fatalf("pickup was not persisted: %+v", order)
	}
	order.Status = models.AwaitingPaymentStatus
	if err := repo.Update(&order); err == nil {
		t.Fatal("stale payment moved processing order backwards")
	}
}

func TestCancelAfterProductDeletionPreservesLedger(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 1)
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := NewProductRepository(db).Delete(&product); err != nil {
		t.Fatal(err)
	}
	if err := repo.Cancel(order.ID); err != nil {
		t.Fatal(err)
	}
	assertStock(t, db, product.ID, 10)
}

func TestConcurrentStockWithdrawalsCannotOversell(t *testing.T) {
	db, profile, product := commerceDB(t)
	results := make(chan error, 2)
	start := make(chan struct{})
	for i := 0; i < 2; i++ {
		go func() {
			<-start
			results <- NewStockRepository(db).Create(&models.Stock{ProfileID: profile.ID, ProductID: product.ID, Quantity: 7, Type: models.StockSaida})
		}()
	}
	close(start)
	var succeeded int
	for i := 0; i < 2; i++ {
		if <-results == nil {
			succeeded++
		}
	}
	if succeeded != 1 {
		t.Fatalf("successful withdrawals = %d, want 1", succeeded)
	}
	assertStock(t, db, product.ID, 3)
}

func TestDeletingCartItemRecalculatesTotal(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	if err := NewShoppingCartItemRepository(db).Delete(cart.ID, product.ID); err != nil {
		t.Fatal(err)
	}
	if err := db.First(&cart, cart.ID).Error; err != nil {
		t.Fatal(err)
	}
	if cart.Total != 0 {
		t.Fatalf("empty cart total = %v", cart.Total)
	}
}

func TestCartCreationRequiresExistingOwnerAndDoesNotInventIDs(t *testing.T) {
	db, _, _ := commerceDB(t)
	repo := NewShoppingCartRepository(db)
	if _, err := repo.FindOrCreateByUserId(999); err == nil {
		t.Fatal("created a cart without an owner")
	}
	if _, err := repo.FindOrCreateById(999); err == nil {
		t.Fatal("created a cart from an arbitrary ID")
	}
}

func TestCheckoutUsesCurrentPriceAndFreezesItForTheOrder(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	product.Price = 6
	if err := NewProductRepository(db).Update(&product); err != nil {
		t.Fatal(err)
	}
	repo := NewOrderRepository(db)
	order, err := repo.FindOrCreate(profile.ID, cart.ID)
	if err != nil {
		t.Fatal(err)
	}
	if order.Total != 14.50 || order.ShoppingCart.Items[0].ItemPrice != 6 {
		t.Fatalf("checkout did not refresh price: %+v", order)
	}
	product.Price = 7
	if err := NewProductRepository(db).Update(&product); err != nil {
		t.Fatal(err)
	}
	order, err = repo.FindById(order.ID)
	if err != nil || order.Total != 14.50 || order.ShoppingCart.Items[0].ItemPrice != 6 {
		t.Fatalf("historical price changed: %+v, %v", order, err)
	}
}

func TestLegacyOrderCancellationDoesNotInventAStockReservation(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	order := models.Order{ProfileID: profile.ID, ShoppingCartID: cart.ID, Status: models.ActiveStatus, PaymentMethod: models.CashPaymentMethod}
	if err := db.Create(&order).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&cart).UpdateColumn("order_id", order.ID).Error; err != nil {
		t.Fatal(err)
	}
	if err := NewOrderRepository(db).Cancel(order.ID); err != nil {
		t.Fatal(err)
	}
	assertStock(t, db, product.ID, 10)
}
