package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/utils"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

const (
	CreateProductPage = "products/create"
	ProductPage       = "products/product"
)

type ProductController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	RenderCreate(ctx *fiber.Ctx) error
	RenderProducts(ctx *fiber.Ctx) error
	RenderProduct(ctx *fiber.Ctx) error
	RenderDetails(ctx *fiber.Ctx) error
	RenderDelete(ctx *fiber.Ctx) error
	JSONProducts(ctx *fiber.Ctx) error
}

type productController struct {
	productService services.ProductService
}

func NewProductController(s services.ProductService) ProductController {
	return &productController{
		productService: s,
	}
}

func (c *productController) RenderCreate(ctx *fiber.Ctx) error {
	return views.Render(ctx, CreateProductPage, nil, views.BaseLayout)
}

func (c *productController) Create(ctx *fiber.Ctx) error {
	product := &models.Product{}
	if err := ctx.BodyParser(product); err != nil {
		return views.RenderError(ctx, CreateProductPage, product, "Dados do produto inválidos: "+err.Error(), views.BaseLayout)
	}
	if err := c.saveProductImage(ctx, product); err != nil {
		return views.RenderError(ctx, ProductPage, product, err.Error(), views.BaseLayout)
	}
	if err := c.productService.Create(product); err != nil {
		return views.RenderError(ctx, CreateProductPage, product, "Falha ao criar produto: "+err.Error(), views.BaseLayout)
	}
	return ctx.Redirect("/products")
}

func (c *productController) RenderDetails(ctx *fiber.Ctx) error {
	product, err := c.getProductByID(ctx)
	if err != nil {
		return ctx.Redirect("/store")
	}
	return views.Render(ctx, "products/details", product, views.StoreLayout)
}

func (c *productController) RenderProducts(ctx *fiber.Ctx) error {
	filter := c.getProductFilterFromQueryParams(ctx)
	products := c.productService.FindAll(filter)
	return views.Render(ctx, "products/products", fiber.Map{"Products": products, "Filter": filter}, views.BaseLayout)
}

func (c *productController) RenderProduct(ctx *fiber.Ctx) error {
	product, err := c.getProductByID(ctx)
	if err != nil {
		return ctx.Redirect("/products")
	}
	return views.Render(ctx, ProductPage, product, views.BaseLayout)
}

func (c *productController) RenderDelete(ctx *fiber.Ctx) error {
	product, err := c.getProductByID(ctx)
	if err != nil {
		return ctx.Redirect("/products")
	}
	return views.Render(ctx, "products/delete", product, views.BaseLayout)
}

func (c *productController) Update(ctx *fiber.Ctx) error {
	product, err := c.getProductByID(ctx)
	if err != nil {
		return ctx.Redirect("/products")
	}
	err = c.updateProductFromRequest(ctx, &product)
	if err != nil {
		return views.RenderError(ctx, ProductPage, &product, err.Error(), views.BaseLayout)
	}
	return ctx.Redirect("/products")
}

func (c *productController) Delete(ctx *fiber.Ctx) error {
	productID, err := c.getProductIdFromParam(ctx)
	if err != nil {
		return ctx.Redirect("/products")
	}
	if err := c.productService.Delete(productID); err != nil {
		return ctx.Redirect("/products")
	}
	return ctx.Redirect("/products")
}

func (c *productController) JSONProducts(ctx *fiber.Ctx) error {
	filter := c.getProductFilterFromQueryParams(ctx)
	products := c.productService.FindAll(filter)
	return ctx.JSON(map[string]interface{}{
		"Products": products,
		"Filter":   filter,
	})
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
		return views.RenderError(ctx, ProductPage, product, "Falha ao atualizar produto.", views.BaseLayout)
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

func (c *productController) getProductIdFromParam(ctx *fiber.Ctx) (uint, error) {
	return utils.StringToId(ctx.Params("id"))
}

func (c *productController) getProductFilterFromQueryParams(ctx *fiber.Ctx) *models.ProductFilter {
	return models.NewProductFilter(ctx.Query("q", ""), ctx.QueryInt("page"), ctx.QueryInt("limit"))
}

func (c *productController) getProductByID(ctx *fiber.Ctx) (models.Product, error) {
	productID, err := c.getProductIdFromParam(ctx)
	if err != nil {
		return models.Product{}, err
	}
	return c.productService.FindById(productID)
}
