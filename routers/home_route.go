package routers

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterHomeRoutes(app *fiber.App, redirectAfterLogin string) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.Redirect().Status(fiber.StatusMovedPermanently).To(redirectAfterLogin)
	})
}
