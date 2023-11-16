package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	AppHost string
	AppPort string
	DbPath  string
	Payment PaymentConfig
}

type PaymentConfig struct {
	IdPix      string
	IdCash     string
	PixActive  bool
	CashActive bool
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			AppHost: getEnvOrDefault("APP_HOST", "localhost"),
			AppPort: getEnvOrDefault("APP_PORT", "8080"),
			DbPath:  getEnvOrDefault("DB_PATH", "db/app.db"),
			Payment: PaymentConfig{
				IdPix:      getEnvOrDefault("PAYMENT_ID_PIX", ""),
				IdCash:     getEnvOrDefault("PAYMENT_ID_CASH", ""),
				PixActive:  getEnvAsBool("PAYMENT_PIX_ACTIVE", false),
				CashActive: getEnvAsBool("PAYMENT_CASH_ACTIVE", false),
			},
		}
	})

	return instance
}

func getEnvOrDefault(key, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return def
}

func getEnvAsBool(key string, def bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return def
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return def
	}
	return boolValue
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Failed to load .env file. Using default or OS environment variables.")
	}
}
