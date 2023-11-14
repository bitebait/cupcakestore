package controllers

import (
	"fmt"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/utils"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type StockController interface {
	RenderCreate(ctx *fiber.Ctx) error
	HandlerCreate(ctx *fiber.Ctx) error
	RenderStocks(ctx *fiber.Ctx) error
	RenderStock(ctx *fiber.Ctx) error
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
	stock := &models.Stock{}
	if err := ctx.BodyParser(stock); err != nil {
		return views.Render(ctx, "stock/create", nil,
			"Dados inválidos: "+err.Error(), baseLayout)
	}

	if err := c.stockService.Create(stock); err != nil {
		return views.Render(ctx, "stock/create", nil,
			"Falha ao adicionar ao estoque: "+err.Error(), baseLayout)
	}

	return ctx.Redirect(fmt.Sprintf("/stock/%v", stock.ProductID))
}

func (c *stockController) RenderStocks(ctx *fiber.Ctx) error {
	return views.Render(ctx, "stock/stocks", nil, "", baseLayout)
}

func (c *stockController) RenderStock(ctx *fiber.Ctx) error {
	productID, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return ctx.Redirect("/stocks")
	}
	page := ctx.QueryInt("page")
	limit := ctx.QueryInt("limit")
	filter := models.NewStockFilter(productID, page, limit)
	stocks := c.stockService.FindByProductId(filter)
	return views.Render(ctx, "stock/stock", fiber.Map{"Stocks": stocks, "Filter": filter}, "", baseLayout) // Updated to use "Stocks" instead of "Stock"
}
