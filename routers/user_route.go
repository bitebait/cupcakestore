package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type UserRouter struct{}

func NewUserRouter() *UserRouter {
	return &UserRouter{}
}

func (ur *UserRouter) InstallRouters(app *fiber.App) {
	// Initialize repositories
	userRepository := repositories.NewUserRepository(database.DB)

	// Initialize services with repositories
	userService := services.NewUserService(userRepository)

	// Initialize controllers with services
	userController := controllers.NewUserController(userService)

	user := app.Group("/user", cors.New())
	user.Get("/create", userController.RenderCreate)
	user.Post("/create", userController.HandleCreate)
	user.Get("/list", userController.RenderList)
	user.Get("/detail/:id", userController.RenderDetail).Name("userDetail")
	user.Post("/update/:id", userController.HandleUpdate)
}
