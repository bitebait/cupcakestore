package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/gofiber/fiber/v3"
)

func RegisterDashboardRoutes(app *fiber.App, controller controllers.DashboardController) {
	dashboard := app.Group("/dashboard").Use(middlewares.LoginAndStaffRequired())
	dashboard.Get("/", controller.RenderDashboard)
}
