package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func CacheControl() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "public, max-age=120")
		return c.Next()
	}
}
