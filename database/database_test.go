package database

import (
	"net/url"
	"path/filepath"
	"testing"

	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

func testDatabase(t *testing.T, cfg *config.Config) *gorm.DB {
	t.Helper()
	cfg.DBType = "sqlite"
	cfg.DBPath = filepath.Join(t.TempDir(), "store.db")
	db, err := Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}

func TestSeedWithoutDefaultAdministrator(t *testing.T) {
	cfg := &config.Config{}
	db := testDatabase(t, cfg)
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("unexpected default users: %d", count)
	}
	var store models.StoreConfig
	if err := db.First(&store).Error; err != nil {
		t.Fatal(err)
	}
	if store.PaymentPixIsActive || store.PixKey != "" {
		t.Fatal("PIX must remain unavailable until configured")
	}
	if err := db.Model(&store).Update("physical_store_email", "shop@example.com").Error; err != nil {
		t.Fatal(err)
	}
	if err := seedDatabase(db, cfg); err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&models.StoreConfig{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("restarting after an email edit created %d store configs", count)
	}
}

func TestSeedAdministratorPreservesExistingAccount(t *testing.T) {
	cfg := &config.Config{AdminEmail: "admin@example.com", AdminPassword: "strong-test-password"}
	db := testDatabase(t, cfg)
	var admin models.User
	if err := db.First(&admin).Error; err != nil {
		t.Fatal(err)
	}
	if !admin.IsStaff || admin.Password == cfg.AdminPassword || admin.CheckPassword(cfg.AdminPassword) != nil {
		t.Fatal("administrator password must be hashed and usable")
	}
	originalHash := admin.Password
	cfg.AdminPassword = "different-test-password"
	if err := db.Model(&admin).Update("is_staff", false).Error; err != nil {
		t.Fatal(err)
	}
	if err := seedDatabase(db, cfg); err != nil {
		t.Fatal(err)
	}
	if err := db.First(&admin).Error; err != nil {
		t.Fatal(err)
	}
	if admin.Password != originalHash || admin.IsStaff {
		t.Fatal("seed changed existing account credentials or permissions")
	}
}

func TestSQLiteForeignKeysEnabled(t *testing.T) {
	db := testDatabase(t, &config.Config{})
	var enabled int
	if err := db.Raw("PRAGMA foreign_keys").Scan(&enabled).Error; err != nil {
		t.Fatal(err)
	}
	if enabled != 1 {
		t.Fatal("SQLite foreign keys must be enforced")
	}
}

func TestPostgresDSNPreservesCredentialsAndTimezone(t *testing.T) {
	cfg := &config.Config{
		DBHost: "::1", DBPort: "5432", DBUser: "a user", DBPassword: "space ' @ ? / &",
		DBName: "store name", DBSSLMode: "require", DBTimezone: "America/Sao_Paulo",
	}
	u, err := url.Parse(postgresDSN(cfg))
	if err != nil {
		t.Fatal(err)
	}
	password, _ := u.User.Password()
	if password != cfg.DBPassword || u.User.Username() != cfg.DBUser || u.Host != "[::1]:5432" || u.Path != "/store name" || u.Query().Get("TimeZone") != cfg.DBTimezone || u.Query().Get("sslmode") != "require" {
		t.Fatalf("PostgreSQL DSN lost or corrupted connection values: %s", u.Redacted())
	}
}
