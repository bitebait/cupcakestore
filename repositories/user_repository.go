package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *models.User) error
	List() ([]*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(database *gorm.DB) UserRepository {
	return &userRepository{
		db: database,
	}
}

func (r *userRepository) Create(u *models.User) error {
	res := r.db.Create(u)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *userRepository) List() ([]*models.User, error) {
	var users []*models.User

	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
