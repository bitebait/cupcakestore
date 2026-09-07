package config

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

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
	AdminEmail          string
	AdminPassword       string
}

var (
	cfg     *Config
	cfgErr  error
	cfgOnce sync.Once
)

// Initialize loads configuration once, allowing startup errors to reach main.
func Initialize() error {
	cfgOnce.Do(func() { cfg, cfgErr = Load() })
	return cfgErr
}

func Get() *Config {
	if err := Initialize(); err != nil {
		panic(err)
	}
	return cfg
}

// Load accepts environment-only deployments; an optional .env never overrides
// variables already supplied by the process environment.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}
	return loadEnvironment()
}

func loadEnvironment() (*Config, error) {
	devMode, err := strconv.ParseBool(envOrDefault("DEV_MODE", "true"))
	if err != nil {
		return nil, errors.New("DEV_MODE must be true or false")
	}
	c := &Config{
		AppHost:             envOrDefault("APP_HOST", "0.0.0.0"),
		AppPort:             envOrDefault("APP_PORT", "8080"),
		DBType:              envOrDefault("DB_TYPE", "sqlite"),
		DBPath:              envOrDefault("DB_PATH", "gorm.db"),
		DBHost:              os.Getenv("DB_HOST"),
		DBUser:              os.Getenv("DB_USER"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBName:              os.Getenv("DB_NAME"),
		DBPort:              envOrDefault("DB_PORT", "5432"),
		DBSSLMode:           envOrDefault("DB_SSLMODE", "require"),
		DBTimezone:          envOrDefault("DB_TIMEZONE", "America/Sao_Paulo"),
		RedirectAfterLogin:  envOrDefault("REDIRECT_AFTER_LOGIN", "/store"),
		RedirectAfterLogout: envOrDefault("REDIRECT_AFTER_LOGOUT", "/store"),
		DevMode:             devMode,
		CertFilePath:        os.Getenv("CERT_FILE_PATH"),
		KeyFilePath:         os.Getenv("KEY_FILE_PATH"),
		AdminEmail:          strings.TrimSpace(os.Getenv("ADMIN_EMAIL")),
		AdminPassword:       os.Getenv("ADMIN_PASSWORD"),
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Validate() error {
	if err := validatePort("APP_PORT", c.AppPort); err != nil {
		return err
	}
	switch c.DBType {
	case "sqlite":
		if c.DBPath == "" {
			return errors.New("DB_PATH is required for SQLite")
		}
	case "postgres":
		if c.DBHost == "" || c.DBUser == "" || c.DBName == "" {
			return errors.New("DB_HOST, DB_USER and DB_NAME are required for PostgreSQL")
		}
		if err := validatePort("DB_PORT", c.DBPort); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported DB_TYPE %q", c.DBType)
	}
	for _, path := range []string{c.RedirectAfterLogin, c.RedirectAfterLogout} {
		u, err := url.Parse(path)
		if err != nil || !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") || strings.ContainsAny(path, "\\\r\n") || u.IsAbs() || u.Host != "" {
			return errors.New("redirect paths must be local absolute paths")
		}
	}
	if !c.DevMode && (c.CertFilePath == "" || c.KeyFilePath == "") {
		return errors.New("CERT_FILE_PATH and KEY_FILE_PATH are required when DEV_MODE=false")
	}
	if (c.AdminEmail == "") != (c.AdminPassword == "") {
		return errors.New("ADMIN_EMAIL and ADMIN_PASSWORD must be configured together")
	}
	if c.AdminEmail != "" {
		address, err := mail.ParseAddress(c.AdminEmail)
		if err != nil || address.Address != c.AdminEmail {
			return errors.New("ADMIN_EMAIL must be a valid email address")
		}
		if len(c.AdminPassword) < 12 || len(c.AdminPassword) > 72 {
			return errors.New("ADMIN_PASSWORD must contain between 12 and 72 bytes")
		}
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func validatePort(key, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("%s must be a port between 1 and 65535", key)
	}
	return nil
}
