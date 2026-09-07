package database

import (
	"fmt"
	"net"
	"net/url"

	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// SetupDatabase is retained for callers using the original bootstrap API.
func SetupDatabase() {
	db, err := Open(config.Get())
	if err != nil {
		panic(err)
	}
	DB = db
}

// Open connects, migrates and seeds one database, returning startup failures.
func Open(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.DBType {
	case "sqlite":
		dialector = sqlite.Open(cfg.DBPath)
	case "postgres":
		dialector = postgres.Open(postgresDSN(cfg))
	default:
		return nil, fmt.Errorf("unsupported database type %q", cfg.DBType)
	}
	db, err := gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection: %w", err)
	}
	// SQLite allows only one writer. One connection also keeps in-memory test
	// databases consistent and ensures PRAGMA settings apply to every query.
	if cfg.DBType == "sqlite" {
		sqlDB.SetMaxOpenConns(1)
		for _, pragma := range []string{"PRAGMA foreign_keys = ON", "PRAGMA busy_timeout = 5000"} {
			if err := db.Exec(pragma).Error; err != nil {
				_ = sqlDB.Close()
				return nil, fmt.Errorf("configure SQLite: %w", err)
			}
		}
	}
	if err := migrateModels(db); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("migrate database: %w", err)
	}
	if err := seedDatabase(db, cfg); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("seed database: %w", err)
	}
	return db, nil
}

func postgresDSN(cfg *config.Config) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   net.JoinHostPort(cfg.DBHost, cfg.DBPort),
		Path:   "/" + cfg.DBName,
	}
	query := u.Query()
	query.Set("sslmode", cfg.DBSSLMode)
	query.Set("TimeZone", cfg.DBTimezone)
	u.RawQuery = query.Encode()
	return u.String()
}

func migrateModels(db *gorm.DB) error {
	return db.AutoMigrate(
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
}
