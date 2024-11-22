package controllers

import (
	"errors"
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/messages"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
)

const (
	RegisterTemplate = "auth/register"
	LoginTemplate    = "auth/login"
)

type AuthController interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	RenderLogin(ctx *fiber.Ctx) error
	RenderRegister(ctx *fiber.Ctx) error
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authService: authService}
}

func parseUserFromContext(ctx *fiber.Ctx) (*models.User, error) {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return nil, errors.New("dados inválidos")
	}
	return user, nil
}

func createProfileFromContext(ctx *fiber.Ctx, user *models.User) *models.Profile {
	return &models.Profile{
		FirstName: ctx.FormValue("firstname"),
		LastName:  ctx.FormValue("lastname"),
		User:      *user,
	}
}

func (c *authController) Register(ctx *fiber.Ctx) error {
	user, err := parseUserFromContext(ctx)

	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/auth/register")
	}

	profile := createProfileFromContext(ctx, user)

	if err := c.authService.Register(profile); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/auth/register")
	}

	return ctx.Redirect("/auth/login")
}

func (c *authController) Login(ctx *fiber.Ctx) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	if err := c.authService.Authenticate(ctx, email, password); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect("/auth/login")
	}

	redirectPath := config.Instance().GetEnvVar("REDIRECT_AFTER_LOGIN", "/")

	return ctx.Redirect(redirectPath)
}

func (c *authController) Logout(ctx *fiber.Ctx) error {
	redirectPath := config.Instance().GetEnvVar("REDIRECT_AFTER_LOGOUT", "/")
	sess, err := session.Store.Get(ctx)

	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect(redirectPath)
	}

	if err := sess.Destroy(); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect(redirectPath)

	}

	return ctx.Redirect(redirectPath)

}

func (c *authController) RenderLogin(ctx *fiber.Ctx) error {
	return ctx.Render(LoginTemplate, fiber.Map{})
}

func (c *authController) RenderRegister(ctx *fiber.Ctx) error {
	return ctx.Render(RegisterTemplate, fiber.Map{})
}
