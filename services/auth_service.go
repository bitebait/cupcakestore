package services

import (
	"errors"
	"strings"
	"time"

	"github.com/bitebait/cupcakestore/models"
)

type AuthService interface {
	Register(profile *models.Profile) error
	Authenticate(email, password string) (models.Profile, error)
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

func (s *authService) Authenticate(email, password string) (models.Profile, error) {
	user, err := s.userService.FindByEmail(strings.ToLower(strings.TrimSpace(email)))

	if err != nil || !user.IsActive {
		return models.Profile{}, errors.New("e-mail ou senha inválidos")
	}

	if err := user.CheckPassword(password); err != nil {
		return models.Profile{}, errors.New("e-mail ou senha inválidos")
	}

	profile, err := s.profileService.FindByUserId(user.ID)

	if err != nil {
		return models.Profile{}, errors.New("e-mail ou senha inválidos")
	}

	if err := s.userService.RecordLogin(user.ID, time.Now()); err != nil {
		return models.Profile{}, errors.New("falha ao registrar a data de login do usuário")
	}

	profile.User.Password = ""
	return profile, nil
}
