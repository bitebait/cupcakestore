package services

import (
	"errors"
	"strings"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
)

type UserService interface {
	Create(user *models.User) error
	FindAll(filter *models.UserFilter) []models.User
	FindById(id uint) (models.User, error)
	FindByEmail(email string) (models.User, error)
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
	s.normalizeUser(user)

	if err := s.userRepository.Create(user); err != nil {
		return errors.New("falha ao cadastrar o usuário, verifique os dados ou tente um e-mail diferente")
	}

	return nil
}

func (s *userService) FindAll(filter *models.UserFilter) []models.User {
	return s.userRepository.FindAll(filter)
}

func (s *userService) FindById(id uint) (models.User, error) {
	user, err := s.userRepository.FindById(id)

	if err != nil {
		err = errors.New("falha ao encontrar o usuário com o id informado")
	}

	return user, err
}

func (s *userService) FindByEmail(email string) (models.User, error) {
	user, err := s.userRepository.FindByEmail(email)

	if err != nil {
		err = errors.New("falha ao encontrar o usuário com o email informado")
	}

	return user, err
}

func (s *userService) Update(user *models.User) error {
	s.normalizeUser(user)

	if err := s.userRepository.Update(user); err != nil {
		return errors.New("falha ao atualizar o usuário")
	}

	return nil
}

func (s *userService) Delete(id uint) error {
	user, err := s.FindById(id)

	if err != nil {
		return errors.New("falha ao encontrar o usuário com o id informado")
	}

	if err := s.userRepository.Delete(&user); err != nil {
		return errors.New("falha ao deletar o usuário")
	}

	return nil
}

func (s *userService) normalizeUser(user *models.User) {
	user.Email = strings.ToLower(user.Email)
}
