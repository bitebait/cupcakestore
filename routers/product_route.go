package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
)

type ProductRouter struct {
	productController controllers.ProductController
}

func NewProductRouter() *ProductRouter {
	productRepository := repositories.NewProductRepository(database.DB)

	productService := services.NewProductService(productRepository)

	productController := controllers.NewProductController(productService)

	return &ProductRouter{
		productController: productController,
	}
}

func (r *ProductRouter) InstallRouters(app *fiber.App) {
	product := app.Group("/products").Use(middlewares.LoginAndStaffRequired())

	product.Get("/create", r.productController.RenderCreate)
	product.Post("/create", r.productController.Create)
	product.Get("/json", r.productController.JSONProducts)
	product.Get("/", r.productController.RenderProducts)
	product.Get("/:id", r.productController.RenderProduct)
	product.Post("/update/:id", r.productController.Update)
	product.Get("/delete/:id", r.productController.RenderDelete)
	product.Post("/delete/:id", r.productController.Delete)
}
