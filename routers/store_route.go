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

type StoreRouter struct {
	storeController controllers.StoreController
}

func NewStoreRouter() *StoreRouter {
	// Initialize repositories
	productRepository := repositories.NewProductRepository(database.DB)

	// Initialize services with repositories
	productService := services.NewProductService(productRepository)

	// Initialize controllers with services
	storeController := controllers.NewStoreController(productService)
	return &StoreRouter{
		storeController: storeController,
	}
}

func (r *StoreRouter) InstallRouters(app *fiber.App) {
	store := app.Group("/store", cors.New()).Use(middlewares.Auth())
	store.Get("/", r.storeController.RenderStore)
}
