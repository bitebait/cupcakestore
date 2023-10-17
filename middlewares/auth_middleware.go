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
			return c.Redirect("/auth/logout")
		}

		user, err := userService.FindByUsername(username.(string))
		if err != nil {
			return c.Redirect("/auth/logout")
		}

		if user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Usuário não encontrado")
		}

		if !user.IsStaff {
			return fiber.NewError(fiber.StatusUnauthorized, "Você não possui autorização de administrador")
		}

		c.Locals("user", user)

		return c.Next()
	}
}
