package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Authenticate(ctx *fiber.Ctx, username, password string) error
}

type authService struct {
	userRepository    repositories.UserRepository
	profileRepository repositories.ProfileRepository
}

func NewAuthService(userRepository repositories.UserRepository, profileRepository repositories.ProfileRepository) AuthService {
	return &authService{
		userRepository:    userRepository,
		profileRepository: profileRepository,
	}
}

func (s *authService) Authenticate(ctx *fiber.Ctx, username, password string) error {
	user, err := s.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}

	if err = user.CheckPassword(password); err != nil {
		return err
	}

	profile, err := s.profileRepository.FindByUserId(user.ID)
	if err != nil {
		return err
	}

	return setUserSession(ctx, &profile)
}

func setUserSession(ctx *fiber.Ctx, profile *models.Profile) error {
	sess, err := session.Store.Get(ctx)
	if err != nil {
		return err
	}

	sess.Set("profile", profile)
	if err = sess.Save(); err != nil {
		return err
	}

	return nil
}
