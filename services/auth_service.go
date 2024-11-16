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

func NewAuthService(u UserService, p ProfileService) AuthService {
	return &authService{
		userService:    u,
		profileService: p,
	}
}

func (s *authService) Register(profile *models.Profile) error {
	if err := s.userService.Create(&profile.User); err != nil {
		return err
	}
	return s.updateUserProfile(profile)
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
	const inactiveUserError = "usuário inativo"

	user, err := s.userService.FindByEmail(email)
	if err != nil {
		return err
	}

	if !user.IsActive {
		return errors.New(inactiveUserError)
	}

	if err = user.CheckPassword(password); err != nil {
		return err
	}

	profile, err := s.profileService.FindByUserId(user.ID)
	if err != nil {
		return err
	}

	if err = setUserSession(ctx, &profile); err != nil {
		return err
	}

	return s.registerUserLoginDate(&user)
}

func (s *authService) registerUserLoginDate(user *models.User) error {
	now := time.Now()
	if user.FirstLogin.IsZero() {
		user.FirstLogin = now
	}
	user.LastLogin = now
	return s.userService.Update(user)
}

func setUserSession(ctx *fiber.Ctx, profile *models.Profile) error {
	sess, err := session.Store.Get(ctx)
	if err != nil {
		return err
	}
	sess.Set("Profile", profile)
	return sess.Save()
}
