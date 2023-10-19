package controllers

import (
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type AuthController interface {
	RenderLogin(ctx *fiber.Ctx) error
	HandlerLogin(ctx *fiber.Ctx) error
	HandlerLogout(ctx *fiber.Ctx) error
}

type authController struct {
	authService      services.AuthService
	templateRenderer views.TemplateRenderer
	store            *session.Store
}

func NewAuthController(authService services.AuthService, templateRenderer views.TemplateRenderer, store *session.Store) AuthController {
	return &authController{
		authService:      authService,
		templateRenderer: templateRenderer,
		store:            store,
	}
}

func (c *authController) RenderLogin(ctx *fiber.Ctx) error {
	return c.templateRenderer.Render(ctx, "auth/login", nil, "")
}

func (c *authController) HandlerLogin(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	err := c.authService.Authenticate(username, password)
	if err != nil {
		return c.templateRenderer.Render(ctx, "auth/login", nil, "Credenciais inválidas")
	}

	sess, err := c.store.Get(ctx)
	if err != nil {
		panic(err)
	}

	sess.Set("username", username)
	if err := sess.Save(); err != nil {
		panic(err)
	}

	return ctx.Redirect("/users")
}

func (c *authController) HandlerLogout(ctx *fiber.Ctx) error {
	sess, err := c.store.Get(ctx)
	if err != nil {
		panic(err)
	}

	sess.Delete("username")
	if err := sess.Destroy(); err != nil {
		panic(err)
	}

	return ctx.Redirect("/auth/login")
}
