package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/utils"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ShoppingCartController interface {
	RenderShoppingCart(ctx *fiber.Ctx) error
	AddShoppingCartItem(ctx *fiber.Ctx) error
	RemoveFromCart(ctx *fiber.Ctx) error
	CountShoppingCart(ctx *fiber.Ctx) error
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
	userID := getUserID(ctx)
	cart, err := c.shoppingCartService.FindByUserId(userID)
	if err != nil {
		return renderErrorMessage(err, "obter o carrinho de compras.")
	}
	return views.Render(ctx, "shoppingcart/shoppingcart", cart, views.StoreLayout)
}

func (c *shoppingCartController) AddShoppingCartItem(ctx *fiber.Ctx) error {
	productID, err := utils.StringToId(ctx.FormValue("id"))
	if err != nil {
		return renderErrorMessage(err, "adicionar o item ao carrinho de compras")
	}
	quantity, err := strconv.Atoi(ctx.FormValue("quantity"))
	if err != nil {
		return renderErrorMessage(err, "adicionar o item ao carrinho de compras")
	}

	userID := getUserID(ctx)

	if err = c.shoppingCartService.AddItemToCart(userID, productID, quantity); err != nil {
		return renderErrorMessage(err, "adicionar o item ao carrinho de compras")
	}

	return ctx.Redirect("/cart")
}

func (c *shoppingCartController) RemoveFromCart(ctx *fiber.Ctx) error {
	productID, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return renderErrorMessage(err, "remover o item do carrinho de compras")
	}

	userID := getUserID(ctx)

	if err = c.shoppingCartService.RemoveFromCart(userID, productID); err != nil {
		return renderErrorMessage(err, "remover o item do carrinho de compras")
	}

	return ctx.Redirect("/cart")
}

func (c *shoppingCartController) CountShoppingCart(ctx *fiber.Ctx) error {
	userID := getUserID(ctx)
	cart, err := c.shoppingCartService.FindByUserId(userID)
	if err != nil {
		return err
	}

	itemCount := cart.CountItems()
	return ctx.JSON(fiber.Map{"itemCount": itemCount})
}
