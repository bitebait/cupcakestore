package views

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2"
)

func RenderTemplate(ctx *fiber.Ctx, templateName string, obj interface{}, baseLayout ...string) error {
	user, ok := ctx.Locals("user").(*models.User)
	if !ok {
		user = nil
	}
	response := &models.Response{
		Error:   false,
		Message: "",
		Object:  obj,
		User:    user,
	}
	return ctx.Render(templateName, response, baseLayout...)
}

func RenderTemplateWithMessage(ctx *fiber.Ctx, templateName string, isError bool, message string, obj interface{}, baseLayout string) error {
	user, ok := ctx.Locals("user").(*models.User)
	if !ok {
		user = nil
	}

	response := &models.Response{
		Error:   isError,
		Message: message,
		Object:  obj,
		User:    user,
	}
	return ctx.Render(templateName, response, baseLayout)
}
