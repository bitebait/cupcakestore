package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v3"
)

func selectLayout(isStaff, isUserProfile bool) string {
	if isStaff {
		return views.BaseLayout
	}
	if isUserProfile {
		return views.StoreLayout
	}
	return views.BaseLayout
}

// getProfileID returns the profile key used by carts, orders and stock.
func getProfileID(ctx fiber.Ctx) uint {
	profile := fiber.Locals[*models.Profile](ctx, "Profile")
	if profile == nil {
		return 0
	}
	return profile.ID
}
