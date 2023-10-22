package routers

import (
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ProfileRouter struct {
	profileController controllers.ProfileController
	profileService    services.ProfileService
	templateRenderer  views.TemplateRenderer
}

func NewProfileRouter() *ProfileRouter {
	templateRenderer := views.NewTemplateRenderer()
	profileRepository := repositories.NewProfileRepository(database.DB)
	profileService := services.NewProfileService(profileRepository)
	profileController := controllers.NewProfileController(profileService, templateRenderer)

	return &ProfileRouter{
		profileController: profileController,
		profileService:    profileService,
		templateRenderer:  templateRenderer,
	}
}

func (r *ProfileRouter) InstallRouters(app *fiber.App) {
	profile := app.Group("/profile", cors.New())
	profile.Post("/update/:id", r.profileController.HandlerUpdate)
}
