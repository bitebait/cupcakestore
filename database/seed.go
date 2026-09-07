package database

import (
	"errors"

	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

type Seeder interface {
	Seed(db *gorm.DB) error
}

type UserAdminSeeder struct {
	Email    string
	Password string
}

func (s UserAdminSeeder) Seed(db *gorm.DB) error {
	// No shared/default administrator credentials are ever created.
	if s.Email == "" && s.Password == "" {
		return nil
	}
	if s.Email == "" || len(s.Password) < 12 || len(s.Password) > 72 {
		return errors.New("administrator seed requires an email and a password of 12 to 72 bytes")
	}
	var admin models.User
	err := db.Unscoped().Where("email = ?", s.Email).First(&admin).Error
	if err == nil {
		// Never overwrite credentials, restore deleted users, or promote an
		// existing customer's account during an application restart.
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return db.Create(&models.User{
		Email: s.Email, Password: s.Password, IsActive: true, IsStaff: true,
	}).Error
}

type StoreConfigSeeder struct{}

func (s StoreConfigSeeder) Seed(db *gorm.DB) error {
	var existing models.StoreConfig
	if err := db.First(&existing).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	storeConfig := &models.StoreConfig{
		DeliveryPrice:       10,
		DeliveryIsActive:    true,
		PaymentCashIsActive: true,
		PixKeyType:          models.PixTypeCPF,
	}
	return db.Create(storeConfig).Error
}

func SeedDatabase(db *gorm.DB) error {
	return seedDatabase(db, config.Get())
}

func seedDatabase(db *gorm.DB, cfg *config.Config) error {
	return db.Transaction(func(tx *gorm.DB) error {
		seeders := []Seeder{
			UserAdminSeeder{Email: cfg.AdminEmail, Password: cfg.AdminPassword},
			StoreConfigSeeder{},
		}
		for _, seeder := range seeders {
			if err := seeder.Seed(tx); err != nil {
				return err
			}
		}
		return nil
	})
}
