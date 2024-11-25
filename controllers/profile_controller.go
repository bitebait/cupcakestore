package controllers

import (
	"errors"
	"github.com/bitebait/cupcakestore/helpers"
	"github.com/bitebait/cupcakestore/messages"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
	profile, userSess, err := c.getAuthorizedProfileAndUser(ctx)
	if err != nil {
		messages.SetErrorMessage(ctx, "ocorreu um erro ao processar o perfil")
		return ctx.Redirect("/")
	}

	layout := selectLayout(userSess.User.IsStaff, profile.UserID == userSess.UserID)
	return ctx.Render("profile/user-profile", fiber.Map{"Object": profile}, layout)
}

func (c *profileController) Update(ctx *fiber.Ctx) error {
	profile, userSess, err := c.getAuthorizedProfileAndUser(ctx)
	if err != nil {
		messages.SetErrorMessage(ctx, "ocorreu um erro ao processar o perfil")
		return ctx.Redirect("/")
	}

	if err := ctx.BodyParser(&profile); err != nil {
		messages.SetErrorMessage(ctx, "ocorreu um erro ao processar o perfil")
		return ctx.Redirect("/")
	}

	if err = c.profileService.Update(&profile); err != nil {
		layout := selectLayout(userSess.User.IsStaff, profile.UserID == userSess.UserID)
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Render("profile/user-profile", fiber.Map{"Object": profile}, layout)
	}

	messages.SetSuccessMessage(ctx, "perfil atualizado com sucesso")
	return ctx.Redirect("/profile/" + strconv.Itoa(int(profile.ID)))
}

func (c *profileController) getAuthorizedProfileAndUser(ctx *fiber.Ctx) (models.Profile, *models.Profile, error) {
	profile, err := c.getProfile(ctx)
	if err != nil {
		return models.Profile{}, nil, err
	}

	userSess, err := c.getUserSession(ctx)
	if err != nil {
		return models.Profile{}, nil, err
	}

	if !userSess.User.IsStaff && profile.UserID != userSess.User.ID {
		return models.Profile{}, nil, errors.New("usuário não autorizado, por favor, efetue o login e tente novamente")
	}

	return profile, userSess, nil
}

func (c *profileController) getProfile(ctx *fiber.Ctx) (models.Profile, error) {
	userID, err := helpers.ParseStringToID(ctx.Params("id"))
	if err != nil {
		return models.Profile{}, errors.New("usuário não autorizado, por favor, efetue o login e tente novamente")
	}
	return c.profileService.FindByUserId(userID)
}

func (c *profileController) getUserSession(ctx *fiber.Ctx) (*models.Profile, error) {
	userSess, ok := ctx.Locals("Profile").(*models.Profile)
	if !ok || userSess == nil {
		return nil, errors.New("usuário não autorizado, por favor, efetue o login e tente novamente")
	}
	return userSess, nil
}
