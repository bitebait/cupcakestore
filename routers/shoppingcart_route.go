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
	shoppginCartRepository := repositories.NewShoppingCartRepository(database.DB)
	shoppginCartItemRepository := repositories.NewShoppingCartItemRepository(database.DB)
	shoppginCartItemService := services.NewShoppingCartItemService(shoppginCartItemRepository)
	shoppginCartService := services.NewShoppingCartService(shoppginCartRepository, shoppginCartItemService)
	shoppginCartController := controllers.NewShoppingCartController(shoppginCartService)
	return &ShoppingCartRouter{
		shoppingCartController: shoppginCartController,
	}
}

func (r *ShoppingCartRouter) InstallRouters(app *fiber.App) {
	cart := app.Group("/cart")
	cart.Get("/", r.shoppingCartController.RenderShoppingCart)
	cart.Post("/", r.shoppingCartController.AddShoppingCartItem)
	cart.Get("/remove/:id", r.shoppingCartController.RemoveFromCart) // Alterado de Post para Get
}
