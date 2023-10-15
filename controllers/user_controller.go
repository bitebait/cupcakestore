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
	RenderUsers(ctx *fiber.Ctx) error
	RenderUser(ctx *fiber.Ctx) error
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

const baseLayout = "layouts/base"

func (c *userController) RenderCreate(ctx *fiber.Ctx) error {
	return ctx.Render("users/create", fiber.Map{}, baseLayout)
}

func (c *userController) HandleCreate(ctx *fiber.Ctx) error {
	user := &models.User{}

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Render("users/create", fiber.Map{"error": "Dados do usuário inválidos: " + err.Error()}, baseLayout)
	}

	if err := c.userService.Create(user); err != nil {
		return ctx.Render("users/create", fiber.Map{"error": "Erro ao criar usuário: " + err.Error()}, baseLayout)
	}

	return ctx.Redirect("/users")
}

func (c *userController) RenderUsers(ctx *fiber.Ctx) error {
	p := models.NewPagination(ctx.QueryInt("page"), ctx.QueryInt("limit"))
	users := c.userService.FindAll(p)

	return ctx.Render("users/users", fiber.Map{"Users": users, "Pagination": p}, baseLayout)
}

func (c *userController) RenderUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return ctx.Redirect("/users")
	}

	user, err := c.userService.FindById(uint(userID))
	if err != nil {
		return ctx.Redirect("/users")
	}

	return ctx.Render("users/user", fiber.Map{"User": user}, baseLayout)
}

func (c *userController) HandleUpdate(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Redirect("/users")
	}

	user, err := c.userService.FindById(uint(id))
	if err != nil {
		return ctx.Redirect("/users")
	}

	if err := ctx.BodyParser(user); err != nil {
		return ctx.Render("users/user", fiber.Map{"User": user, "error": err.Error()}, baseLayout)
	}

	oldPassword := ctx.FormValue("oldPassword")
	newPassword := ctx.FormValue("newPassword")
	if oldPassword != "" && newPassword != "" {
		if err := user.UpdatePassword(oldPassword, newPassword); err != nil {
			return ctx.Render("users/user", fiber.Map{"User": user, "error": "Não foi possível atualizar a senha do usuário. Por favor, verifique os dados."}, baseLayout)
		}
	}

	if ctx.FormValue("isStaff") != "on" {
		user.IsStaff = false
	}
	if ctx.FormValue("isActive") != "on" {
		user.IsActive = false
	}

	if err := c.userService.Update(user); err != nil {
		return ctx.Render("users/user", fiber.Map{"User": user, "error": "Falha ao atualizar usuário."}, baseLayout)
	}

	return ctx.Redirect("/users")
}

func (c *userController) RenderDelete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Redirect("/users")
	}

	user, err := c.userService.FindById(uint(id))
	if err != nil {
		return ctx.Redirect("/users")
	}
	return ctx.Render("users/delete", fiber.Map{"User": user}, baseLayout)
}

func (c *userController) HandleDelete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Redirect("/users")
	}

	err = c.userService.Delete(uint(id))
	if err != nil {
		return ctx.Redirect("/users")
	}

	return ctx.Redirect("/users")
}
