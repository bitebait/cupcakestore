package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/gofiber/fiber/v3"
)

func RegisterProfileRoutes(app *fiber.App, controller controllers.ProfileController) {
	profile := app.Group("/profile").Use(middlewares.LoginRequired())

	profile.Get("/:id", controller.RenderProfile)
	profile.Post("/update/:id", controller.Update)
}
