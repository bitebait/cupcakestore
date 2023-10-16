package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll(p *models.Pagination, filter string) []*models.User
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

func (r *userRepository) FindAll(p *models.Pagination, filter string) []*models.User {
	offset := (p.Page - 1) * p.Limit

	query := r.db.Model(&models.User{}).Omit("Password")

	if filter != "" {
		filterPattern := "%" + filter + "%"
		query = query.Where("username LIKE ? OR email LIKE ?", filterPattern, filterPattern)
	}

	query.Count(&p.Total)

	var users []*models.User
	query.Offset(offset).Limit(p.Limit).Order("username, email").Find(&users)

	return users
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
		query := r.db.Omit("Password")
		return query.Save(user).Error
	}

	return r.db.Save(user).Error
}

func (r *userRepository) Delete(user *models.User) error {
	return r.db.Delete(user).Error
}
