package views

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, templateName string, obj interface{}, message string, baseLayout ...string) error {
	user := getUserFromContext(ctx)
	response := createResponse(message, obj, user)
	if message != "" {
		response.Error = true
	}
	return ctx.Render(templateName, response, baseLayout...)
}

func getUserFromContext(ctx *fiber.Ctx) *models.User {
	user, ok := ctx.Locals("user").(*models.User)
	if !ok {
		return nil
	}
	return user
}

func createResponse(message string, obj interface{}, user *models.User) *models.Response {
	return &models.Response{
		Error:   message != "",
		Message: message,
		Object:  obj,
		User:    user,
	}
}
