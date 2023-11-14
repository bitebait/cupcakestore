package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/utils"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type ProductController interface {
	RenderCreate(ctx *fiber.Ctx) error
	HandlerCreate(ctx *fiber.Ctx) error
	RenderProducts(ctx *fiber.Ctx) error
	RenderProduct(ctx *fiber.Ctx) error
	HandlerUpdate(ctx *fiber.Ctx) error
	RenderDelete(ctx *fiber.Ctx) error
	HandlerDelete(ctx *fiber.Ctx) error
	JSONProducts(ctx *fiber.Ctx) error
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

	if err := c.saveProductImage(ctx, product); err != nil {
		return views.Render(ctx, "products/product", product, err.Error(), baseLayout)
	}

	if err := c.productService.Create(product); err != nil {
		errorMessage := "Falha ao criar produto: " + err.Error()
		return views.Render(ctx, "products/create", nil, errorMessage, baseLayout)
	}

	return ctx.Redirect("/products")
}

func (c *productController) RenderProducts(ctx *fiber.Ctx) error {
	filter := &models.ProductFilter{
		Product: &models.Product{
			Name: ctx.Query("q", ""),
		},
		Pagination: models.NewPagination(ctx.QueryInt("page"), ctx.QueryInt("limit")),
	}

	products := c.productService.FindAll(filter)

	data := map[string]interface{}{
		"Products": products,
		"Filter":   filter,
	}

	return views.Render(ctx, "products/products", data, "", baseLayout)
}

func (c *productController) RenderProduct(ctx *fiber.Ctx) error {
	productID, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return ctx.Redirect("/products")
	}

	product, err := c.productService.FindById(productID)
	if err != nil {
		return ctx.Redirect("/products")
	}

	return views.Render(ctx, "products/product", product, "", baseLayout)
}

func (c *productController) HandlerUpdate(ctx *fiber.Ctx) error {
	id, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return ctx.Redirect("/products")
	}

	product, err := c.productService.FindById(id)
	if err != nil {
		return ctx.Redirect("/products")
	}

	err = c.updateProductFromRequest(ctx, &product)
	if err != nil {
		return views.Render(ctx, "products/product", product, err.Error(), baseLayout)
	}

	return ctx.Redirect("/products")
}

func (c *productController) updateProductFromRequest(ctx *fiber.Ctx, product *models.Product) error {
	oldImage := product.Image

	if err := ctx.BodyParser(product); err != nil {
		return err
	}

	product.IsActive = ctx.FormValue("isActive") == "on"
	if err := c.saveProductImage(ctx, product); err != nil {
		product.Image = oldImage
	}

	if err := c.productService.Update(product); err != nil {
		return views.Render(ctx, "products/product", product, "Falha ao atualizar produto.", baseLayout)
	}
	return nil
}

func (c *productController) saveProductImage(ctx *fiber.Ctx, product *models.Product) error {
	imageFile, err := ctx.FormFile("image")
	if err != nil {
		return err
	}

	img := &models.ProductImage{}
	if err := img.Save(imageFile); err != nil {
		return err
	}

	product.Image = img.Path
	return nil
}

func (c *productController) RenderDelete(ctx *fiber.Ctx) error {
	id, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return ctx.Redirect("/products")
	}

	product, err := c.productService.FindById(id)
	if err != nil {
		return ctx.Redirect("/products")
	}

	return views.Render(ctx, "products/delete", product, "", baseLayout)
}

func (c *productController) HandlerDelete(ctx *fiber.Ctx) error {
	id, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return ctx.Redirect("/products")
	}

	err = c.productService.Delete(id)
	if err != nil {
		return ctx.Redirect("/products")
	}

	return ctx.Redirect("/products")
}

func (c *productController) JSONProducts(ctx *fiber.Ctx) error {
	filter := &models.ProductFilter{
		Product: &models.Product{
			Name: ctx.Query("q", ""),
		},
		Pagination: models.NewPagination(ctx.QueryInt("page"), ctx.QueryInt("limit")),
	}

	products := c.productService.FindAll(filter)

	data := map[string]interface{}{
		"Products": products,
		"Filter":   filter,
	}

	return ctx.JSON(data)
}
