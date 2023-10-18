package views

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2"
)

func RenderTemplate(ctx *fiber.Ctx, templateName string, obj interface{}, baseLayout ...string) error {
	user := getUserFromContext(ctx)
	response := createResponse(false, "", obj, user)
	return ctx.Render(templateName, response, baseLayout...)
}

func RenderTemplateWithError(ctx *fiber.Ctx, templateName string, obj interface{}, message string, baseLayout ...string) error {
	user := getUserFromContext(ctx)
	response := createResponse(true, message, obj, user)
	return ctx.Render(templateName, response, baseLayout...)
}

func RenderTemplateWithSuccess(ctx *fiber.Ctx, templateName string, obj interface{}, message string, baseLayout ...string) error {
	user := getUserFromContext(ctx)
	response := createResponse(false, message, obj, user)
	return ctx.Render(templateName, response, baseLayout...)
}

func getUserFromContext(ctx *fiber.Ctx) *models.User {
	user, ok := ctx.Locals("user").(*models.User)
	if !ok {
		return nil
	}
	return user
}

func createResponse(err bool, message string, obj interface{}, user *models.User) *models.Response {
	return &models.Response{
		Error:   err,
		Message: message,
		Object:  obj,
		User:    user,
	}
}
