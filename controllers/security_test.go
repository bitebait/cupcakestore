package controllers

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitebait/cupcakestore/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type requestValidator struct {
	validator *validator.Validate
}

func (v requestValidator) Validate(value any) error { return v.validator.Struct(value) }

func newBindingTestApp() *fiber.App {
	return fiber.New(fiber.Config{StructValidator: requestValidator{validator: validator.New()}})
}

func TestRegistrationIgnoresInternalAccountFields(t *testing.T) {
	for _, input := range []struct{ name, contentType, body string }{
		{"form", "application/x-www-form-urlencoded", "email=user%40example.com&password=password123&ID=88&isStaff=true&isActive=false"},
		{"json", "application/json", `{"email":"user@example.com","password":"password123","ID":88,"isStaff":true,"isActive":false}`},
	} {
		t.Run(input.name, func(t *testing.T) {
			app := newBindingTestApp()
			app.Post("/", func(c fiber.Ctx) error {
				user, err := parseUserFromContext(c)
				if err != nil {
					return err
				}
				if user.ID != 0 || user.IsStaff || !user.IsActive || user.Password != "password123" || user.Email != "user@example.com" {
					t.Errorf("registration accepted internal fields: ID=%d staff=%v active=%v", user.ID, user.IsStaff, user.IsActive)
				}
				return c.SendStatus(200)
			})
			req := httptest.NewRequest("POST", "/", strings.NewReader(input.body))
			req.Header.Set("Content-Type", input.contentType)
			response, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != 200 {
				t.Fatal(response.StatusCode)
			}
		})
	}
}

func TestProductFormPreservesIdentityStockAndImage(t *testing.T) {
	app := newBindingTestApp()
	product := models.Product{Model: gorm.Model{ID: 3}, CurrentStock: 7, Image: "/images/original.jpg"}
	app.Post("/", func(c fiber.Ctx) error { return readProductForm(c, &product) })
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

func TestMultipartProductFormPreservesServerFields(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for name, value := range map[string]string{"name": "Cupcake", "price": "12.50", "ID": "99", "currentStock": "500", "image": "https://example.com/tracker"} {
		if err := writer.WriteField(name, value); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	product := models.Product{Model: gorm.Model{ID: 3}, CurrentStock: 7, Image: "/images/original.jpg"}
	app := newBindingTestApp()
	app.Post("/", func(c fiber.Ctx) error { return readProductForm(c, &product) })
	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 || product.ID != 3 || product.CurrentStock != 7 || product.Image != "/images/original.jpg" || product.Price != 12.5 {
		t.Fatalf("invalid multipart binding: status=%d product=%+v", response.StatusCode, product)
	}
}

func TestBindingValidatesRegistrationAndStockBeforePersistence(t *testing.T) {
	for _, input := range []struct{ name, path, body string }{
		{"invalid email", "/register", "email=invalid&password=password123"},
		{"short password", "/register", "email=user%40example.com&password=short"},
		{"negative quantity", "/stock", "productID=1&quantity=-1&type=entrada"},
		{"missing product", "/stock", "quantity=1&type=entrada"},
		{"invalid stock type", "/stock", "productID=1&quantity=1&type=unknown"},
	} {
		t.Run(input.name, func(t *testing.T) {
			app := newBindingTestApp()
			app.Post("/register", func(c fiber.Ctx) error {
				if _, err := parseUserFromContext(c); err != nil {
					return fiber.ErrBadRequest
				}
				return c.SendStatus(fiber.StatusOK)
			})
			// A nil service makes accidental persistence immediately fail the test.
			app.Post("/stock", NewStockController(nil).Create)
			req := httptest.NewRequest("POST", input.path, strings.NewReader(input.body))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			response, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if response.StatusCode != fiber.StatusBadRequest {
				t.Fatalf("status=%d, expected bad request", response.StatusCode)
			}
		})
	}
}

type storeSettingsStub struct {
	config models.StoreConfig
}

func (s *storeSettingsStub) GetStoreConfig() (models.StoreConfig, error) { return s.config, nil }
func (s *storeSettingsStub) Update(config *models.StoreConfig) error {
	s.config = *config
	return nil
}

func TestSettingsBindingPreservesPartialFieldsAndIdentity(t *testing.T) {
	settings := &storeSettingsStub{config: models.StoreConfig{
		Model: gorm.Model{ID: 7}, DeliveryPrice: 5, DeliveryIsActive: true,
		PaymentCashIsActive: true, PhysicalStoreAddress: "Rua original", PixKey: "original@example.com",
	}}
	app := newBindingTestApp()
	app.Post("/", NewStoreConfigController(settings).Update)
	req := httptest.NewRequest("POST", "/", strings.NewReader("deliveryPrice=0&deliveryIsActive=0&ID=999"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != fiber.StatusFound || response.Header.Get("Location") != "/dashboard" {
		t.Fatalf("redirect changed: %d %q", response.StatusCode, response.Header.Get("Location"))
	}
	got := settings.config
	if got.ID != 7 || got.DeliveryPrice != 0 || got.DeliveryIsActive || !got.PaymentCashIsActive || got.PhysicalStoreAddress != "Rua original" || got.PixKey != "original@example.com" {
		t.Fatalf("partial settings update changed unrelated fields: %+v", got)
	}
}
