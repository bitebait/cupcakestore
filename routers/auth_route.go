package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"time"
)

type AuthRouter struct {
	authController controllers.AuthController
}

func NewAuthRouter() *AuthRouter {
	// Initialize repositories
	userRepository := repositories.NewUserRepository(database.DB)
	profileRepository := repositories.NewProfileRepository(database.DB)

	// Initialize services with repositories
	userService := services.NewUserService(userRepository)
	profileService := services.NewProfileService(profileRepository)
	authService := services.NewAuthService(userService, profileService)

	// Initialize controllers with services
	authController := controllers.NewAuthController(authService)

	return &AuthRouter{
		authController: authController,
	}
}

func (r *AuthRouter) InstallRouters(app *fiber.App) {
	auth := app.Group("/auth")

	auth.Get("/login", r.authController.RenderLogin)
	auth.Post("/login", limiter.New(limiter.Config{Max: 10, Expiration: time.Minute}), r.authController.Login)
	auth.Get("/register", r.authController.RenderRegister)
	auth.Post("/register", limiter.New(limiter.Config{Max: 5, Expiration: time.Minute}), r.authController.Register)
	auth.Post("/logout", middlewares.LoginRequired(), r.authController.Logout)
}
