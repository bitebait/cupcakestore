package controllers

import (
	"github.com/bitebait/cupcakestore/helpers"
	"github.com/bitebait/cupcakestore/messages"
	"github.com/bitebait/cupcakestore/services"
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
	cart, err := c.shoppingCartService.FindOrCreateByUserId(userID)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/")
	}
	return ctx.Render("shoppingcart/shoppingcart", fiber.Map{"Object": cart}, views.StoreLayout)
}

func (c *shoppingCartController) AddShoppingCartItem(ctx *fiber.Ctx) error {
	productIDStr := ctx.FormValue("id")
	quantityStr := ctx.FormValue("quantity", "1")

	productID, err := helpers.ParseStringToID(productIDStr)
	if err != nil {
		messages.SetErrorMessage(ctx, "falha ao identificar o produto: "+err.Error())
		return ctx.Redirect("/cart")
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		messages.SetErrorMessage(ctx, "falha ao identificar a quantidade: "+err.Error())
		return ctx.Redirect("/cart")
	}

	userID := getUserID(ctx)
	if err = c.shoppingCartService.AddItemToCart(userID, productID, quantity); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/cart")
	}

	messages.SetSuccessMessage(ctx, "produto adicionado ao carrinho")
	return ctx.Redirect("/cart")
}

func (c *shoppingCartController) RemoveFromCart(ctx *fiber.Ctx) error {
	productIDStr := ctx.Params("id")

	productID, err := helpers.ParseStringToID(productIDStr)
	if err != nil {
		messages.SetErrorMessage(ctx, "falha ao identificar o produto: "+err.Error())
		return ctx.Redirect("/cart")
	}

	userID := getUserID(ctx)
	if err = c.shoppingCartService.RemoveFromCart(userID, productID); err != nil {

		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/cart")
	}

	messages.SetSuccessMessage(ctx, "produto removido do carrinho")
	return ctx.Redirect("/cart")
}

func (c *shoppingCartController) CountShoppingCart(ctx *fiber.Ctx) error {
	userID := getUserID(ctx)
	cart, err := c.shoppingCartService.FindOrCreateByUserId(userID)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/cart")
	}
	itemCount := cart.CountItems()
	return ctx.JSON(fiber.Map{"itemCount": itemCount})
}
