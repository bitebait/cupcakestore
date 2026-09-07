package controllers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type quantityCartService struct {
	services.ShoppingCartService
	profileID, productID uint
	quantity             int
}

func (s *quantityCartService) SetItemQuantity(profileID, productID uint, quantity int) error {
	s.profileID, s.productID, s.quantity = profileID, productID, quantity
	return nil
}

func TestCartQuantityFormUsesAuthenticatedProfileAndRouteProduct(t *testing.T) {
	service := &quantityCartService{}
	controller := NewShoppingCartController(service)
	app := fiber.New()
	app.Post("/cart/items/:id/quantity", func(ctx fiber.Ctx) error {
		ctx.Locals("Profile", &models.Profile{Model: gorm.Model{ID: 7}})
		return controller.SetItemQuantity(ctx)
	})
	req := httptest.NewRequest("POST", "/cart/items/2/quantity", strings.NewReader("quantity=3&profileID=99&productID=88&shoppingCartID=77&itemPrice=0.01"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != fiber.StatusSeeOther || res.Header.Get("Location") != "/cart" || service.profileID != 7 || service.productID != 2 || service.quantity != 3 {
		t.Fatalf("form changed identity: status=%d profile=%d product=%d quantity=%d", res.StatusCode, service.profileID, service.productID, service.quantity)
	}
}

func TestCartQuantityFormRejectsInvalidValuesBeforePersistence(t *testing.T) {
	for _, tc := range []struct{ id, quantity string }{
		{"0", "1"}, {"invalid", "1"}, {"2", ""}, {"2", "0"}, {"2", "-1"}, {"2", "1.5"}, {"2", "999999999999999999999999999999"},
	} {
		app := fiber.New()
		// A nil service makes accidental persistence fail immediately.
		app.Post("/cart/items/:id/quantity", NewShoppingCartController(nil).SetItemQuantity)
		req := httptest.NewRequest("POST", "/cart/items/"+tc.id+"/quantity", strings.NewReader("quantity="+tc.quantity))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != fiber.StatusSeeOther {
			t.Fatalf("unexpected response for %+v: %d", tc, res.StatusCode)
		}
	}
}
