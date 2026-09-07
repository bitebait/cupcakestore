package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
	"log/slog"
)

type ProfileRepository interface {
	Create(profile *models.Profile) error
	FindByUserId(userID uint) (models.Profile, error)
	Update(profile *models.Profile) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) Create(profile *models.Profile) error {
	if err := r.db.Create(profile).Error; err != nil {
		slog.Error("ProfileRepository Create", "error", err)
		return err
	}

	return nil
}

func (r *profileRepository) FindByUserId(userID uint) (models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&profile).Error

	if err != nil {
		slog.Error("ProfileRepository FindOrCreateByUserId", "error", err)
	}

	return profile, err
}

func (r *profileRepository) Update(profile *models.Profile) error {
	if profile == nil || profile.ID == 0 {
		return gorm.ErrRecordNotFound
	}
	result := r.db.Model(profile).
		Select("FirstName", "LastName", "Address", "City", "State", "PostalCode", "PhoneNumber").Updates(profile)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
