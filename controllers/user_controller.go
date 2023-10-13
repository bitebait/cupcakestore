package controllers

import (
	"strconv"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	RenderCreate(c *fiber.Ctx) error
	HandleCreate(c *fiber.Ctx) error
	RenderList(c *fiber.Ctx) error
	RenderDetail(c *fiber.Ctx) error
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{userService}
}

func (c *userController) RenderCreate(ctx *fiber.Ctx) error {
	return ctx.Render("user/create", fiber.Map{}, "layouts/base")
}

func (c *userController) HandleCreate(ctx *fiber.Ctx) error {
	user := new(models.User)

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
	return ctx.Render("user/list", fiber.Map{"Users": c.userService.List()}, "layouts/base")
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
