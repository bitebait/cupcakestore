package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/gofiber/fiber/v3"
)

func RegisterStoreRoutes(app *fiber.App, controller controllers.StoreController) {
	store := app.Group("/store")
	store.Get("/", controller.RenderStore)
}
