package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ProfileRouter struct {
	profileController controllers.ProfileController
	profileService    services.ProfileService
}

func NewProfileRouter() *ProfileRouter {
	profileRepository := repositories.NewProfileRepository(database.DB)
	profileService := services.NewProfileService(profileRepository)
	profileController := controllers.NewProfileController(profileService)

	return &ProfileRouter{
		profileController: profileController,
		profileService:    profileService,
	}
}

func (r *ProfileRouter) InstallRouters(app *fiber.App) {
	profile := app.Group("/profile", cors.New())
	profile.Post("/update/:id", r.profileController.HandlerUpdate)
}
