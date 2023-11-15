package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/utils"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type ProfileController interface {
	HandlerUpdate(ctx *fiber.Ctx) error
	RenderProfile(ctx *fiber.Ctx) error
}

type profileController struct {
	profileService services.ProfileService
}

func NewProfileController(p services.ProfileService) ProfileController {
	return &profileController{
		profileService: p,
	}
}

func (c *profileController) RenderProfile(ctx *fiber.Ctx) error {
	userID, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return ctx.Redirect("/users")
	}

	profile, err := c.profileService.FindByUserId(userID)
	if err != nil {
		return ctx.Redirect("/users")
	}

	return views.Render(ctx, "profile/user-profile", profile, "", baseLayout)
}

func (c *profileController) HandlerUpdate(ctx *fiber.Ctx) error {
	id, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return c.redirectToUsers(ctx)
	}

	profile, err := c.profileService.FindByUserId(id)
	if err != nil {
		return c.redirectToUsers(ctx)
	}

	if err := ctx.BodyParser(&profile); err != nil {
		return views.Render(ctx, "users/user", profile, err.Error(), baseLayout)
	}

	if err := c.profileService.Update(&profile); err != nil {
		return views.Render(ctx, "users/user", profile, "Falha ao atualizar perfil do usuário.", baseLayout)
	}

	return ctx.Redirect("/users")
}

func (c *profileController) redirectToUsers(ctx *fiber.Ctx) error {
	return ctx.Redirect("/users")
}
