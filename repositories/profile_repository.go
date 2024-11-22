package repositories

import (
	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
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
		log.Errorf("ProfileRepository Create: %s", err.Error())
		return err
	}

	return nil
}

func (r *profileRepository) FindByUserId(userID uint) (models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&profile).Error

	if err != nil {
		log.Errorf("ProfileRepository FindOrCreateByUserId: %s", err.Error())
	}

	return profile, err
}

func (r *profileRepository) Update(profile *models.Profile) error {
	if err := r.db.Save(profile).Error; err != nil {
		log.Errorf("ProfileRepository Update: %s", err.Error())
		return err
	}

	return nil
}
