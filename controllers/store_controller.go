package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type StoreController interface {
	RenderStore(ctx *fiber.Ctx) error
}

type storeController struct {
	productService services.ProductService
}

func NewStoreController(productService services.ProductService) StoreController {
	return &storeController{
		productService: productService,
	}
}

func (c *storeController) RenderStore(ctx *fiber.Ctx) error {
	query := ctx.Query("q", "")
	page := ctx.QueryInt("page")
	limit := ctx.QueryInt("limit")
	filter := models.NewProductFilter(query, page, limit)
	products := c.productService.FindActiveWithStock(filter)

	data := fiber.Map{
		"Products": products,
		"Filter":   filter,
	}

	return ctx.Render("store/store", fiber.Map{"Object": data}, views.StoreLayout)
}
