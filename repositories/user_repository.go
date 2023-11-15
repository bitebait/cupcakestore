package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll(filter *models.UserFilter) []models.User
	FindById(id uint) (models.User, error)
	FindByUsername(username string) (models.User, error)
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

func (r *userRepository) FindAll(filter *models.UserFilter) []models.User {
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit

	query := r.db.Model(&models.User{}).Omit("Password")

	if filter.User.Username != "" {
		filterPattern := "%" + filter.User.Username + "%"
		query = query.Where("username LIKE ? OR email LIKE ?", filterPattern, filterPattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil
	}
	filter.Pagination.Total = total

	var users []models.User
	query.Offset(offset).Limit(filter.Pagination.Limit).Order("created_at desc").Find(&users)

	return users
}

func (r *userRepository) FindById(id uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *userRepository) FindByUsername(username string) (models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *userRepository) Update(user *models.User) error {
	if user.Password == "" {
		query := r.db.Omit("Password")
		return query.Save(user).Error
	}

	return r.db.Save(user).Error
}

func (r *userRepository) Delete(user *models.User) error {
	return r.db.Select("Profile").Delete(user).Error
}
