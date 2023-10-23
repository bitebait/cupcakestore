package controllers

import (
	"fmt"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type ProductController interface {
	RenderCreate(ctx *fiber.Ctx) error
	HandlerCreate(ctx *fiber.Ctx) error
	RenderProducts(ctx *fiber.Ctx) error
}

type productController struct {
	productService services.ProductService
}

func NewProductController(p services.ProductService) ProductController {
	return &productController{
		productService: p,
	}
}

func (c *productController) RenderCreate(ctx *fiber.Ctx) error {
	return views.Render(ctx, "products/create", nil, "", baseLayout)
}

func (c *productController) HandlerCreate(ctx *fiber.Ctx) error {
	product := &models.Product{}
	if err := ctx.BodyParser(product); err != nil {
		errorMessage := "Dados do produto inválidos: " + err.Error()
		return views.Render(ctx, "products/create", nil, errorMessage, baseLayout)
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return err
	}

	if err := ctx.SaveFile(file, fmt.Sprintf("./web/images/%s", file.Filename)); err != nil {
		return err
	}
	product.Image = fmt.Sprintf("/images/%s", file.Filename)

	if err := c.productService.Create(product); err != nil {
		errorMessage := "Falha ao criar produto: " + err.Error()
		return views.Render(ctx, "products/create", nil, errorMessage, baseLayout)
	}

	return ctx.Redirect("/products")
}

func (c *productController) RenderProducts(ctx *fiber.Ctx) error {
	query := ctx.Query("q", "")

	pagination := models.NewPagination(ctx.QueryInt("page"), ctx.QueryInt("limit"))
	products := c.productService.FindAll(pagination, query)
	data := fiber.Map{
		"Products":   products,
		"Pagination": pagination,
	}

	return views.Render(ctx, "products/products", data, "", baseLayout)
}
