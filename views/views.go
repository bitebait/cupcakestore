package views

import (
	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, templateName string, data interface{}, message string, baseLayout ...string) error {
	ctx.Locals("Error", message != "")
	ctx.Locals("Message", message)

	renderData := fiber.Map{"Object": data}

	return ctx.Render(templateName, renderData, baseLayout...)
}
