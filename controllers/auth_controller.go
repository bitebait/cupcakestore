package controllers

import (
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/session"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

const (
	RegisterTemplate      = "auth/register"
	LoginTemplate         = "auth/login"
	DefaultLoginRedirect  = "/"
	DefaultLogoutRedirect = "/"
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

func (c *authController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return views.RenderError(ctx, RegisterTemplate, nil, "Dados da conta inválidos: "+err.Error())
	}

	profile := &models.Profile{
		FirstName: ctx.FormValue("firstname"),
		LastName:  ctx.FormValue("lastname"),
		User:      *user,
	}

	if err := c.authService.Register(profile); err != nil {
		return views.RenderError(ctx, RegisterTemplate, nil, "Falha ao criar usuário: "+err.Error())
	}

	return ctx.Redirect("/auth/login")
}

func (c *authController) Login(ctx *fiber.Ctx) error {
	email, password := ctx.FormValue("email"), ctx.FormValue("password")
	if err := c.authService.Authenticate(ctx, email, password); err != nil {
		return views.RenderError(ctx, LoginTemplate, nil, "Credenciais inválidas ou usuário inativo.")
	}
	return ctx.Redirect(config.Instance().GetEnvVar("REDIRECT_AFTER_LOGIN", DefaultLoginRedirect))
}

func (c *authController) Logout(ctx *fiber.Ctx) error {
	sess, err := session.Store.Get(ctx)
	if err != nil {
		return err
	}

	if err := sess.Destroy(); err != nil {
		return err
	}
	return ctx.Redirect(config.Instance().GetEnvVar("REDIRECT_AFTER_LOGOUT", DefaultLogoutRedirect))
}

func (c *authController) RenderLogin(ctx *fiber.Ctx) error {
	return views.Render(ctx, LoginTemplate, nil)
}

func (c *authController) RenderRegister(ctx *fiber.Ctx) error {
	return views.Render(ctx, RegisterTemplate, nil)
}
