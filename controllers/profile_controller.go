package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ProfileController interface {
	HandlerUpdate(ctx *fiber.Ctx) error
}

type profileController struct {
	profileService   services.ProfileService
	templateRenderer views.TemplateRenderer
}

func NewProfileController(p services.ProfileService, t views.TemplateRenderer) ProfileController {
	return &profileController{
		profileService:   p,
		templateRenderer: t,
	}
}

func (c *profileController) HandlerUpdate(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Redirect("/users")
	}

	profile, err := c.profileService.FindByUserId(uint(id))
	if err != nil {
		return ctx.Redirect("/users")
	}

	if err := ctx.BodyParser(profile); err != nil {
		return c.templateRenderer.Render(ctx, "users/user", nil, err.Error(), baseLayout)
	}

	if err := c.profileService.Update(profile); err != nil {
		return c.templateRenderer.Render(ctx, "users/user", nil,
			"Falha ao atualizar perfil do usuário.", baseLayout)
	}

	return ctx.Redirect("/users")
}
