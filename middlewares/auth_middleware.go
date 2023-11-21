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
			return c.Status(fiber.StatusInternalServerError).SendString("Erro interno do servidor")
		}

		if profile := sess.Get("profile"); profile != nil {
			c.Locals("profile", profile.(*models.Profile))
		} else if c.Path() != "/auth/login" && c.Path() != "/store" && c.Path() != "/auth/register" {
			return c.Redirect("/auth/login")
		}

		return c.Next()
	}
}

func LoginAndStaffRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		profile, ok := c.Locals("profile").(*models.Profile)
		if !ok || profile == nil || !profile.User.IsStaff || !profile.User.IsActive {
			return redirectToLogout(c)
		}

		return c.Next()
	}
}

func redirectToLogout(c *fiber.Ctx) error {
	return c.Redirect("/auth/logout")
}
