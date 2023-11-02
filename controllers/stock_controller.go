package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type StockController interface {
	RenderCreate(ctx *fiber.Ctx) error
	HandlerCreate(ctx *fiber.Ctx) error
}

type stockController struct {
	stockService services.StockService
}

func NewStockController(s services.StockService) StockController {
	return &stockController{
		stockService: s,
	}
}

func (c *stockController) RenderCreate(ctx *fiber.Ctx) error {
	return views.Render(ctx, "stock/create", nil, "", baseLayout)
}

func (c *stockController) HandlerCreate(ctx *fiber.Ctx) error {
	return ctx.Redirect("/stock")
}
