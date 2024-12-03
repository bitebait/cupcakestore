package database

import (
	"fmt"
	"log"
	"os"

	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB           *gorm.DB
	dbLoggerMode = logger.Silent
)

func SetupDatabase() {
	dbType := config.Instance().GetEnvVar("DB_TYPE", "sqlite")
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
	dbPath := config.Instance().GetEnvVar("DB_PATH", "gorm.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	handleError("Falha ao conectar ao banco de dados SQLite", err)

	handleError("Falha ao migrar os modelos", migrateModels(db))
	handleError("Falha ao popular o banco de dados", SeedDatabase(db))

	db.Logger = logger.Default.LogMode(dbLoggerMode)
	DB = db
}

func setupPostgresDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_TIMEZONE"),
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
