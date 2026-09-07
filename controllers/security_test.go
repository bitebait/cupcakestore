package controllers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TestRegistrationIgnoresInternalAccountFields(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		user, err := parseUserFromContext(c)
		if err != nil {
			return err
		}
		if user.ID != 0 || user.IsStaff || !user.IsActive || user.Password != "password123" || user.Email != "user@example.com" {
			t.Errorf("registration accepted internal fields: ID=%d staff=%v active=%v", user.ID, user.IsStaff, user.IsActive)
		}
		return c.SendStatus(200)
	})
	req := httptest.NewRequest("POST", "/", strings.NewReader("email=user%40example.com&password=password123&ID=88&isStaff=true&isActive=false"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatal(response.StatusCode)
	}
}

func TestProductFormPreservesIdentityStockAndImage(t *testing.T) {
	app := fiber.New()
	product := models.Product{Model: gorm.Model{ID: 3}, CurrentStock: 7, Image: "/images/original.jpg"}
	app.Post("/", func(c *fiber.Ctx) error { return readProductForm(c, &product) })
	req := httptest.NewRequest("POST", "/", strings.NewReader("name=Cupcake&price=12.50&ID=99&currentStock=500&image=https://example.com/tracker"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatal(response.StatusCode)
	}
	if product.ID != 3 || product.CurrentStock != 7 || product.Image != "/images/original.jpg" || product.Price != 12.5 {
		t.Fatalf("unexpected product: %+v", product)
	}
}
