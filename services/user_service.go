package services

import (
	"github.com/bitebait/cupcakestore/repositories"
)

type UserService interface {
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
