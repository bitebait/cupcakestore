package views

import (
	"github.com/gofiber/fiber/v2"
)

const (
	BaseLayout  = "layouts/base"
	StoreLayout = "layouts/store"
	ErrorKey    = "Error"
	MessageKey  = "Message"
	ObjectKey   = "Object"
)

func setContextLocals(ctx *fiber.Ctx, message string) {
	ctx.Locals(ErrorKey, message != "")
	ctx.Locals(MessageKey, message)
}

func renderTemplate(ctx *fiber.Ctx, templateName string, data interface{}, message string, baseLayout ...string) error {
	setContextLocals(ctx, message)
	return ctx.Render(templateName, fiber.Map{ObjectKey: data}, baseLayout...)
}

func Render(ctx *fiber.Ctx, template string, data interface{}, baseLayout ...string) error {
	return renderTemplate(ctx, template, data, "", baseLayout...)
}

func RenderError(ctx *fiber.Ctx, template string, data interface{}, message string, baseLayout ...string) error {
	return renderTemplate(ctx, template, data, message, baseLayout...)
}
