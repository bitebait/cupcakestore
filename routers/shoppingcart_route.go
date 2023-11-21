package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
)

type ShoppingCartRouter struct {
	shoppingCartController controllers.ShoppingCartController
}

func NewShoppingCartRouter() *ShoppingCartRouter {
	shoppingCartRepository := repositories.NewShoppingCartRepository(database.DB)
	shoppingCartItemRepository := repositories.NewShoppingCartItemRepository(database.DB)
	shoppingCartItemService := services.NewShoppingCartItemService(shoppingCartItemRepository)
	shoppingCartService := services.NewShoppingCartService(shoppingCartRepository, shoppingCartItemService)

	shoppingCartController := controllers.NewShoppingCartController(shoppingCartService)
	return &ShoppingCartRouter{
		shoppingCartController: shoppingCartController,
	}
}

func (r *ShoppingCartRouter) InstallRouters(app *fiber.App) {
	cart := app.Group("/cart")
	cart.Get("/", r.shoppingCartController.RenderShoppingCart)
	cart.Post("/", r.shoppingCartController.AddShoppingCartItem)
	cart.Get("/remove/:id", r.shoppingCartController.RemoveFromCart)

}
