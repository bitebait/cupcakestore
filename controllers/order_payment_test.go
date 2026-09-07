package controllers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func TestPaymentReopensLocalOrderWithoutFollowingHistoricalPixURLs(t *testing.T) {
	for _, tc := range []struct{ path, location string }{
		{"/invoice/123", "/orders/order/12"},
		{".attacker.example/invoice", "/orders/order/12"},
		{"//attacker.example/invoice", "/orders/order/12"},
		{"/\\attacker.example/invoice", "/orders/order/12"},
		{"https://attacker.example/invoice", "/orders/order/12"},
		{"", "/orders/order/12"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			controller := &orderController{}
			app := fiber.New()
			app.Get("/payment", func(ctx fiber.Ctx) error {
				return controller.processPaymentGet(ctx, &models.Order{Model: gorm.Model{ID: 12}, Status: models.AwaitingPaymentStatus, PaymentMethod: models.PixPaymentMethod, PixURL: tc.path})
			})
			res, err := app.Test(httptest.NewRequest("GET", "/payment", nil))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != fiber.StatusFound || res.Header.Get("Location") != tc.location {
				t.Fatalf("status=%d location=%q; want %q", res.StatusCode, res.Header.Get("Location"), tc.location)
			}
		})
	}
}

func TestPaymentFormRequiresExplicitFulfillmentAndCheckoutVersion(t *testing.T) {
	for _, body := range []string{
		"paymentMethod=Dinheiro", "paymentMethod=Dinheiro&isDelivery=unknown",
		"paymentMethod=Dinheiro&isDelivery=0", "paymentMethod=Dinheiro&isDelivery=0&storeVersion=invalid",
	} {
		app := fiber.New()
		app.Post("/payment", func(ctx fiber.Ctx) error {
			// No service: invalid input must be rejected before any payment work.
			if _, err := (&orderController{}).processPaymentPost(ctx, &models.Order{}); err == nil {
				t.Error("accepted incomplete checkout form")
			}
			return ctx.SendStatus(fiber.StatusBadRequest)
		})
		req := httptest.NewRequest("POST", "/payment", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != fiber.StatusBadRequest {
			t.Fatal(res.StatusCode)
		}
	}
}
