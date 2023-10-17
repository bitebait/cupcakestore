package controllers

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AuthController interface {
	RenderLogin(ctx *fiber.Ctx) error
	HandlerLogin(ctx *fiber.Ctx) error
}

type authController struct {
	authService services.AuthService
	store       *session.Store
}

func NewAuthController(authService services.AuthService, store *session.Store) AuthController {
	return &authController{authService: authService, store: store}
}

func (c *authController) RenderLogin(ctx *fiber.Ctx) error {
	return ctx.Render("auth/login", models.NewResponse(false, nil, ""))
}

func (c *authController) HandlerLogin(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	// Autenticar usuário
	err := c.authService.Authenticate(username, password)
	if err != nil {
		return ctx.Render("auth/login", models.NewResponse(true, nil, "Credenciais inválidas"))
	}

	sess, err := c.store.Get(ctx)
	if err != nil {
		panic(err)
	}

	if sess.Fresh() {
		sess.Set("username", username)
		if err := sess.Save(); err != nil {
			panic(err)
		}

	}

	return ctx.Redirect("/users/list")
}
