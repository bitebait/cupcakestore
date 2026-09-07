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

type confirmationOrderService struct {
	services.OrderService
	orderID, staffID uint
}

func (s *confirmationOrderService) ConfirmPayment(orderID, staffID uint) error {
	s.orderID, s.staffID = orderID, staffID
	return nil
}

func TestBankConfirmationRequiresExplicitCheckAndAuthenticatedStaff(t *testing.T) {
	for _, tc := range []struct {
		name   string
		staff  bool
		body   string
		want   int
		called bool
	}{
		{"customer", false, "payment_received=on", fiber.StatusForbidden, false},
		{"unchecked", true, "", fiber.StatusSeeOther, false},
		{"confirmed", true, "payment_received=on&staffProfileID=999&orderID=777", fiber.StatusSeeOther, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := &confirmationOrderService{}
			app := fiber.New()
			app.Post("/orders/order/:id/confirm-payment", func(ctx fiber.Ctx) error {
				ctx.Locals("Profile", &models.Profile{Model: gorm.Model{ID: 8}, User: models.User{IsStaff: tc.staff}})
				return NewOrderController(service, nil).ConfirmPayment(ctx)
			})
			req := httptest.NewRequest("POST", "/orders/order/12/confirm-payment", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			res, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.want || (service.orderID != 0) != tc.called {
				t.Fatalf("status=%d recorded=%d", res.StatusCode, service.orderID)
			}
			if tc.called && (service.orderID != 12 || service.staffID != 8) {
				t.Fatal("form forged confirmation identity")
			}
		})
	}
}
