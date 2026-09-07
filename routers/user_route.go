package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(app *fiber.App, controller controllers.UserController) {
	userGroup := app.Group("/users").Use(middlewares.LoginRequired())
	userGroup.Get("/user/:id", controller.RenderUser)
	userGroup.Post("/user/update/:id", controller.Update)

	adminGroup := app.Group("/users").Use(middlewares.LoginAndStaffRequired())
	adminGroup.Get("/create", controller.RenderCreate)
	adminGroup.Post("/create", controller.Create)
	adminGroup.Get("/delete/:id", controller.RenderDelete)
	adminGroup.Post("/delete/:id", controller.Delete)
	adminGroup.Get("/", controller.RenderUsers)

}
