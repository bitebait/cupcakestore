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
		redirectURL := config.Instance().GetEnvVar("REDIRECT_AFTER_LOGIN", "/")
		return c.Redirect(redirectURL, fiber.StatusMovedPermanently)
	})
}
