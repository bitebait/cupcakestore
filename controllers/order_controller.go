package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type OrderController interface {
	RenderOrders(ctx *fiber.Ctx) error
}

type orderController struct {
	orderService services.OrderService
}

func NewOrderController(orderService services.OrderService) OrderController {
	return &orderController{
		orderService: orderService,
	}
}

func (c *orderController) RenderOrders(ctx *fiber.Ctx) error {
	profileID := c.getUserID(ctx)
	page := ctx.QueryInt("page")
	limit := ctx.QueryInt("limit")
	filter := models.NewShoppingCartFilter(profileID, page, limit)
	orders := c.orderService.FindAll(filter)

	return views.Render(ctx, "shoppingcart/orders", fiber.Map{"Orders": orders, "Filter": filter}, "", storeLayout)
}

func (c *orderController) getUserID(ctx *fiber.Ctx) uint {
	return ctx.Locals("profile").(*models.Profile).ID
}
