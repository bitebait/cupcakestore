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
	profileService services.ProfileService
	userService    services.UserService
}

func NewProfileController(p services.ProfileService, u services.UserService) ProfileController {
	return &profileController{
		profileService: p,
		userService:    u,
	}
}

func (c *profileController) HandlerUpdate(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return c.redirectToUsers(ctx)
	}

	user, err := c.userService.FindById(uint(id))
	if err != nil {
		return c.redirectToUsers(ctx)
	}

	profile, err := c.profileService.FindByUserId(uint(id))
	if err != nil {
		return c.redirectToUsers(ctx)
	}

	data := fiber.Map{
		"User":    user,
		"Profile": profile,
	}

	if err := ctx.BodyParser(profile); err != nil {
		return views.Render(ctx, "users/user", data, err.Error(), baseLayout)
	}

	if err := c.profileService.Update(profile); err != nil {
		return views.Render(ctx, "users/user", data, "Falha ao atualizar perfil do usuário.", baseLayout)
	}

	return ctx.Redirect("/users")
}

func (c *profileController) redirectToUsers(ctx *fiber.Ctx) error {
	return ctx.Redirect("/users")
}
