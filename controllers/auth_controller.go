package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/session"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	RenderLogin(ctx *fiber.Ctx) error
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(a services.AuthService) AuthController {
	return &authController{
		authService: a,
	}
}

func (c *authController) RenderLogin(ctx *fiber.Ctx) error {
	return views.Render(ctx, "auth/login", nil, "")
}

func (c *authController) Login(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	err := c.authService.Authenticate(ctx, username, password)
	if err != nil {
		return views.Render(ctx, "auth/login", nil, "Credenciais inválidas")
	}

	return ctx.Redirect("/users")
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
