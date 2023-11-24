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
	switch {
	case isStaff:
		return baseLayout
	case isUserProfile:
		return storeLayout
	default:
		return baseLayout
	}
}

func selectRedirectPath(isStaff bool) string {
	if isStaff {
		return userPath
	}
	return rootPath
}

func getUserID(ctx *fiber.Ctx) uint {
	profile, ok := ctx.Locals("profile").(*models.Profile)
	if !ok {
		return 0
	}
	return profile.ID
}

func renderErrorMessage(ctx *fiber.Ctx, err error, action string) error {
	errorMessage := "Houve um erro ao " + action
	if err != nil {
		errorMessage += ": " + err.Error()
	}
	return ctx.Status(fiber.StatusInternalServerError).SendString(errorMessage)
}
