package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	storeConfig := app.Group("/config", cors.New()).Use(middlewares.LoginAndStaffRequired())
	storeConfig.Get("/", r.storeConfigController.RenderStoreConfig)
}
