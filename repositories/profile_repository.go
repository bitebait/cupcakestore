package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	Create(profile *models.Profile) error
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
