package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/utils"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ShoppingCartController interface {
	RenderShoppingCart(ctx *fiber.Ctx) error
	AddShoppingCartItem(ctx *fiber.Ctx) error
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

func (c *shoppingCartController) AddShoppingCartItem(ctx *fiber.Ctx) error {
	productID, err := utils.StringToId(ctx.FormValue("id"))
	if err != nil {
		errorMessage := "Houve um erro ao adicionar o item ao carrinho de compras: " + err.Error()
		return c.renderError(ctx, errorMessage)
	}
	quantityStr := ctx.FormValue("quantity")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		errorMessage := "Houve um erro ao adicionar o item ao carrinho de compras: " + err.Error()
		return c.renderError(ctx, errorMessage)
	}

	userID := c.getUserID(ctx)

	err = c.shoppingCartService.AddItemToCart(userID, productID, quantity)
	if err != nil {
		errorMessage := "Houve um erro ao adicionar o item ao carrinho de compras: " + err.Error()
		return c.renderError(ctx, errorMessage)
	}

	return ctx.Redirect("/cart")
}

func (c *shoppingCartController) getUserID(ctx *fiber.Ctx) uint {
	return ctx.Locals("profile").(*models.Profile).ID
}

func (c *shoppingCartController) renderError(ctx *fiber.Ctx, errorMessage string) error {
	return views.Render(ctx, "shoppingcart/shoppingcart", nil, errorMessage, storeLayout)
}
