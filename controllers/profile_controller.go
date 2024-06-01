package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/utils"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type ProfileController interface {
	Update(ctx *fiber.Ctx) error
	RenderProfile(ctx *fiber.Ctx) error
}

type profileController struct {
	profileService services.ProfileService
}

func NewProfileController(profileService services.ProfileService) ProfileController {
	return &profileController{profileService: profileService}
}

func (c *profileController) RenderProfile(ctx *fiber.Ctx) error {
	profile, err := c.getProfile(ctx)
	if err != nil {
		return err
	}

	userSess, err := c.getUserSession(ctx)
	if err != nil {
		return err
	}

	layout := selectLayout(userSess.User.IsStaff, profile.UserID == userSess.UserID)
	if layout == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	return views.Render(ctx, "profile/user-profile", profile, layout)
}

func (c *profileController) Update(ctx *fiber.Ctx) error {
	profile, err := c.getProfile(ctx)
	if err != nil {
		return err
	}

	if err = ctx.BodyParser(&profile); err != nil {
		return views.RenderError(ctx, "profile/user-profile", profile, err.Error())
	}

	userSess, err := c.getUserSession(ctx)
	if err != nil {
		return err
	}

	if !userSess.User.IsStaff && profile.UserID != userSess.UserID {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	if err = c.profileService.Update(&profile); err != nil {
		return views.RenderError(ctx, "profile/user-profile", profile,
			"Falha ao atualizar perfil do usuário.", selectLayout(userSess.User.IsStaff, profile.UserID == userSess.UserID))
	}

	redirectPath := selectRedirectPath(userSess.User.IsStaff)
	return ctx.Redirect(redirectPath)
}

func (c *profileController) getProfile(ctx *fiber.Ctx) (models.Profile, error) {
	userID, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return models.Profile{}, ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return c.profileService.FindByUserId(userID)
}

func (c *profileController) getUserSession(ctx *fiber.Ctx) (*models.Profile, error) {
	userSess, ok := ctx.Locals("Profile").(*models.Profile)
	if !ok || userSess == nil {
		return nil, fiber.ErrUnauthorized
	}
	return userSess, nil
}
