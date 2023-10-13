package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *models.User) error
	List() ([]*models.User, error)
	FindById(id uint) (*models.User, error)
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
	err := r.db.Create(u).Error
	return err
}

func (r *userRepository) List() ([]*models.User, error) {
	var users []*models.User

	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindById(id uint) (*models.User, error) {
	var user models.User

	err := r.db.First(&user, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &user, err
}
