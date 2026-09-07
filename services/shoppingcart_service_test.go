package services

import (
	"testing"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"gorm.io/gorm"
)

type quantityCartRepository struct {
	repositories.ShoppingCartRepository
	requestedProfileID uint
}

func (r *quantityCartRepository) FindOrCreateByUserId(profileID uint) (models.ShoppingCart, error) {
	r.requestedProfileID = profileID
	return models.ShoppingCart{Model: gorm.Model{ID: 7}, ProfileID: profileID, Items: []models.ShoppingCartItem{{Model: gorm.Model{ID: 3}, ShoppingCartID: 7, ProductID: 2, Quantity: 5, ItemPrice: 10}}}, nil
}

type quantityItemService struct {
	ShoppingCartItemService
	updated *models.ShoppingCartItem
}

func (s *quantityItemService) Update(item *models.ShoppingCartItem) error {
	s.updated = item
	return nil
}

func TestSetCartQuantityUsesExistingItemFromCurrentProfile(t *testing.T) {
	repo, items := &quantityCartRepository{}, &quantityItemService{}
	service := NewShoppingCartService(repo, items)
	if err := service.SetItemQuantity(11, 2, 3); err != nil {
		t.Fatal(err)
	}
	if repo.requestedProfileID != 11 || items.updated == nil || items.updated.ID != 3 || items.updated.ShoppingCartID != 7 || items.updated.ProductID != 2 || items.updated.Quantity != 3 {
		t.Fatalf("quantity update changed identity or added quantity: profile=%d item=%+v", repo.requestedProfileID, items.updated)
	}
}

func TestSetCartQuantityRejectsInvalidAndMissingItems(t *testing.T) {
	for _, tc := range []struct {
		profileID, productID uint
		quantity             int
	}{
		{0, 2, 3}, {11, 0, 3}, {11, 2, 0}, {11, 2, -1}, {11, 88, 1},
	} {
		repo, items := &quantityCartRepository{}, &quantityItemService{}
		if err := NewShoppingCartService(repo, items).SetItemQuantity(tc.profileID, tc.productID, tc.quantity); err == nil || items.updated != nil {
			t.Fatalf("invalid update accepted: %+v err=%v", tc, err)
		}
	}
}
