package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type AuthRouter struct {
	authController controllers.AuthController
}

func NewAuthRouter() *AuthRouter {
	// Initialize repositories
	userRepository := repositories.NewUserRepository(database.DB)

	// Initialize services with repositories
	authService := services.NewAuthService(userRepository)

	// Initialize controllers with services
	authController := controllers.NewAuthController(authService, session.Store)

	return &AuthRouter{
		authController: authController,
	}
}

func (r *AuthRouter) InstallRouters(app *fiber.App) {
	auth := app.Group("/auth", cors.New())

	auth.Get("/login", r.authController.RenderLogin)
	auth.Post("/login", r.authController.HandlerLogin)
}
