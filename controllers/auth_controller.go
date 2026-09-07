package controllers

import (
	"errors"
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/messages"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v3"
)

const (
	RegisterTemplate = "auth/register"
	LoginTemplate    = "auth/login"
)

type AuthController interface {
	Register(ctx fiber.Ctx) error
	Login(ctx fiber.Ctx) error
	Logout(ctx fiber.Ctx) error
	RenderLogin(ctx fiber.Ctx) error
	RenderRegister(ctx fiber.Ctx) error
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authService: authService}
}

func parseUserFromContext(ctx fiber.Ctx) (*models.User, error) {
	var input struct {
		Email    string `form:"email" json:"email" validate:"required,email"`
		Password string `form:"password" json:"password" validate:"required,min=8,max=72"`
	}
	if err := ctx.Bind().Body(&input); err != nil {
		return nil, errors.New("dados inválidos")
	}
	user := &models.User{Email: input.Email, Password: input.Password, IsActive: true}
	return user, nil
}

func createProfileFromContext(ctx fiber.Ctx, user *models.User) *models.Profile {
	return &models.Profile{
		FirstName: ctx.FormValue("firstname"),
		LastName:  ctx.FormValue("lastname"),
		User:      *user,
	}
}

func (c *authController) Register(ctx fiber.Ctx) error {
	user, err := parseUserFromContext(ctx)

	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect().Status(fiber.StatusFound).To("/auth/register")
	}

	profile := createProfileFromContext(ctx, user)

	if err := c.authService.Register(profile); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect().Status(fiber.StatusFound).To("/auth/register")
	}

	return ctx.Redirect().Status(fiber.StatusFound).To("/auth/login")
}

func (c *authController) Login(ctx fiber.Ctx) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	profile, err := c.authService.Authenticate(email, password)
	if err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect().Status(fiber.StatusFound).To("/auth/login")
	}
	if err := session.Login(ctx, profile.UserID); err != nil {
		messages.SetErrorMessage(ctx, "não foi possível iniciar a sessão")
		return ctx.Redirect().Status(fiber.StatusFound).To("/auth/login")
	}

	return ctx.Redirect().Status(fiber.StatusFound).To(config.Get().RedirectAfterLogin)
}

func (c *authController) Logout(ctx fiber.Ctx) error {
	redirectPath := config.Get().RedirectAfterLogout
	sess := session.FromContext(ctx)
	if sess == nil {
		return fiber.ErrInternalServerError
	}
	if err := sess.Reset(); err != nil {
		messages.SetErrorMessage(ctx, err.Error())
		return ctx.Redirect().Status(fiber.StatusFound).To(redirectPath)

	}

	return ctx.Redirect().Status(fiber.StatusFound).To(redirectPath)

}

func (c *authController) RenderLogin(ctx fiber.Ctx) error {
	return ctx.Render(LoginTemplate, fiber.Map{})
}

func (c *authController) RenderRegister(ctx fiber.Ctx) error {
	return ctx.Render(RegisterTemplate, fiber.Map{})
}
