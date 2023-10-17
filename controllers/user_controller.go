package controllers

import (
	"fmt"
	"strconv"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	RenderCreate(ctx *fiber.Ctx) error
	HandlerCreate(ctx *fiber.Ctx) error
	RenderUsers(ctx *fiber.Ctx) error
	RenderUser(ctx *fiber.Ctx) error
	HandlerUpdate(ctx *fiber.Ctx) error
	RenderDelete(ctx *fiber.Ctx) error
	HandlerDelete(ctx *fiber.Ctx) error
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{userService: userService}
}

const baseLayout = "layouts/base"

func (c *userController) RenderCreate(ctx *fiber.Ctx) error {
	return views.RenderTemplate(ctx, "users/create", nil, baseLayout)
}

func (c *userController) HandlerCreate(ctx *fiber.Ctx) error {
	user := &models.User{}
	if err := ctx.BodyParser(user); err != nil {
		return views.RenderTemplateWithMessage(ctx, "users/create", true, "Dados do usuário inválidos: "+err.Error(), nil, baseLayout)
	}
	if err := c.userService.Create(user); err != nil {
		return views.RenderTemplateWithMessage(ctx, "users/create", true, "Erro ao criar usuário: "+err.Error(), nil, baseLayout)
	}
	return views.RenderTemplateWithMessage(ctx, "users/create", false, fmt.Sprintf("Usuário %s criado com sucesso!", user.Username), nil, baseLayout)
}

func (c *userController) RenderUsers(ctx *fiber.Ctx) error {
	pagination := models.NewPagination(ctx.QueryInt("page"), ctx.QueryInt("limit"))
	query := ctx.Query("q", "")
	users := c.userService.FindAll(pagination, query)
	data := fiber.Map{
		"Users":      users,
		"Pagination": pagination,
	}
	return views.RenderTemplate(ctx, "users/users", data, baseLayout)
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
	return views.RenderTemplate(ctx, "users/user", user, baseLayout)
}

func (c *userController) HandlerUpdate(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Redirect("/users")
	}
	user, err := c.userService.FindById(uint(id))
	if err != nil {
		return ctx.Redirect("/users")
	}
	if err := ctx.BodyParser(user); err != nil {
		return views.RenderTemplateWithMessage(ctx, "users/user", true, err.Error(), nil, baseLayout)
	}
	oldPassword := ctx.FormValue("oldPassword")
	newPassword := ctx.FormValue("newPassword")
	if oldPassword != "" && newPassword != "" {
		if err := user.UpdatePassword(oldPassword, newPassword); err != nil {
			return views.RenderTemplateWithMessage(ctx, "users/user", true, "Não foi possível atualizar a senha do usuário. Por favor, verifique os dados.", nil, baseLayout)
		}
	}
	user.IsStaff = ctx.FormValue("isStaff") == "on"
	user.IsActive = ctx.FormValue("isActive") == "on"
	if err := c.userService.Update(user); err != nil {
		return views.RenderTemplateWithMessage(ctx, "users/user", true, "Falha ao atualizar usuário.", nil, baseLayout)
	}
	return views.RenderTemplateWithMessage(ctx, "users/user", false, fmt.Sprintf("Usuário %s atualizado com sucesso!", user.Username), user, baseLayout)
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
	return views.RenderTemplate(ctx, "users/delete", user, baseLayout)
}

func (c *userController) HandlerDelete(ctx *fiber.Ctx) error {
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
