package database

import (
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
	"log"
)

const (
	adminEmail               = "admin@admin.com"
	adminPassword            = "admin@admin.com"
	storePhysicalEmail       = "foo@bar.com"
	storePhysicalAddress     = "Foo Bar"
	storePhysicalCity        = "Foo Bar"
	storePhysicalState       = "Foo Bar"
	storePhysicalPostalCode  = "00000-000"
	storePhysicalPhoneNumber = "(00)00000-0000"
	storePixKey              = "000.000.000-00"
)

type Seeder interface {
	Seed(db *gorm.DB) error
}

type UserAdminSeeder struct{}

func (s UserAdminSeeder) Seed(db *gorm.DB) error {
	admin := &models.User{
		Email:    adminEmail,
		Password: adminPassword,
		IsActive: true,
		IsStaff:  true,
	}
	return createRecordIfNotExists(db, admin, "email = ?", admin.Email, "AdminUser")
}

type StoreConfigSeeder struct{}

func (s StoreConfigSeeder) Seed(db *gorm.DB) error {
	storeConfig := &models.StoreConfig{
		DeliveryPrice:            10,
		DeliveryIsActive:         true,
		PhysicalStoreEmail:       storePhysicalEmail,
		PhysicalStoreAddress:     storePhysicalAddress,
		PhysicalStoreCity:        storePhysicalCity,
		PhysicalStoreState:       storePhysicalState,
		PhysicalStorePostalCode:  storePhysicalPostalCode,
		PhysicalStorePhoneNumber: storePhysicalPhoneNumber,
		PaymentCashIsActive:      true,
		PaymentPixIsActive:       true,
		PixKey:                   storePixKey,
		PixKeyType:               models.PixTypeCPF,
	}
	return createRecordIfNotExists(db, storeConfig, "physical_store_email = ?", storeConfig.PhysicalStoreEmail, "StoreConfig")
}

func createRecordIfNotExists(db *gorm.DB, value interface{}, query string, args ...interface{}) error {
	if err := db.FirstOrCreate(value, append([]interface{}{query}, args...)...).Error; err != nil {
		log.Printf("Failed to create %T: %v", value, err)
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
