package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
)

type StoreConfigRouter struct {
	storeConfigController controllers.StoreConfigController
}

func NewStoreConfigRouter() *StoreConfigRouter {
	// Initialize repositories
	storeConfigRepository := repositories.NewStoreConfigRepository(database.DB)

	// Initialize services with repositories
	storeConfigService := services.NewStoreConfigService(storeConfigRepository)

	// Initialize controllers with services
	storeConfigController := controllers.NewStoreConfigController(storeConfigService)

	return &StoreConfigRouter{
		storeConfigController: storeConfigController,
	}
}

func (r *StoreConfigRouter) InstallRouters(app *fiber.App) {
	storeConfig := app.Group("/config").Use(middlewares.LoginAndStaffRequired())
	storeConfig.Get("/address", func(ctx *fiber.Ctx) error {
		return r.storeConfigController.RenderStoreConfig(ctx, "address")
	})
	storeConfig.Get("/delivery", func(ctx *fiber.Ctx) error {
		return r.storeConfigController.RenderStoreConfig(ctx, "delivery")
	})
	storeConfig.Get("/payment", func(ctx *fiber.Ctx) error {
		return r.storeConfigController.RenderStoreConfig(ctx, "payment")
	})
	storeConfig.Get("/pix", func(ctx *fiber.Ctx) error {
		return r.storeConfigController.RenderStoreConfig(ctx, "pix")
	})
	storeConfig.Post("/", r.storeConfigController.Update)
}
