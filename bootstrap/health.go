package bootstrap

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"gorm.io/gorm"
)

// Probes run before sessions and CSRF so polling never creates user sessions.
func registerHealthChecks(app *fiber.App, db *gorm.DB) {
	app.Add([]string{fiber.MethodGet, fiber.MethodHead}, healthcheck.LivenessEndpoint, healthcheck.New())
	app.Add([]string{fiber.MethodGet, fiber.MethodHead}, healthcheck.ReadinessEndpoint, healthcheck.New(healthcheck.Config{
		Probe: func(c fiber.Ctx) bool {
			sqlDB, err := db.DB()
			if err != nil {
				return false
			}
			ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
			defer cancel()
			return sqlDB.PingContext(ctx) == nil
		},
	}))
}
