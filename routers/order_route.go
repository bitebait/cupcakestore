package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/gofiber/fiber/v3"
)

func RegisterOrderRoutes(app *fiber.App, controller controllers.OrderController) {
	order := app.Group("/orders").Use(middlewares.LoginRequired())
	order.Get("/checkout/:id", controller.Checkout)
	order.Post("/checkout/:id", controller.Checkout)
	order.Post("/payment/:id", controller.Payment)
	order.Get("/payment/:id", controller.Payment)
	order.Get("/cancel/:id", controller.RenderCancel)
	order.Post("/cancel/:id", controller.Cancel)
	order.Get("/", controller.RenderAllOrders)
	order.Get("/order/:id", controller.RenderOrder)
	order.Post("/order/:id", middlewares.LoginAndStaffRequired(), controller.Update)
}
