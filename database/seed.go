package database

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
	"log"
)

type Seeder interface {
	Seed(db *gorm.DB) error
}

type UserAdminSeeder struct{}

func (s UserAdminSeeder) Seed(db *gorm.DB) error {
	admin := &models.User{
		Email:    "admin@admin.com",
		Password: "admin@admin.com",
		IsActive: true,
		IsStaff:  true,
	}

	if err := db.FirstOrCreate(&admin, "email = ?", admin.Email).Error; err != nil {
		log.Fatalf("Failed to create AdminUser: %v", err)
		return err
	}

	return nil
}

type StoreConfigSeeder struct{}

func (s StoreConfigSeeder) Seed(db *gorm.DB) error {
	storeConfig := &models.StoreConfig{
		DeliveryPrice:            10,
		DeliveryIsActive:         true,
		PhysicalStoreEmail:       "foo@bar.com",
		PhysicalStoreAddress:     "Foo Bar",
		PhysicalStoreCity:        "Foo Bar",
		PhysicalStoreState:       "Foo Bar",
		PhysicalStorePostalCode:  "00000-000",
		PhysicalStorePhoneNumber: "(00)00000-0000",
		PaymentCashIsActive:      true,
		PaymentPixIsActive:       true,
		PixKey:                   "000.000.000-00",
		PixKeyType:               models.PixTypeCPF,
	}

	if err := db.FirstOrCreate(&storeConfig, "physical_store_email = ?", storeConfig.PhysicalStoreEmail).Error; err != nil {
		log.Fatalf("Failed to create StoreConfig: %v", err)
		return err
	}

	return nil
}

func SeedDatabase(db *gorm.DB) error {
	seeders := []Seeder{
		UserAdminSeeder{},
		StoreConfigSeeder{},
	}

	for _, seeder := range seeders {
		if err := seeder.Seed(db); err != nil {
			return err
		}
	}
	return nil
}
