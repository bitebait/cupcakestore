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
	err := s.userService.Create(&profile.User)
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}

	if !user.IsActive {
		return errors.New("usuário inativo")
	}

	if err = user.CheckPassword(password); err != nil {
		return err
	}

	profile, err := s.profileService.FindByUserId(user.ID)
	if err != nil {
		return err
	}

	err = setUserSession(ctx, &profile)
	if err != nil {
		return err
	}

	return s.registerUserLoginDate(&user)
}

func (s *authService) registerUserLoginDate(user *models.User) error {
	if user.FirstLogin.IsZero() {
		user.FirstLogin = time.Now()
	}

	user.LastLogin = time.Now()
	return s.userService.Update(user)
}

func setUserSession(ctx *fiber.Ctx, profile *models.Profile) error {
	sess, err := session.Store.Get(ctx)
	if err != nil {
		return err
	}

	sess.Set("Profile", profile)
	if err = sess.Save(); err != nil {
		return err
	}

	return nil
}
