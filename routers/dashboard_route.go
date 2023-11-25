package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
)

type DashboardRouter struct {
	dashboardController controllers.DashboardController
}

func NewDashboardRouter() *DashboardRouter {
	// Initialize repositories
	dashboardRepository := repositories.NewDashboardRepository(database.DB)

	// Initialize services with repositories
	dashboardService := services.NewDashboardService(dashboardRepository)

	// Initialize controllers with services
	dashboardController := controllers.NewDashboardController(dashboardService)

	return &DashboardRouter{
		dashboardController: dashboardController,
	}
}

func (r *DashboardRouter) InstallRouters(app *fiber.App) {
	dashboard := app.Group("/dashboard").Use(middlewares.LoginAndStaffRequired())
	dashboard.Get("/", r.dashboardController.RenderDashboard)
}
