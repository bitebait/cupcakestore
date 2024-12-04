package routers

import (
	"github.com/bitebait/cupcakestore/config"
	"github.com/gofiber/fiber/v2"
)

type HomeRouter struct{}

func NewHomeRouter() *HomeRouter {
	return &HomeRouter{}
}

func (r *HomeRouter) InstallRouters(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect(config.Get().RedirectAfterLogin, fiber.StatusMovedPermanently)
	})
}
