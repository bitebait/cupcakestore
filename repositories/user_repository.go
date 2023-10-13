package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	List() ([]*models.User, error)
	FindById(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(database *gorm.DB) UserRepository {
	return &userRepository{
		db: database,
	}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) List() ([]*models.User, error) {
	var users []*models.User

	err := r.db.Omit("Password").Find(&users).Error
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

func (r *userRepository) Update(user *models.User) error {
	if user.Password == "" {
		return r.db.Omit("Password").Save(user).Error
	}

	return r.db.Save(user).Error
}

func (r *userRepository) Delete(user *models.User) error {
	return r.db.Delete(user).Error
}
