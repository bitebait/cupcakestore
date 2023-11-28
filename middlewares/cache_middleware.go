package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func CacheControl() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-store, max-age=0")
		return c.Next()
	}
}
