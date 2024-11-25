package controllers

import (
	"errors"
	"github.com/bitebait/cupcakestore/helpers"
	"github.com/bitebait/cupcakestore/messages"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

const (
	tplCriarUsuario   = "users/create"
	tplUsuarios       = "users/users"
	tplUsuario        = "users/user"
	tplExcluirUsuario = "users/delete"
	layoutBase        = views.BaseLayout
)

type UserController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	RenderCreate(ctx *fiber.Ctx) error
	RenderUsers(ctx *fiber.Ctx) error
	RenderUser(ctx *fiber.Ctx) error
	RenderDelete(ctx *fiber.Ctx) error
}

type userController struct {
	userService services.UserService
}

func NewUserController(u services.UserService) UserController {
	return &userController{
		userService: u,
	}
}

func (c *userController) RenderCreate(ctx *fiber.Ctx) error {
	return ctx.Render(tplCriarUsuario, fiber.Map{}, layoutBase)
}

func (c *userController) Create(ctx *fiber.Ctx) error {
	var user models.User

	if err := ctx.BodyParser(&user); err != nil {
		messages.SetErrorMessage(ctx, "erro ao processar os dados do usuário")
		return ctx.Redirect("/users/create")
	}

	c.extractUserFormData(ctx, &user)

	if err := c.userService.Create(&user); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/users/create")
	}

	messages.SetSuccessMessage(ctx, "usuário criado com sucesso")
	return ctx.Redirect("/users")
}

func (c *userController) RenderUsers(ctx *fiber.Ctx) error {
	filter := models.NewUserFilter(ctx.Query("q", ""), ctx.QueryInt("page"), ctx.QueryInt("limit"))
	users := c.userService.FindAll(filter)

	data := fiber.Map{"Users": users, "Filter": filter}

	return ctx.Render(tplUsuarios, fiber.Map{"Object": data}, layoutBase)
}

func (c *userController) RenderUser(ctx *fiber.Ctx) error {
	user, err := c.getUser(ctx)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/users")
	}

	userSess, err := c.getUserSession(ctx)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/users")
	}

	if user.IsStaff && user.ID != userSess.ID {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/users")
	}

	layout := selectLayout(userSess.IsStaff, user.ID == userSess.ID)
	if layout == "" {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/users")
	}

	return ctx.Render(tplUsuario, fiber.Map{"Object": user}, layout)
}

func (c *userController) getUser(ctx *fiber.Ctx) (models.User, error) {
	userID, err := helpers.ParseStringToID(ctx.Params("id"))

	if err != nil {
		return models.User{}, err
	}

	return c.userService.FindById(userID)
}

func (c *userController) getUserSession(ctx *fiber.Ctx) (*models.User, error) {
	userSess, ok := ctx.Locals("Profile").(*models.Profile)

	if !ok || userSess == nil {
		return nil, errors.New("acesso negado")
	}

	return &userSess.User, nil
}

func (c *userController) Update(ctx *fiber.Ctx) error {
	user, err := c.getUserAndCheckAccess(ctx)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/users")
	}

	isProfileUser := user.ID == ctx.Locals("Profile").(*models.Profile).UserID
	isStaff := user.IsStaff
	layout := selectLayout(isStaff, isProfileUser)

	if err := c.updateUserFromRequest(ctx, user); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Render(tplUsuario, fiber.Map{"Object": user}, layout)
	}

	if err := c.updateUserPassword(ctx, user); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Render(tplUsuario, fiber.Map{"Object": user}, layout)
	}

	if err := c.userService.Update(user); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Render(tplUsuario, fiber.Map{"Object": user}, layout)
	}

	if isProfileUser {
		return ctx.Redirect("/auth/logout")
	}

	messages.SetSuccessMessage(ctx, "usuário atualizado com sucesso")
	return ctx.Redirect("/users/user/" + strconv.Itoa(int(user.ID)))
}

func (c *userController) getUserAndCheckAccess(ctx *fiber.Ctx) (*models.User, error) {
	id, err := helpers.ParseStringToID(ctx.Params("id"))
	if err != nil {
		return nil, err
	}

	user, err := c.userService.FindById(id)
	if err != nil {
		return nil, err
	}

	userSess, err := c.getUserSession(ctx)
	if err != nil {
		return nil, errors.New("acesso negado")
	}

	if user.IsStaff && user.ID != userSess.ID {
		return nil, errors.New("acesso não autorizado para staff")
	}

	if !userSess.IsStaff && user.ID != userSess.ID {
		return nil, errors.New("acesso negado")
	}

	return &user, nil
}

func (c *userController) updateUserFromRequest(ctx *fiber.Ctx, user *models.User) error {
	if err := ctx.BodyParser(user); err != nil {
		return errors.New("erro ao processar os dados da requisição")
	}

	c.extractUserFormData(ctx, user)

	return nil
}

func (c *userController) updateUserPassword(ctx *fiber.Ctx, user *models.User) error {
	oldPassword := ctx.FormValue("oldPassword")
	newPassword := ctx.FormValue("newPassword")

	if oldPassword != "" && newPassword != "" {
		if err := user.UpdatePassword(oldPassword, newPassword); err != nil {
			return err
		}
	}

	return nil
}

func (c *userController) RenderDelete(ctx *fiber.Ctx) error {
	user, err := c.getUser(ctx)

	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/users")
	}

	return ctx.Render(tplExcluirUsuario, fiber.Map{"Object": user})
}

func (c *userController) Delete(ctx *fiber.Ctx) error {
	id, err := helpers.ParseStringToID(ctx.Params("id"))

	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/users")
	}

	if err := c.userService.Delete(id); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/users")
	}

	messages.SetSuccessMessage(ctx, "usuário deletado com sucesso")
	return ctx.Redirect("/users")
}

func (c *userController) extractUserFormData(ctx *fiber.Ctx, user *models.User) {
	user.IsStaff = ctx.FormValue("isStaff") == "on"
	user.IsActive = ctx.FormValue("isActive") == "on"
}
