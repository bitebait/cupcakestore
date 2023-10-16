package controllers

import (
	"fmt"
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
		return ctx.Render("users/create", models.NewResponse(true, nil, "Dados do usuário inválidos: "+err.Error()), baseLayout)
	}

	if err := c.userService.Create(user); err != nil {
		return ctx.Render("users/create", models.NewResponse(true, nil, "Erro ao criar usuário: "+err.Error()), baseLayout)
	}

	return ctx.Render("users/create", models.NewResponse(false, nil, fmt.Sprintf("Usuário %s criado com sucesso!", user.Username)), baseLayout)
}

func (c *userController) RenderUsers(ctx *fiber.Ctx) error {
	pagination := models.NewPagination(ctx.QueryInt("page"), ctx.QueryInt("limit"))
	query := ctx.Query("q", "")

	users := c.userService.FindAll(pagination, query)

	data := fiber.Map{
		"Users":      users,
		"Pagination": pagination,
	}
	return ctx.Render("users/users", models.NewResponse(false, data, ""), baseLayout)
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

	return ctx.Render("users/user", models.NewResponse(false, user, ""), baseLayout)
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
		return ctx.Render("users/user", models.NewResponse(true, user, err.Error()), baseLayout)
	}

	oldPassword := ctx.FormValue("oldPassword")
	newPassword := ctx.FormValue("newPassword")
	if oldPassword != "" && newPassword != "" {
		if err := user.UpdatePassword(oldPassword, newPassword); err != nil {
			return ctx.Render("users/user",
				models.NewResponse(true, user,
					"Não foi possível atualizar a senha do usuário. Por favor, verifique os dados."),
				baseLayout)
		}
	}

	user.IsStaff = ctx.FormValue("isStaff") == "on"
	user.IsActive = ctx.FormValue("isActive") == "on"

	if err := c.userService.Update(user); err != nil {
		return ctx.Render("users/user", models.NewResponse(true, user, "Falha ao atualizar usuário."), baseLayout)
	}

	return ctx.Render("users/user", models.NewResponse(false, user, fmt.Sprintf("Usuário %s atualizado com sucesso!", user.Username)), baseLayout)

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

	return ctx.Render("users/delete", models.NewResponse(false, user, ""), baseLayout)
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
