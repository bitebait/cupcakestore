package views

import (
	"github.com/gofiber/fiber/v2"
)

const (
	BaseLayout  = "layouts/base"
	StoreLayout = "layouts/store"
)

func render(ctx *fiber.Ctx, templateName string, data interface{}, message string, baseLayout ...string) error {
	ctx.Locals("Error", message != "")
	ctx.Locals("Message", message)
	return ctx.Render(templateName, fiber.Map{"Object": data}, baseLayout...)
}

func Render(ctx *fiber.Ctx, template string, data interface{}, baseLayout ...string) error {
	return render(ctx, template, data, "", baseLayout...)
}

func RenderError(ctx *fiber.Ctx, template string, data interface{}, message string, baseLayout ...string) error {
	return render(ctx, template, data, message, baseLayout...)
}
