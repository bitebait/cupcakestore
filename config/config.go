package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppHost             string
	AppPort             string
	DBType              string
	DBPath              string
	DBHost              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBPort              string
	DBSSLMode           string
	DBTimezone          string
	RedirectAfterLogin  string
	RedirectAfterLogout string
	DevMode             bool
	CertFilePath        string
	KeyFilePath         string
}

var cfg *Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	cfg = &Config{
		AppHost:             os.Getenv("APP_HOST"),
		AppPort:             os.Getenv("APP_PORT"),
		DBType:              os.Getenv("DB_TYPE"),
		DBPath:              os.Getenv("DB_PATH"),
		DBHost:              os.Getenv("DB_HOST"),
		DBUser:              os.Getenv("DB_USER"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBName:              os.Getenv("DB_NAME"),
		DBPort:              os.Getenv("DB_PORT"),
		DBSSLMode:           os.Getenv("DB_SSLMODE"),
		DBTimezone:          os.Getenv("DB_TIMEZONE"),
		RedirectAfterLogin:  os.Getenv("REDIRECT_AFTER_LOGIN"),
		RedirectAfterLogout: os.Getenv("REDIRECT_AFTER_LOGOUT"),
		DevMode:             mustParseBool(os.Getenv("DEV_MODE")),
		CertFilePath:        os.Getenv("CERT_FILE_PATH"),
		KeyFilePath:         os.Getenv("KEY_FILE_PATH"),
	}
}

func Get() *Config {
	return cfg
}

func mustParseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatalf("unable to parse bool: %e", err)
	}
	return b
}
