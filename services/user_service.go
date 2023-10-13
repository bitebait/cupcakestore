package services

import (
	"log"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type UserService interface {
	Create(user *models.User) error
	List() []*models.User
	FindById(id uint) (*models.User, error)
	Update(user *models.User) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) Create(user *models.User) error {
	return s.userRepository.Create(user)
}

func (s *userService) List() []*models.User {
	users, err := s.userRepository.List()
	if err != nil {
		log.Println(err)
	}
	return users
}

func (s *userService) FindById(id uint) (*models.User, error) {
	return s.userRepository.FindById(id)
}

func (s *userService) Update(user *models.User) error {
	return s.userRepository.Update(user)
}
