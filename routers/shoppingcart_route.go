package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/gofiber/fiber/v3"
)

func RegisterShoppingCartRoutes(app *fiber.App, controller controllers.ShoppingCartController) {
	cart := app.Group("/cart").Use(middlewares.LoginRequired())
	cart.Get("/", controller.RenderShoppingCart)
	cart.Post("/", controller.AddShoppingCartItem)
	cart.Get("/count", controller.CountShoppingCart)
	cart.Post("/remove/:id", controller.RemoveFromCart)
	cart.Post("/items/:id/quantity", controller.SetItemQuantity)
}
