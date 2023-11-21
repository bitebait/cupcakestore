package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
)

type OrderRouter struct {
	orderController controllers.OrderController
}

func NewOrderRouter() *OrderRouter {
	orderRepository := repositories.NewOrderRepository(database.DB)
	orderService := services.NewOrderService(orderRepository)
	orderController := controllers.NewOrderController(orderService)

	return &OrderRouter{
		orderController: orderController,
	}
}

func (r *OrderRouter) InstallRouters(app *fiber.App) {
	order := app.Group("/orders")
	order.Get("/", r.orderController.RenderOrders)
}
