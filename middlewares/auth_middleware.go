package middlewares

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := session.Store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
		}

		if profile := sess.Get("profile"); profile != nil {
			c.Locals("profile", profile.(*models.Profile))
			if !profile.(*models.Profile).User.IsActive {
				return c.Redirect("/auth/logout")
			}
			return c.Next()
		}

		switch c.Path() {
		case "/auth/login", "/store", "/auth/register":
			return c.Next()
		default:
			return c.Redirect("/auth/login")
		}
	}
}

func LoginAndStaffRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		profile, ok := c.Locals("profile").(*models.Profile)
		if !(ok && profile != nil && profile.User.IsStaff && profile.User.IsActive) {
			return c.Redirect("/auth/logout")
		}

		return c.Next()
	}
}
