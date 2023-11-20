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
	RemoveFromCart(ctx *fiber.Ctx) error
	Checkout(ctx *fiber.Ctx) error
	Payment(ctx *fiber.Ctx) error
}

type shoppingCartController struct {
	shoppingCartService services.ShoppingCartService
	storeConfigService  services.StoreConfigService
}

func NewShoppingCartController(shoppingCartService services.ShoppingCartService, storeConfigService services.StoreConfigService) ShoppingCartController {
	return &shoppingCartController{
		shoppingCartService: shoppingCartService,
		storeConfigService:  storeConfigService,
	}
}

func (c *shoppingCartController) RenderShoppingCart(ctx *fiber.Ctx) error {
	userID := c.getUserID(ctx)
	cart, err := c.shoppingCartService.FindByUserId(userID)
	if err != nil {
		return c.renderError(ctx, "Erro ao obter o carrinho de compras.")
	}
	return views.Render(ctx, "shoppingcart/shoppingcart", cart, "", storeLayout)
}

func (c *shoppingCartController) AddShoppingCartItem(ctx *fiber.Ctx) error {
	productID, err := utils.StringToId(ctx.FormValue("id"))
	if err != nil {
		return c.renderErrorMessage(ctx, err, "adicionar o item ao carrinho de compras")
	}
	quantity, err := strconv.Atoi(ctx.FormValue("quantity"))
	if err != nil {
		return c.renderErrorMessage(ctx, err, "adicionar o item ao carrinho de compras")
	}

	userID := c.getUserID(ctx)

	if err = c.shoppingCartService.AddItemToCart(userID, productID, quantity); err != nil {
		return c.renderErrorMessage(ctx, err, "adicionar o item ao carrinho de compras")
	}

	return ctx.Redirect("/cart")
}

func (c *shoppingCartController) RemoveFromCart(ctx *fiber.Ctx) error {
	productID, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return c.renderErrorMessage(ctx, err, "remover o item do carrinho de compras")
	}

	userID := c.getUserID(ctx)

	if err = c.shoppingCartService.RemoveFromCart(userID, productID); err != nil {
		return c.renderErrorMessage(ctx, err, "remover o item do carrinho de compras")
	}

	return ctx.Redirect("/cart")
}

func (c *shoppingCartController) renderErrorMessage(ctx *fiber.Ctx, err error, action string) error {
	errorMessage := "Houve um erro ao " + action + ": " + err.Error()
	return c.renderError(ctx, errorMessage)
}

func (c *shoppingCartController) getUserID(ctx *fiber.Ctx) uint {
	return ctx.Locals("profile").(*models.Profile).ID
}

func (c *shoppingCartController) renderError(ctx *fiber.Ctx, errorMessage string) error {
	return views.Render(ctx, "shoppingcart/shoppingcart", nil, errorMessage, storeLayout)
}

func (c *shoppingCartController) Checkout(ctx *fiber.Ctx) error {
	cartID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return c.renderError(ctx, "Erro ao processar o ID do carrinho.")
	}

	cart, err := c.shoppingCartService.FindById(uint(cartID))
	if err != nil {
		return c.renderError(ctx, "Erro ao obter o carrinho de compras.")
	}

	if cart.Total <= 0 {
		return ctx.Redirect("/cart")
	}

	if !(cart.Status == models.ActiveStatus || cart.Status == models.AwaitingPaymentStatus) {
		return ctx.Redirect("/cart")
	}

	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		return c.renderError(ctx, "Erro interno no servidor.")
	}

	data := fiber.Map{
		"ShoppingCart": cart,
		"StoreConfig":  storeConfig,
	}
	return views.Render(ctx, "shoppingcart/checkout", data, "", storeLayout)
}

func (c *shoppingCartController) Payment(ctx *fiber.Ctx) error {
	cartID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return c.renderErrorMessage(ctx, err, "processar o checkout do carrinho")
	}
	cart, err := c.shoppingCartService.FindById(uint(cartID))
	if err != nil {
		return c.renderErrorMessage(ctx, err, "processar o checkout do carrinho")
	}

	if cart.Total <= 0 {
		return ctx.Redirect("/cart")
	}

	if ctx.Method() == fiber.MethodPost && cart.PaymentMethod == models.PixPaymentMethod && (cart.Status == models.ActiveStatus || cart.Status == models.AwaitingPaymentStatus) {
		if err := ctx.BodyParser(&cart); err != nil {
			return c.renderErrorMessage(ctx, err, "processar os dados de pagamento")
		}

		if err := c.shoppingCartService.Update(&cart); err != nil {
			return c.renderErrorMessage(ctx, err, "atualizar o carrinho para pagamento")
		}

		if err := c.shoppingCartService.Payment(&cart); err != nil {
			return c.renderErrorMessage(ctx, err, "realizar o pagamento do carrinho")
		}

		return ctx.Redirect("https://pix.ae" + cart.PixURL)

	} else if ctx.Method() == fiber.MethodGet && cart.PaymentMethod == models.PixPaymentMethod && cart.Status == models.AwaitingPaymentStatus {
		return ctx.Redirect("https://pix.ae" + cart.PixURL)
	} else {
		return ctx.Redirect("/")
	}
}
