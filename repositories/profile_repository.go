package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	Create(profile *models.Profile) error
	FindByUserId(id uint) (models.Profile, error)
	Update(profile *models.Profile) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(database *gorm.DB) ProfileRepository {
	return &profileRepository{
		db: database,
	}
}

func (r *profileRepository) Create(profile *models.Profile) error {
	return r.db.Create(profile).Error
}

func (r *profileRepository) FindByUserId(id uint) (models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("user_id = ?", id).Preload("User").First(&profile).Error
	return profile, err
}

func (r *profileRepository) Update(profile *models.Profile) error {
	return r.db.Save(profile).Error
}
