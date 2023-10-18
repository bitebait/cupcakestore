package middlewares

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func LoginAndStaffRequired(store *session.Store, userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			panic(err)
		}

		username := sess.Get("username")
		if username == nil {
			return redirectToLogout(c)
		}

		user, err := userService.FindByUsername(username.(string))
		if err != nil {
			return redirectToLogout(c)
		}

		if user == nil || !user.IsStaff || !user.IsActive {
			return redirectToLogout(c)
		}

		c.Locals("user", user)

		return c.Next()
	}
}

func redirectToLogout(c *fiber.Ctx) error {
	return c.Redirect("/auth/logout")
}
