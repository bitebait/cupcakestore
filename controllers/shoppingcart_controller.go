package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type ShoppingCartController interface {
	RenderShoppingCart(ctx *fiber.Ctx) error
}

type shoppingCartController struct {
	shoppingCartService services.ShoppingCartService
}

func NewShoppingCartController(shoppingCartService services.ShoppingCartService) ShoppingCartController {
	return &shoppingCartController{
		shoppingCartService: shoppingCartService,
	}
}

func (c *shoppingCartController) RenderShoppingCart(ctx *fiber.Ctx) error {
	userID := ctx.Locals("profile").(*models.Profile).ID
	cart, err := c.shoppingCartService.FindByUserId(userID)
	if err != nil {
		return ctx.Redirect("/")
	}
	return views.Render(ctx, "shoppingcart/shoppingcart", cart, "", storeLayout)
}
