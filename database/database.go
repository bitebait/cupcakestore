package database

import (
	"github.com/bitebait/cupcakestore/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.GetEnv("DB_PATH", "database.db")), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
}
