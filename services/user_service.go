package services

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type UserService interface {
	Create(user *models.User) error
	FindAll(p *models.Pagination) []*models.User
	FindById(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
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

func (s *userService) FindAll(p *models.Pagination) []*models.User {
	return s.userRepository.FindAll(p)
}

func (s *userService) FindById(id uint) (*models.User, error) {
	return s.userRepository.FindById(id)
}

func (s *userService) Update(user *models.User) error {
	return s.userRepository.Update(user)
}

func (s *userService) Delete(id uint) error {
	user, err := s.FindById(id)
	if err != nil {
		return err
	}
	return s.userRepository.Delete(user)
}
