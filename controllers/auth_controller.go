package controllers

import (
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/session"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
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
	return &authController{
		authService: authService,
	}
}

func (c *authController) Register(ctx *fiber.Ctx) error {
	user := &models.User{}
	if err := ctx.BodyParser(user); err != nil {
		return views.Render(ctx, "auth/register", nil,
			"Dados da conta inválidos: "+err.Error(), "")
	}

	profile := &models.Profile{
		FirstName: ctx.FormValue("firstname"),
		LastName:  ctx.FormValue("lastname"),
		User:      *user,
	}

	if err := c.authService.Register(profile); err != nil {
		return views.Render(ctx, "auth/register", nil,
			"Falha ao criar usuário: "+err.Error(), "")
	}

	return ctx.Redirect("/auth/login")
}

func (c *authController) Login(ctx *fiber.Ctx) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	err := c.authService.Authenticate(ctx, email, password)
	if err != nil {
		return views.Render(ctx, "auth/login", nil, "Credenciais inválidas")
	}

	return ctx.Redirect(config.GetEnv("REDIRECT_AFTER_LOGIN", "/"))
}

func (c *authController) Logout(ctx *fiber.Ctx) error {
	sess, err := session.Store.Get(ctx)
	if err != nil {
		panic(err)
	}

	sess.Delete("profile")
	if err := sess.Destroy(); err != nil {
		panic(err)
	}

	return ctx.Redirect("/auth/login")
}

func (c *authController) RenderLogin(ctx *fiber.Ctx) error {
	return views.Render(ctx, "auth/login", nil, "")
}

func (c *authController) RenderRegister(ctx *fiber.Ctx) error {
	return views.Render(ctx, "auth/register", nil, "")
}
