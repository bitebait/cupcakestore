package repositories

import (
	"testing"

	"github.com/bitebait/cupcakestore/models"
)

func TestCartQuantityReplacesQuantityAndPreservesTotalsOnFailure(t *testing.T) {
	db, profile, product := commerceDB(t)
	cart := cartWithItem(t, db, profile.ID, product.ID, 2)
	var item models.ShoppingCartItem
	if err := db.Where("shopping_cart_id = ?", cart.ID).First(&item).Error; err != nil {
		t.Fatal(err)
	}
	repo := NewShoppingCartItemRepository(db)
	item.Quantity = 3
	if err := repo.Update(&item); err != nil {
		t.Fatal(err)
	}
	if err := db.Preload("Items").First(&cart, cart.ID).Error; err != nil {
		t.Fatal(err)
	}
	if len(cart.Items) != 1 || cart.Items[0].Quantity != 3 || cart.Total != 15.75 {
		t.Fatalf("incorrect replacement: %+v", cart)
	}
	item.Quantity = 11
	if err := repo.Update(&item); err == nil {
		t.Fatal("allowed quantity over available stock")
	}
	if err := db.Preload("Items").First(&cart, cart.ID).Error; err != nil {
		t.Fatal(err)
	}
	if cart.Items[0].Quantity != 3 || cart.Total != 15.75 {
		t.Fatal("failed update changed cart")
	}
	if _, err := NewOrderRepository(db).FindOrCreate(profile.ID, cart.ID); err != nil {
		t.Fatal(err)
	}
	item.Quantity = 1
	if err := repo.Update(&item); err == nil {
		t.Fatal("changed quantity after checkout")
	}
	assertStock(t, db, product.ID, 7)
}
