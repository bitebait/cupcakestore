package services

import (
	"errors"
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
	if err := s.userService.Create(&profile.User); err != nil {
		return errors.New("falha ao criar o usuário")
	}

	if err := s.updateUserProfile(profile); err != nil {
		return errors.New("falha ao atualizar o perfil do usuário")
	}

	return nil
}

func (s *authService) updateUserProfile(profile *models.Profile) error {
	p, err := s.profileService.FindByUserId(profile.User.ID)

	if err != nil {
		return err
	}

	p.FirstName = profile.FirstName
	p.LastName = profile.LastName

	return s.profileService.Update(&p)
}

func (s *authService) Authenticate(ctx *fiber.Ctx, email, password string) error {
	user, err := s.userService.FindByEmail(email)

	if err != nil || !user.IsActive {
		return errors.New("usuário não encontrado ou inativo")
	}

	if err := user.CheckPassword(password); err != nil {
		return errors.New("senha de acesso incorreta")
	}

	profile, err := s.profileService.FindByUserId(user.ID)

	if err != nil {
		return errors.New("usuário não encontrado ou inativo")
	}

	if err := setupUserSession(ctx, &profile); err != nil {
		return errors.New("falha ao criar sessão de usuário")
	}

	if err := s.registerUserLoginDate(&user); err != nil {
		return errors.New("falha ao registrar a data de login do usuário")
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

	sess.Set("Profile", profile)

	return sess.Save()
}
