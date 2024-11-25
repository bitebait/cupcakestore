package controllers

import (
	"github.com/bitebait/cupcakestore/helpers"
	"github.com/bitebait/cupcakestore/messages"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
	return ctx.Render("products/create", fiber.Map{}, views.BaseLayout)
}

func (c *productController) Create(ctx *fiber.Ctx) error {
	var product models.Product

	if err := ctx.BodyParser(&product); err != nil {
		messages.SetErrorMessage(ctx, "erro ao processar os dados do produto")
		return ctx.Redirect("/products/create")
	}

	if err := c.saveProductImage(ctx, &product); err != nil {
		messages.SetErrorMessage(ctx, "erro ao processar a imagem do produto")
		return ctx.Redirect("/products/create")
	}

	if err := c.productService.Create(&product); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/products/create")
	}

	messages.SetSuccessMessage(ctx, "produto criado com sucesso")
	return ctx.Redirect("/products")
}

func (c *productController) RenderDetails(ctx *fiber.Ctx) error {
	product, err := c.getProductByID(ctx)
	if err != nil {
		messages.SetErrorMessage(ctx, "falha ao identificar o produto: "+err.Error())
		return ctx.Redirect("/store")
	}

	return ctx.Render("products/details", fiber.Map{"Object": product}, views.StoreLayout)
}

func (c *productController) RenderProducts(ctx *fiber.Ctx) error {
	filter := c.getProductFilterFromQueryParams(ctx)
	products := c.productService.FindAll(filter)
	data := fiber.Map{"Products": products, "Filter": filter}

	return ctx.Render("products/products", fiber.Map{"Object": data}, views.BaseLayout)
}

func (c *productController) RenderProduct(ctx *fiber.Ctx) error {
	product, err := c.getProductByID(ctx)
	if err != nil {
		messages.SetErrorMessage(ctx, "falha ao identificar o produto: "+err.Error())
		return ctx.Redirect("/products")
	}
	return ctx.Render("products/product", fiber.Map{"Object": product}, views.BaseLayout)
}

func (c *productController) RenderDelete(ctx *fiber.Ctx) error {
	product, err := c.getProductByID(ctx)
	if err != nil {
		messages.SetErrorMessage(ctx, "falha ao identificar o produto: "+err.Error())
		return ctx.Redirect("/products")
	}
	return ctx.Render("products/delete", fiber.Map{"Object": product}, views.BaseLayout)
}

func (c *productController) Update(ctx *fiber.Ctx) error {
	product, err := c.getProductByID(ctx)
	if err != nil {
		messages.SetErrorMessage(ctx, "falha ao identificar o produto: "+err.Error())
		return ctx.Redirect("/products")
	}

	err = c.updateProductFromRequest(ctx, &product)
	if err != nil {
		messages.SetErrorMessage(ctx, "falha ao atualizar produto")
		return ctx.Render("products/product", fiber.Map{"Object": product}, views.BaseLayout)
	}

	messages.SetSuccessMessage(ctx, "produto atualizado com sucesso")
	return ctx.Redirect("/products/" + strconv.Itoa(int(product.ID)))

}

func (c *productController) Delete(ctx *fiber.Ctx) error {
	productID, err := c.getProductIdFromParam(ctx)
	if err != nil {
		messages.SetErrorMessage(ctx, "falha ao identificar o produto: "+err.Error())
		return ctx.Redirect("/products")
	}

	if err := c.productService.Delete(productID); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/products")
	}

	messages.SetSuccessMessage(ctx, "produto deletado com sucesso")
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
		return err
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
	return helpers.ParseStringToID(ctx.Params("id"))
}

func (c *productController) getProductFilterFromQueryParams(ctx *fiber.Ctx) *models.ProductFilter {
	return models.NewProductFilter(ctx.Query("q", ""), ctx.QueryInt("page"), ctx.QueryInt("limit"))
}

func (c *productController) getProductByID(ctx *fiber.Ctx) (models.Product, error) {
	productID, err := helpers.ParseStringToID(ctx.Params("id"))
	if err != nil {
		return models.Product{}, err
	}
	return c.productService.FindById(productID)
}
