package controllers

import (
	"fmt"
	"github.com/bitebait/cupcakestore/helpers"
	"github.com/bitebait/cupcakestore/messages"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v3"
)

const (
	StockCreateView = "stock/create"
	StocksView      = "stock/stocks"
	StockView       = "stock/stock"
)

type StockController interface {
	Create(ctx fiber.Ctx) error
	RenderCreate(ctx fiber.Ctx) error
	RenderStocks(ctx fiber.Ctx) error
	RenderStock(ctx fiber.Ctx) error
}

type stockController struct {
	stockService services.StockService
}

func NewStockController(s services.StockService) StockController {
	return &stockController{
		stockService: s,
	}
}

func (c *stockController) RenderCreate(ctx fiber.Ctx) error {
	return ctx.Render(StockCreateView, fiber.Map{}, views.BaseLayout)
}

func (c *stockController) Create(ctx fiber.Ctx) error {
	var input struct {
		ProductID uint   `form:"productID" validate:"required,gt=0"`
		Quantity  int    `form:"quantity" validate:"required,gt=0"`
		Type      string `form:"type" validate:"required,oneof=entrada saída"`
	}
	if err := ctx.Bind().Body(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "dados de estoque inválidos")
	}
	stock := models.Stock{ProductID: input.ProductID, Quantity: input.Quantity}
	switch input.Type {
	case string(models.StockEntrada):
		stock.Type = models.StockEntrada
	case string(models.StockSaida):
		stock.Type = models.StockSaida
	default:
		return fiber.NewError(fiber.StatusBadRequest, "tipo de movimentação inválido")
	}

	profile := fiber.Locals[*models.Profile](ctx, "Profile")
	if profile == nil || profile.ID == 0 {
		messages.SetErrorMessage(ctx, "por favor faça o login e tente novamente")
		return ctx.Render(StockCreateView, fiber.Map{}, views.BaseLayout)
	}

	stock.ProfileID = profile.ID
	if err := c.stockService.Create(&stock); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Render(StockCreateView, fiber.Map{}, views.BaseLayout)
	}

	messages.SetSuccessMessage(ctx, "estoque adicionado com sucesso")
	return ctx.Redirect().Status(fiber.StatusFound).To(fmt.Sprintf("/stock/%v", stock.ProductID))
}

func (c *stockController) RenderStocks(ctx fiber.Ctx) error {
	return ctx.Render(StocksView, fiber.Map{}, views.BaseLayout)
}

func (c *stockController) RenderStock(ctx fiber.Ctx) error {
	productID, err := helpers.ParseStringToID(ctx.Params("id"))

	if err != nil {
		messages.SetErrorMessage(ctx, "falha ao identificar o produto: "+err.Error())
		return ctx.Redirect().Status(fiber.StatusFound).To("/stocks")
	}

	page := fiber.Query[int](ctx, "page")
	limit := fiber.Query[int](ctx, "limit")
	filter := models.NewStockFilter(productID, page, limit)
	stocks := c.stockService.FindByProductId(filter)
	data := fiber.Map{
		"Stocks": stocks,
		"Filter": filter,
	}

	return ctx.Render(StockView, fiber.Map{"Object": data}, views.BaseLayout)
}
