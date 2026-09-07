package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"time"
)

func RegisterAuthRoutes(app *fiber.App, controller controllers.AuthController) {
	auth := app.Group("/auth")

	auth.Get("/login", controller.RenderLogin)
	auth.Post("/login", limiter.New(limiter.Config{Max: 10, Expiration: time.Minute}), controller.Login)
	auth.Get("/register", controller.RenderRegister)
	auth.Post("/register", limiter.New(limiter.Config{Max: 5, Expiration: time.Minute}), controller.Register)
	auth.Post("/logout", middlewares.LoginRequired(), controller.Logout)
}
