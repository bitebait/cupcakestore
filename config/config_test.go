package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func cleanEnvironment(t *testing.T) {
	t.Helper()
	for _, key := range []string{"APP_HOST", "APP_PORT", "DEV_MODE", "DB_TYPE", "DB_PATH", "DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "DB_SSLMODE", "DB_TIMEZONE", "REDIRECT_AFTER_LOGIN", "REDIRECT_AFTER_LOGOUT", "CERT_FILE_PATH", "KEY_FILE_PATH", "ADMIN_EMAIL", "ADMIN_PASSWORD"} {
		t.Setenv(key, "")
	}
}

func TestLoadWithoutDotEnv(t *testing.T) {
	cleanEnvironment(t)
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AppPort != "8080" || cfg.DBType != "sqlite" || !cfg.DevMode || cfg.AdminPassword != "" {
		t.Fatalf("unexpected local defaults: %#v", cfg)
	}
	t.Setenv("APP_PORT", "9090")
	if err := os.WriteFile(filepath.Join(".", ".env"), []byte("APP_PORT=1234\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err = Load()
	if err != nil || cfg.AppPort != "9090" {
		t.Fatalf("environment should take precedence: cfg=%v, err=%v", cfg, err)
	}
}

func TestInvalidEnvironment(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
		want string
	}{
		{"boolean", map[string]string{"DEV_MODE": "oops"}, "DEV_MODE"},
		{"port", map[string]string{"APP_PORT": "70000"}, "APP_PORT"},
		{"database", map[string]string{"DB_TYPE": "mysql"}, "DB_TYPE"},
		{"postgres", map[string]string{"DB_TYPE": "postgres"}, "DB_HOST"},
		{"tls", map[string]string{"DEV_MODE": "false"}, "CERT_FILE_PATH"},
		{"redirect host", map[string]string{"REDIRECT_AFTER_LOGIN": "//evil.test"}, "redirect"},
		{"redirect scheme", map[string]string{"REDIRECT_AFTER_LOGOUT": "https://evil.test"}, "redirect"},
		{"redirect backslash", map[string]string{"REDIRECT_AFTER_LOGIN": "/\\evil.test"}, "redirect"},
		{"incomplete admin", map[string]string{"ADMIN_EMAIL": "admin@example.com"}, "configured together"},
		{"admin email", map[string]string{"ADMIN_EMAIL": "invalid", "ADMIN_PASSWORD": "strong-password"}, "ADMIN_EMAIL"},
		{"short password", map[string]string{"ADMIN_EMAIL": "admin@example.com", "ADMIN_PASSWORD": "short"}, "ADMIN_PASSWORD"},
		{"long password", map[string]string{"ADMIN_EMAIL": "admin@example.com", "ADMIN_PASSWORD": strings.Repeat("a", 73)}, "ADMIN_PASSWORD"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cleanEnvironment(t)
			for key, value := range tc.env {
				t.Setenv(key, value)
			}
			_, err := loadEnvironment()
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("wanted error containing %q; got %v", tc.want, err)
			}
		})
	}
}
