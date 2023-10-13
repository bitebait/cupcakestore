package services

import (
	"errors"
	"log"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/utils"
)

type UserService interface {
	Create(u *models.User) error
	List() []*models.User
	FindById(id uint) (*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) Create(u *models.User) error {
	if u == nil || u.Password == "" {
		return errors.New("invalid user or empty password")
	}

	hash, err := utils.PasswordHasher(u.Password)
	if err != nil {
		return err
	}

	u.Password = hash

	return s.userRepository.Create(u)
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
