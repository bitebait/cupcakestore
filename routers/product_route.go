package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/gofiber/fiber/v3"
)

func RegisterProductRoutes(app *fiber.App, controller controllers.ProductController) {
	product := app.Group("/products")
	product.Get("/details/:id", controller.RenderDetails)

	productAdmin := app.Group("/products").Use(middlewares.LoginAndStaffRequired())
	productAdmin.Get("/create", controller.RenderCreate)
	productAdmin.Post("/create", controller.Create)
	productAdmin.Get("/json", controller.JSONProducts)
	productAdmin.Post("/update/:id", controller.Update)
	productAdmin.Get("/delete/:id", controller.RenderDelete)
	productAdmin.Post("/delete/:id", controller.Delete)
	productAdmin.Get("/", controller.RenderProducts)
	productAdmin.Get("/:id", controller.RenderProduct)
}
