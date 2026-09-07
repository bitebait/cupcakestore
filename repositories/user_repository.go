package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type UserRepository interface {
	Create(user *models.User) error
	CreateWithProfile(profile *models.Profile) error
	FindAll(filter *models.UserFilter) []models.User
	FindById(id uint) (models.User, error)
	FindByEmail(email string) (models.User, error)
	Update(user *models.User) error
	RecordLogin(id uint, at time.Time) error
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
		slog.Error("UserRepository Create", "error", err)
		return err
	}

	return nil
}

func (r *userRepository) FindAll(filter *models.UserFilter) []models.User {
	var total int64
	query := r.buildFilteredQuery(filter)

	if err := query.Count(&total).Error; err != nil {
		slog.Error("UserRepository FindAll", "error", err)
	}

	var users []models.User
	filter.Pagination.Total = total
	offset := (filter.Pagination.Page - 1) * filter.Pagination.Limit

	if err := query.Offset(offset).Limit(filter.Pagination.Limit).Order("created_at desc").Find(&users).Error; err != nil {
		slog.Error("UserRepository FindAll", "error", err)
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
		slog.Error("UserRepository FindOrCreateById", "error", err)
	}

	return user, err
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		slog.Error("UserRepository FindByEmail", "error", err)
	}

	return user, err
}

func (r *userRepository) Update(user *models.User) error {
	if user == nil || user.ID == 0 {
		return gorm.ErrRecordNotFound
	}
	fields := []string{"Email", "IsActive", "IsStaff"}
	if user.Password != "" {
		fields = append(fields, "Password")
	}
	result := r.db.Model(user).Where("updated_at = ?", user.UpdatedAt).Select(fields).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *userRepository) Delete(user *models.User) error {
	if err := r.db.Select("Profile").Delete(user).Error; err != nil {
		slog.Error("UserRepository Delete", "error", err)
		return err
	}

	return nil
}

// CreateWithProfile commits the account and registration details together.
func (r *userRepository) CreateWithProfile(profile *models.Profile) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&profile.User).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Profile{}).Where("user_id = ?", profile.User.ID).
			Updates(map[string]interface{}{"first_name": profile.FirstName, "last_name": profile.LastName}).Error; err != nil {
			return err
		}
		return tx.Preload("User").Where("user_id = ?", profile.User.ID).First(profile).Error
	})
}

// RecordLogin never writes credential or permission fields from a stale snapshot.
func (r *userRepository) RecordLogin(id uint, at time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&models.User{}).Where("id = ? AND is_active = ?", id, true).
			UpdateColumn("last_login", at)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		return tx.Model(&models.User{}).Where("id = ? AND (first_login IS NULL OR first_login = ?)", id, time.Time{}).
			UpdateColumn("first_login", at).Error
	})
}
