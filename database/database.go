package database

import (
	"log"

	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB           *gorm.DB
	dbLoggerMode = logger.Silent
)

func SetupDatabase() {
	dbPath := config.Instance().GetEnvVar("DB_PATH", "database.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	handleError("Failed to connect to the database", err)

	handleError("Failed to migrate models", migrateModels(db))
	handleError("Failed to seed database", SeedDatabase(db))

	db.Logger = logger.Default.LogMode(dbLoggerMode)
	DB = db
}

func handleError(msg string, err error) {
	if err != nil {
		log.Panicf("%s: %v", msg, err)
	}
}

func migrateModels(db *gorm.DB) error {
	modelsToMigrate := []interface{}{
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
	return db.AutoMigrate(modelsToMigrate...)
}
