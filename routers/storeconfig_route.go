package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/gofiber/fiber/v3"
)

func RegisterStoreConfigRoutes(app *fiber.App, controller controllers.StoreConfigController) {
	storeConfig := app.Group("/config").Use(middlewares.LoginAndStaffRequired())
	storeConfig.Get("/address", func(ctx fiber.Ctx) error {
		return controller.RenderStoreConfig(ctx, "address")
	})
	storeConfig.Get("/delivery", func(ctx fiber.Ctx) error {
		return controller.RenderStoreConfig(ctx, "delivery")
	})
	storeConfig.Get("/payment", func(ctx fiber.Ctx) error {
		return controller.RenderStoreConfig(ctx, "payment")
	})
	storeConfig.Get("/pix", func(ctx fiber.Ctx) error {
		return controller.RenderStoreConfig(ctx, "pix")
	})
	storeConfig.Post("/", controller.Update)
}
