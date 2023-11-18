package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type UserRouter struct {
	userController controllers.UserController
}

func NewUserRouter() *UserRouter {
	userRepository := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	return &UserRouter{
		userController: userController,
	}
}

func (r *UserRouter) InstallRouters(app *fiber.App) {
	userGroup := app.Group("/users", cors.New())
	userGroup.Get("/:id", r.userController.RenderUser)
	userGroup.Post("/update/:id", r.userController.Update)

	adminGroup := app.Group("/users", cors.New(), middlewares.LoginAndStaffRequired())
	adminGroup.Get("/create", r.userController.RenderCreate)
	adminGroup.Post("/create", r.userController.Create)
	adminGroup.Get("/", r.userController.RenderUsers)
	adminGroup.Get("/delete/:id", r.userController.RenderDelete)
	adminGroup.Post("/delete/:id", r.userController.Delete)
}
