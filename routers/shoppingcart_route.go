package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ShoppingCartRouter struct {
	shoppingCartController controllers.ShoppingCartController
}

func NewShoppingCartRouter() *ShoppingCartRouter {
	// Initialize controllers with services
	shoppingCartController := controllers.NewShoppingCartController()
	return &ShoppingCartRouter{
		shoppingCartController: shoppingCartController,
	}
}

func (r *ShoppingCartRouter) InstallRouters(app *fiber.App) {
	cart := app.Group("/cart", cors.New())
	cart.Get("/", r.shoppingCartController.RenderShoppingCart)
}
