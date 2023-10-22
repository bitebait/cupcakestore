package services

import (
	"log"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	Authenticate(ctx *fiber.Ctx, username, password string) error
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
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

	return setUserSession(ctx, &user)
}

func setUserSession(ctx *fiber.Ctx, user *models.User) error {
	sess, err := session.Store.Get(ctx)
	if err != nil {
		return err
	}

	sess.Set("user", user)
	if err = sess.Save(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
