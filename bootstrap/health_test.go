package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestHealthChecksReportDatabaseFailureWithoutSessionCookies(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "health.db")), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	app := fiber.New()
	registerHealthChecks(app, db)
	// Any accidental fallthrough would reach session/auth middleware in the app.
	app.Use(func(c fiber.Ctx) error { return c.SendStatus(http.StatusUnauthorized) })
	check := func(path string, want int) {
		t.Helper()
		for _, method := range []string{http.MethodGet, http.MethodHead} {
			response, err := app.Test(httptest.NewRequest(method, path, nil))
			if err != nil {
				t.Fatal(err)
			}
			response.Body.Close()
			if response.StatusCode != want || len(response.Cookies()) != 0 {
				t.Fatalf("%s %s: status=%d cookies=%v", method, path, response.StatusCode, response.Cookies())
			}
		}
	}
	check("/readyz", http.StatusOK)
	check("/livez", http.StatusOK)
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
	check("/readyz", http.StatusServiceUnavailable)
	check("/livez", http.StatusOK)
}
