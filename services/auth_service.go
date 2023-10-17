package services

import (
	"github.com/bitebait/cupcakestore/repositories"
)

type AuthService interface {
	Authenticate(username, password string) error
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func (s *authService) Authenticate(username, password string) error {
	user, err := s.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}

	return user.CheckPassword(password)
}
