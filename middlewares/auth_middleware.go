package middlewares

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
	"log"
)

func LoginAndStaffRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := session.Store.Get(c)
		if err != nil {
			panic(err)
		}

		user, ok := sess.Get("user").(*models.User)
		if !ok || user == nil {
			return redirectToLogout(c)
		}

		log.Println(user)

		if !user.IsStaff || !user.IsActive {
			return redirectToLogout(c)
		}

		c.Locals("user", user)

		return c.Next()
	}
}

func redirectToLogout(c *fiber.Ctx) error {
	return c.Redirect("/auth/logout")
}
