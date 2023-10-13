package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type UserRouter struct {
	userController controllers.UserController
}

func NewUserRouter() *UserRouter {
	// Initialize repositories
	userRepository := repositories.NewUserRepository(database.DB)

	// Initialize services with repositories
	userService := services.NewUserService(userRepository)

	// Initialize controllers with services
	userController := controllers.NewUserController(userService)

	return &UserRouter{
		userController: userController,
	}
}

func (ur *UserRouter) InstallRouters(app *fiber.App) {
	user := app.Group("/user", cors.New())

	user.Get("/create", ur.userController.RenderCreate)
	user.Post("/create", ur.userController.HandleCreate)
	user.Get("/list", ur.userController.RenderList)
	user.Get("/detail/:id", ur.userController.RenderDetail).Name("userDetail")
	user.Post("/update/:id", ur.userController.HandleUpdate)
	user.Get("/delete/:id", ur.userController.RenderDelete)
	user.Post("/delete/:id", ur.userController.HandleDelete)
}
