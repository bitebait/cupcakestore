package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

const (
	usersPath = "/users"
	rootPath  = "/"
)

func selectLayout(isStaff, isUserProfile bool) string {
	if isStaff {
		return views.BaseLayout
	}
	if isUserProfile {
		return views.StoreLayout
	}
	return views.BaseLayout
}

func selectRedirectPath(isStaff bool) string {
	if isStaff {
		return usersPath
	}
	return rootPath
}

func getUserID(ctx *fiber.Ctx) uint {
	profile, ok := ctx.Locals("Profile").(*models.Profile)
	if !ok {
		return 0
	}
	return profile.ID
}

func renderErrorMessage(err error, action string) error {
	errorMessage := "Houve um erro ao " + action
	if err != nil {
		errorMessage += ": " + err.Error()
	}
	return fiber.NewError(fiber.StatusBadRequest, errorMessage)
}
