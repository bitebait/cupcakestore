package controllers

import (
	"github.com/bitebait/cupcakestore/models"
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
	return baseLayout // Default layout if none is applicable
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

func renderErrorMessage(ctx *fiber.Ctx, err error, action string) error {
	errorMessage := "Houve um erro ao " + action
	if err != nil {
		errorMessage += ": " + err.Error()
	}
	return ctx.Status(fiber.StatusInternalServerError).SendString(errorMessage)
}
