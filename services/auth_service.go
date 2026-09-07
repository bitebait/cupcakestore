package services

import (
	"errors"
	"strings"
	"time"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Register(profile *models.Profile) error
	Authenticate(ctx *fiber.Ctx, email, password string) error
}

type authService struct {
	userService    UserService
	profileService ProfileService
}

func NewAuthService(userService UserService, profileService ProfileService) AuthService {
	return &authService{
		userService:    userService,
		profileService: profileService,
	}
}

func (s *authService) Register(profile *models.Profile) error {
	profile.User.IsStaff = false
	profile.User.IsActive = true
	return s.userService.Register(profile)
}

func (s *authService) Authenticate(ctx *fiber.Ctx, email, password string) error {
	user, err := s.userService.FindByEmail(strings.ToLower(strings.TrimSpace(email)))

	if err != nil || !user.IsActive {
		return errors.New("e-mail ou senha inválidos")
	}

	if err := user.CheckPassword(password); err != nil {
		return errors.New("e-mail ou senha inválidos")
	}

	profile, err := s.profileService.FindByUserId(user.ID)

	if err != nil {
		return errors.New("e-mail ou senha inválidos")
	}

	if err := s.registerUserLoginDate(&user); err != nil {
		return errors.New("falha ao registrar a data de login do usuário")
	}
	if err := setupUserSession(ctx, &profile); err != nil {
		return errors.New("falha ao criar sessão de usuário")
	}

	return nil
}

func (s *authService) registerUserLoginDate(user *models.User) error {
	now := time.Now()

	if user.FirstLogin.IsZero() {
		user.FirstLogin = now
	}

	user.LastLogin = now

	return s.userService.Update(user)
}

func setupUserSession(ctx *fiber.Ctx, profile *models.Profile) error {
	sess, err := session.Store.Get(ctx)

	if err != nil {
		return err
	}

	if err := sess.Regenerate(); err != nil {
		return err
	}
	// Never serialize a password hash into the session.
	identity := &models.Profile{UserID: profile.UserID}
	sess.Set("Profile", identity)

	return sess.Save()
}
