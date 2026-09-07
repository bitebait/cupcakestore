package routers

import (
	"github.com/bitebait/cupcakestore/config"
	"github.com/gofiber/fiber/v3"
)

type HomeRouter struct{}

func NewHomeRouter() *HomeRouter {
	return &HomeRouter{}
}

func (r *HomeRouter) InstallRouters(app *fiber.App) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.Redirect().Status(fiber.StatusMovedPermanently).To(config.Get().RedirectAfterLogin)
	})
}
