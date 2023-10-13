package controllers

import (
	"strconv"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	RenderCreate(ctx *fiber.Ctx) error
	HandleCreate(ctx *fiber.Ctx) error
	RenderList(ctx *fiber.Ctx) error
	RenderDetail(ctx *fiber.Ctx) error
	HandleUpdate(ctx *fiber.Ctx) error
	RenderDelete(ctx *fiber.Ctx) error
	HandleDelete(ctx *fiber.Ctx) error
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{userService: userService}
}

func (c *userController) RenderCreate(ctx *fiber.Ctx) error {
	return ctx.Render("user/create", fiber.Map{}, "layouts/base")
}

func (c *userController) HandleCreate(ctx *fiber.Ctx) error {
	user := &models.User{}

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Render("user/create", fiber.Map{"error": "Dados do usuário inválidos: " + err.Error()}, "layouts/base")
	}

	if validationErr := user.Validate(); validationErr != nil {
		return ctx.Render("user/create", fiber.Map{"error": "Dados do usuário inválidos: " + validationErr.Error()}, "layouts/base")
	}

	if err := c.userService.Create(user); err != nil {
		return ctx.Render("user/create", fiber.Map{"error": "Erro ao criar usuário: " + err.Error()}, "layouts/base")
	}

	return ctx.Redirect("/user/list")
}

func (c *userController) RenderList(ctx *fiber.Ctx) error {
	users := c.userService.List()
	return ctx.Render("user/list", fiber.Map{"Users": users}, "layouts/base")
}

func (c *userController) RenderDetail(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return ctx.Redirect("/user/list")
	}

	user, err := c.userService.FindById(uint(userID))
	if err != nil {
		return ctx.Redirect("/user/list")
	}

	return ctx.Render("user/detail", fiber.Map{"User": user}, "layouts/base")
}

func (c *userController) HandleUpdate(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Redirect("/user/list")
	}

	user, err := c.userService.FindById(uint(id))
	if err != nil {
		return ctx.Redirect("/user/list")
	}

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Render("user/detail", fiber.Map{"User": user, "error": err.Error()}, "layouts/base")
	}

	oldPassword := ctx.FormValue("oldPassword")
	newPassword := ctx.FormValue("newPassword")
	if oldPassword != "" && newPassword != "" {
		if err := user.UpdatePassword(oldPassword, newPassword); err != nil {
			return ctx.Render("user/detail", fiber.Map{"User": user, "error": "Não foi possível atualizar a senha do usuário. Por favor, verifique os dados."}, "layouts/base")
		}
	}

	if err := c.userService.Update(user); err != nil {
		return ctx.Render("user/detail", fiber.Map{"User": user, "error": "Falha ao atualizar usuário."}, "layouts/base")
	}

	return ctx.Redirect("/user/list")
}

func (c *userController) RenderDelete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Redirect("/user/list")
	}

	user, err := c.userService.FindById(uint(id))
	if err != nil {
		return ctx.Redirect("/user/list")
	}
	return ctx.Render("user/delete", fiber.Map{"User": user}, "layouts/base")
}

func (c *userController) HandleDelete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Redirect("/user/list")
	}

	err = c.userService.Delete(uint(id))
	if err != nil {
		return ctx.Redirect("/user/list")
	}

	return ctx.Redirect("/user/list")
}
