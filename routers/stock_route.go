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

type StockRouter struct {
	stockController controllers.StockController
}

func NewStockRouter() *StockRouter {
	stockRepository := repositories.NewStockRepository(database.DB)

	stockService := services.NewStockService(stockRepository)

	stockController := controllers.NewStockController(stockService)

	return &StockRouter{
		stockController: stockController,
	}
}

func (r *StockRouter) InstallRouters(app *fiber.App) {
	stock := app.Group("/stock", cors.New()).Use(middlewares.LoginAndStaffRequired())

	stock.Get("/create", r.stockController.RenderCreate)
	stock.Post("/create", r.stockController.HandlerCreate)

}