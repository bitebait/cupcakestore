package views

import (
	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, templateName string, data interface{}, message string, baseLayout ...string) error {
	ctx.Locals("Error", message != "")
	ctx.Locals("Message", message)
	return ctx.Render(templateName, fiber.Map{"Object": data}, baseLayout...)
}
