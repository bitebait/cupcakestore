package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

const (
	baseLayout  = "layouts/base"
	storeLayout = "layouts/store"
	userPath    = "/users"
	rootPath    = "/"
)

func selectLayout(isStaff, isUserProfile bool) string {
	if isStaff {
		return baseLayout
	}
	if isUserProfile {
		return storeLayout
	}
	return ""
}

func selectRedirectPath(isStaff bool) string {
	if isStaff {
		return userPath
	}
	return rootPath
}

func getUserID(ctx *fiber.Ctx) uint {
	return ctx.Locals("profile").(*models.Profile).ID
}

func renderErrorMessage(ctx *fiber.Ctx, err error, action, templateName string) error {
	errorMessage := "Houve um erro ao " + action + ": " + err.Error()
	return views.Render(ctx, templateName, fiber.Map{}, errorMessage, selectLayout(false, true))
}
