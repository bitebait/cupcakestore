package views

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, templateName string, obj interface{}, message string, baseLayout ...string) error {
	profile := getProfileFromContext(ctx)
	response := createResponse(message, &obj, profile)
	if message != "" {
		response.Error = true
	}
	return ctx.Render(templateName, response, baseLayout...)
}

func getProfileFromContext(ctx *fiber.Ctx) *models.Profile {
	profile, ok := ctx.Locals("profile").(*models.Profile)
	if !ok {
		return nil
	}
	return profile
}

func createResponse(message string, obj *interface{}, profile *models.Profile) *models.Response {
	return &models.Response{
		Error:   message != "",
		Message: message,
		Object:  *obj,
		Profile: profile,
	}
}
