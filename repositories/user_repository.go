package repositories

import (
	"gorm.io/gorm"
)

type UserRepository interface {
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(database *gorm.DB) UserRepository {
	return &userRepository{
		db: database,
	}
}
