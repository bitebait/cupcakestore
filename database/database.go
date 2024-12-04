package database

import (
	"fmt"
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var (
	DB           *gorm.DB
	dbLoggerMode = logger.Silent
)

func SetupDatabase() {
	dbType := config.Get().DBType
	switch dbType {
	case "sqlite":
		setupSQLiteDatabase()
	case "postgres":
		setupPostgresDatabase()
	default:
		log.Panicf("Tipo de banco de dados não suportado: %s", dbType)
	}
}

func setupSQLiteDatabase() {
	db, err := gorm.Open(sqlite.Open(config.Get().DBPath), &gorm.Config{})
	handleError("Falha ao conectar ao banco de dados SQLite", err)

	handleError("Falha ao migrar os modelos", migrateModels(db))
	handleError("Falha ao popular o banco de dados", SeedDatabase(db))

	db.Logger = logger.Default.LogMode(dbLoggerMode)
	DB = db
}

func setupPostgresDatabase() {
	cfg := config.Get()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	handleError("Falha ao conectar ao banco de dados PostgreSQL", err)

	handleError("Falha ao migrar os modelos", migrateModels(db))
	handleError("Falha ao popular o banco de dados", SeedDatabase(db))

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
