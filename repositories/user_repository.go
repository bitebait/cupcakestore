package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll(filter *models.UserFilter) []models.User
	FindById(id uint) (models.User, error)
	FindByEmail(email string) (models.User, error)
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
	if err := r.db.Create(user).Error; err != nil {
		log.Errorf("UserRepository Create: %s", err.Error())
		return err
	}

	return nil
}

func (r *userRepository) FindAll(filter *models.UserFilter) []models.User {
	var total int64
	query := r.buildFilteredQuery(filter)

	if err := query.Count(&total).Error; err != nil {
		log.Errorf("UserRepository FindAll: %s", err.Error())
	}

	var users []models.User
	filter.Pagination.Total = total
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit

	if err := query.Offset(offset).Limit(filter.Pagination.Limit).Order("created_at desc").Find(&users).Error; err != nil {
		log.Errorf("UserRepository FindAll: %s", err.Error())
	}

	return users
}

func (r *userRepository) buildFilteredQuery(filter *models.UserFilter) *gorm.DB {
	query := r.db.Model(&models.User{}).Omit("Password")

	if filter.User.Email != "" {
		query = query.Where("email LIKE ?", "%"+filter.User.Email+"%")
	}

	return query
}

func (r *userRepository) FindById(id uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error

	if err != nil {
		log.Errorf("UserRepository FindOrCreateById: %s", err.Error())
	}

	return user, err
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		log.Errorf("UserRepository FindByEmail: %s", err.Error())
	}

	return user, err
}

func (r *userRepository) Update(user *models.User) error {
	query := r.db

	if user.Password == "" {
		query = query.Omit("Password")
	}

	if err := query.Save(user).Error; err != nil {
		log.Errorf("UserRepository Update: %s", err.Error())
		return err
	}

	return nil
}

func (r *userRepository) Delete(user *models.User) error {
	if err := r.db.Select("Profile").Delete(user).Error; err != nil {
		log.Errorf("UserRepository Delete: %s", err.Error())
		return err
	}

	return nil
}
