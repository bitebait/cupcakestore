package routers

import (
	"github.com/gofiber/fiber/v2"
)

type Routers interface {
	InstallRouters(app *fiber.App)
}

func InstallRouters(app *fiber.App) {
	setup(app,
		NewAuthRouter(),
		NewUserRouter(),
		NewProfileRouter(),
		NewProductRouter(),
		NewStockRouter(),
		NewStoreConfigRouter(),
		NewStoreRouter(),
		NewHomeRouter(),
		NewShoppingCartRouter(),
		NewOrderRouter(),
		NewDashboardRouter(),
	)
}

func setup(app *fiber.App, routes ...Routers) {
	for _, r := range routes {
		r.InstallRouters(app)
	}
}
