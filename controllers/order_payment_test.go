package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v3"
)

func TestPaymentRedirectRejectsUnsafeHistoricalPixPaths(t *testing.T) {
	for _, tc := range []struct{ path, location string }{
		{"/invoice/123", "https://pix.ae/invoice/123"},
		{".attacker.example/invoice", "/orders"},
		{"//attacker.example/invoice", "/orders"},
		{"/\\attacker.example/invoice", "/orders"},
		{"https://attacker.example/invoice", "/orders"},
		{"", "/orders"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			controller := &orderController{}
			app := fiber.New()
			app.Get("/payment", func(ctx fiber.Ctx) error {
				return controller.processPaymentGet(ctx, &models.Order{Status: models.AwaitingPaymentStatus, PaymentMethod: models.PixPaymentMethod, PixURL: tc.path})
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
