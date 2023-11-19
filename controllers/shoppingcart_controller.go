package controllers

import (
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type ShoppingCartController interface {
	RenderShoppingCart(ctx *fiber.Ctx) error
}

type shoppingCartController struct {
}

func NewShoppingCartController() ShoppingCartController {
	return &shoppingCartController{}
}

func (c *shoppingCartController) RenderShoppingCart(ctx *fiber.Ctx) error {

	return views.Render(ctx, "shoppingcart/shoppingcart", nil, "", storeLayout)
}
