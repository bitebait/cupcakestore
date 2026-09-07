package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/gofiber/fiber/v3"
)

func RegisterStockRoutes(app *fiber.App, controller controllers.StockController) {
	stock := app.Group("/stock").Use(middlewares.LoginAndStaffRequired())

	stock.Get("/create", controller.RenderCreate)
	stock.Post("/create", controller.Create)
	stock.Get("/", controller.RenderStocks)
	stock.Get("/:id", controller.RenderStock)
}
