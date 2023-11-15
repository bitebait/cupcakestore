package middlewares

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
)

func LoginAndStaffRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := session.Store.Get(c)
		if err != nil {
			panic(err)
		}

		profile, ok := sess.Get("profile").(*models.Profile)
		if !ok || profile == nil {
			return redirectToLogout(c)
		}

		if !profile.User.IsStaff || !profile.User.IsActive {
			return redirectToLogout(c)
		}

		c.Locals("profile", profile)

		return c.Next()
	}
}

func redirectToLogout(c *fiber.Ctx) error {
	return c.Redirect("/auth/logout")
}
