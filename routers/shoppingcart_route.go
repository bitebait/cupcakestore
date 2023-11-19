package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ShoppingCartRouter struct {
	shoppingCartController controllers.ShoppingCartController
}

func NewShoppingCartRouter() *ShoppingCartRouter {
	shoppginCartRepository := repositories.NewShoppingCartRepository(database.DB)
	shoppginCartService := services.NewShoppingCartService(shoppginCartRepository)
	shoppginCartController := controllers.NewShoppingCartController(shoppginCartService)
	return &ShoppingCartRouter{
		shoppingCartController: shoppginCartController,
	}
}

func (r *ShoppingCartRouter) InstallRouters(app *fiber.App) {
	cart := app.Group("/cart", cors.New())
	cart.Get("/", r.shoppingCartController.RenderShoppingCart)
}
