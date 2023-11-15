package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/session"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
	"log"
)

type AuthController interface {
	RenderLogin(ctx *fiber.Ctx) error
	HandlerLogin(ctx *fiber.Ctx) error
	HandlerLogout(ctx *fiber.Ctx) error
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

func (c *authController) HandlerLogin(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	err := c.authService.Authenticate(ctx, username, password)
	log.Println(err)
	if err != nil {
		return views.Render(ctx, "auth/login", nil, "Credenciais inválidas")
	}

	return ctx.Redirect("/users")
}

func (c *authController) HandlerLogout(ctx *fiber.Ctx) error {
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
