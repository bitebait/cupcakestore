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
	var err error
	DB, err = gorm.Open(sqlite.Open(config.GetEnv("DB_PATH", "database.db")), &gorm.Config{})
	if err != nil {
		panic("Falha ao conectar ao banco de dados: " + err.Error())
	}

	migrateModels(DB)
	seedDatabase(DB)
	DB.Logger = logger.Default.LogMode(logger.Silent)
}

func migrateModels(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Product{},
		&models.Stock{},
		&models.StoreConfig{},
		&models.Order{},
		&models.OrderDeliveryDetail{},
		&models.ShoppingCart{},
		&models.ShoppingCartItem{},
	)
	if err != nil {
		log.Panic("erro ao migrar os modelos")
	}
}
