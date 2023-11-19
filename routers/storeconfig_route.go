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
	storeConfigRepository := repositories.NewStoreConfigRepository(database.DB)
	storeConfigService := services.NewStoreConfigService(storeConfigRepository)
	storeConfigController := controllers.NewStoreConfigController(storeConfigService)

	return &StoreConfigRouter{
		storeConfigController: storeConfigController,
	}
}

func (r *StoreConfigRouter) InstallRouters(app *fiber.App) {
	storeConfig := app.Group("/config").Use(middlewares.LoginAndStaffRequired())
	storeConfig.Get("/address", r.storeConfigController.RenderStoreConfigAddress)
	storeConfig.Get("/delivery", r.storeConfigController.RenderStoreConfigDelivery)
	storeConfig.Get("/payment", r.storeConfigController.RenderStoreConfigPayment)
	storeConfig.Get("/pix", r.storeConfigController.RenderStoreConfigPix)
	storeConfig.Post("/", r.storeConfigController.Update)
}
