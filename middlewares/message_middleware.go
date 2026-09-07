package middlewares

import (
	"github.com/bitebait/cupcakestore/messages"
	"github.com/gofiber/fiber/v3"
)

func Message() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		msgs := messages.LoadMessages(ctx)
		ctx.Locals("Messages", msgs)

		return ctx.Next()
	}
}
