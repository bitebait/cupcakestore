package routers

import (
	"github.com/gofiber/fiber/v2"
)

type Routers interface {
	InstallRouters(app *fiber.App)
}

func InstallRouters(app *fiber.App) {
	setup(app, NewUserRouter())
}

func setup(app *fiber.App, routes ...Routers) {
	for _, r := range routes {
		r.InstallRouters(app)
	}
}
