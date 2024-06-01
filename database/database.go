package database

import (
	"log"

	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func SetupDatabase() {
	dbPath := config.Instance().GetEnvVar("DB_PATH", "database.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect to the database: %v", err)
	}

	if err := migrateModels(db); err != nil {
		log.Panicf("Failed to migrate models: %v", err)
	}

	if err := SeedDatabase(db); err != nil {
		log.Panicf("Failed to seed database: %v", err)
	}

	db.Logger = logger.Default.LogMode(logger.Silent)
	DB = db
}

func migrateModels(db *gorm.DB) error {
	m := []interface{}{
		&models.User{},
		&models.Profile{},
		&models.Product{},
		&models.Stock{},
		&models.StoreConfig{},
		&models.Order{},
		&models.OrderDeliveryDetail{},
		&models.ShoppingCart{},
		&models.ShoppingCartItem{},
	}

	return db.AutoMigrate(m...)
}
