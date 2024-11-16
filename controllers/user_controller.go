package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/utils"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
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
	return views.Render(ctx, tplCriarUsuario, nil, layoutBase)
}

func (c *userController) Create(ctx *fiber.Ctx) error {
	user := &models.User{}
	if err := ctx.BodyParser(user); err != nil {
		return c.renderUserError(ctx, tplCriarUsuario, nil, "Dados de usuário inválidos: "+err.Error(), layoutBase)
	}
	c.extractUserFormData(ctx, user)

	if err := c.userService.Create(user); err != nil {
		return c.renderUserError(ctx, tplCriarUsuario, nil, "Falha ao criar usuário: "+err.Error(), layoutBase)
	}
	return ctx.Redirect("/users")
}

func (c *userController) RenderUsers(ctx *fiber.Ctx) error {
	filter := models.NewUserFilter(ctx.Query("q", ""), ctx.QueryInt("page"), ctx.QueryInt("limit"))
	users := c.userService.FindAll(filter)
	return views.Render(ctx, tplUsuarios, fiber.Map{"Users": users, "Filter": filter}, layoutBase)
}

func (c *userController) RenderUser(ctx *fiber.Ctx) error {
	user, err := c.getUser(ctx)
	if err != nil {
		return err
	}

	userSess, err := c.getUserSession(ctx)
	if err != nil {
		return err
	}
	if user.IsStaff && user.ID != userSess.ID {
		return ctx.Redirect("/users")
	}

	layout := selectLayout(userSess.IsStaff, user.ID == userSess.ID)
	if layout == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}
	return views.Render(ctx, tplUsuario, user, layout)
}

func (c *userController) getUser(ctx *fiber.Ctx) (models.User, error) {
	userID, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return models.User{}, ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return c.userService.FindById(userID)
}

func (c *userController) getUserSession(ctx *fiber.Ctx) (*models.User, error) {
	userSess, ok := ctx.Locals("Profile").(*models.Profile)
	if !ok || userSess == nil {
		return nil, fiber.ErrUnauthorized
	}
	return &userSess.User, nil
}

func (c *userController) Update(ctx *fiber.Ctx) error {
	user, err := c.getUserAndCheckAccess(ctx)
	if err != nil {
		return err
	}
	if err := c.updateUserFromRequest(ctx, user); err != nil {
		return c.renderUserError(ctx, tplUsuario, user, err.Error(), selectLayout(user.IsStaff, user.ID == ctx.Locals("Profile").(*models.Profile).UserID))
	}
	if err := c.updateUserPassword(ctx, user); err != nil {
		return c.renderUserError(ctx, tplUsuario, user, "Falha ao atualizar a senha. Certifique-se de que está inserido corretamente.", selectLayout(user.IsStaff, user.ID == ctx.Locals("Profile").(*models.Profile).UserID))
	}
	if err := c.userService.Update(user); err != nil {
		return c.renderUserError(ctx, tplUsuario, user, "Falha ao atualizar o usuário.", selectLayout(user.IsStaff, user.ID == ctx.Locals("Profile").(*models.Profile).UserID))
	}
	if user.ID == ctx.Locals("Profile").(*models.Profile).UserID {
		return ctx.Redirect("/auth/logout")
	}
	return ctx.Redirect(selectRedirectPath(user.IsStaff))
}

func (c *userController) getUserAndCheckAccess(ctx *fiber.Ctx) (*models.User, error) {
	id, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return nil, ctx.Redirect("/users")
	}
	user, err := c.userService.FindById(id)
	if err != nil {
		return nil, ctx.Redirect("/users")
	}
	userSess, err := c.getUserSession(ctx)
	if err != nil {
		return nil, err
	}
	if user.IsStaff && user.ID != userSess.ID {
		return nil, ctx.SendStatus(fiber.StatusUnauthorized)
	}
	if !userSess.IsStaff && user.ID != userSess.ID {
		return nil, ctx.SendStatus(fiber.StatusUnauthorized)
	}
	return &user, nil
}

func (c *userController) updateUserFromRequest(ctx *fiber.Ctx, user *models.User) error {
	if err := ctx.BodyParser(user); err != nil {
		return err
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
		return ctx.Redirect("/users")
	}
	return views.Render(ctx, tplExcluirUsuario, user, layoutBase)
}

func (c *userController) Delete(ctx *fiber.Ctx) error {
	id, err := utils.StringToId(ctx.Params("id"))
	if err != nil {
		return ctx.Redirect("/users")
	}
	if err := c.userService.Delete(id); err != nil {
		return ctx.Redirect("/users")
	}
	return ctx.Redirect("/users")
}

func (c *userController) renderUserError(ctx *fiber.Ctx, template string, user interface{}, message string, layout string) error {
	return views.RenderError(ctx, template, user, message, layout)
}

func (c *userController) extractUserFormData(ctx *fiber.Ctx, user *models.User) {
	user.IsStaff = ctx.FormValue("isStaff") == "on"
	user.IsActive = ctx.FormValue("isActive") == "on"
}
