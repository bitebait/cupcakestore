package controllers

import (
	"log"
	"strconv"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type OrderController interface {
	RenderOrders(ctx *fiber.Ctx) error
	Checkout(ctx *fiber.Ctx) error
	Payment(ctx *fiber.Ctx) error
}

type orderController struct {
	orderService       services.OrderService
	storeConfigService services.StoreConfigService
}

func NewOrderController(orderService services.OrderService, storeConfigService services.StoreConfigService) OrderController {
	return &orderController{
		orderService:       orderService,
		storeConfigService: storeConfigService,
	}
}

func (c *orderController) Checkout(ctx *fiber.Ctx) error {
	profileID := c.getUserID(ctx)
	cartID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return c.renderError(ctx, "Erro ao processar o ID do carrinho.")
	}

	order, err := c.orderService.FindOrCreate(profileID, uint(cartID))
	if err != nil {
		return c.renderError(ctx, "Erro ao obter o carrinho de compras.")
	}

	if order.ShoppingCart.Total <= 0 || !order.IsActiveOrAwaitingPayment() {
		return ctx.Redirect("/cart")
	}

	storeConfig, err := c.storeConfigService.GetStoreConfig()
	if err != nil {
		return c.renderError(ctx, "Erro interno no servidor.")
	}

	data := fiber.Map{
		"Orders":      order,
		"StoreConfig": storeConfig,
	}
	return views.Render(ctx, "orders/checkout", data, "", storeLayout)
}

func (c *orderController) Payment(ctx *fiber.Ctx) error {
	cartID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return c.renderErrorMessage(ctx, err, "processar o checkout do carrinho")
	}

	order, err := c.orderService.FindByCartId(uint(cartID))
	if err != nil {
		return c.renderErrorMessage(ctx, err, "processar o checkout do carrinho")
	}

	if order.ShoppingCart.Total <= 0 {
		return ctx.Redirect("/cart")
	}

	switch ctx.Method() {
	case fiber.MethodPost:
		if !order.IsActiveOrAwaitingPayment() {
			return ctx.Redirect("/cart")
		}

		if err := ctx.BodyParser(order); err != nil {
			log.Println(err)
			return c.renderErrorMessage(ctx, err, "processar os dados de pagamento")
		}

		if err := c.orderService.Update(order); err != nil {
			return c.renderErrorMessage(ctx, err, "atualizar o carrinho para pagamento")
		}

		if err := c.orderService.Payment(order); err != nil {
			return c.renderErrorMessage(ctx, err, "realizar o pagamento do carrinho")
		}

		if order.PaymentMethod == models.PixPaymentMethod {
			return ctx.Redirect("https://pix.ae" + order.PixURL)
		}
		return ctx.Redirect("/orders")
	case fiber.MethodGet:
		if order.Status == models.AwaitingPaymentStatus && order.PaymentMethod == models.PixPaymentMethod {
			return ctx.Redirect("https://pix.ae" + order.PixURL)
		}
	default:
		return ctx.Redirect("/orders")
	}

	return ctx.Redirect("/orders")
}

func (c *orderController) RenderOrders(ctx *fiber.Ctx) error {
	profileID := c.getUserID(ctx)
	page := ctx.QueryInt("page")
	limit := ctx.QueryInt("limit")
	filter := models.NewOrderFilter(profileID, page, limit)
	orders := c.orderService.FindAll(filter)

	return views.Render(ctx, "orders/orders", fiber.Map{"Orders": orders, "Filter": filter}, "", storeLayout)
}

func (c *orderController) getUserID(ctx *fiber.Ctx) uint {
	return ctx.Locals("profile").(*models.Profile).ID
}

func (c *orderController) renderErrorMessage(ctx *fiber.Ctx, err error, action string) error {
	errorMessage := "Houve um erro ao " + action + ": " + err.Error()
	return c.renderError(ctx, errorMessage)
}

func (c *orderController) renderError(ctx *fiber.Ctx, errorMessage string) error {
	return views.Render(ctx, "orders/order", nil, errorMessage, storeLayout)
}
