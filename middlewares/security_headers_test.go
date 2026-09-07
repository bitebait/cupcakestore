package middlewares

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestBrowserSecurityHeadersPreserveAssetCaching(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())
	app.Get("/account", func(c fiber.Ctx) error { c.Type("html"); return c.SendString("<p>Conta</p>") })
	app.Get("/data", func(c fiber.Ctx) error { return c.JSON(fiber.Map{"count": 1}) })
	app.Get("/style.css", func(c fiber.Ctx) error { c.Type("css"); return c.SendString("body{}") })
	for _, path := range []string{"/account", "/data", "/style.css"} {
		response, err := app.Test(httptest.NewRequest("GET", path, nil))
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		policy := response.Header.Get("Content-Security-Policy")
		for _, directive := range []string{"script-src 'self'", "frame-ancestors 'none'", "form-action 'self'"} {
			if !strings.Contains(policy, directive) {
				t.Errorf("%s: missing %s", path, directive)
			}
		}
		if strings.Contains(policy, "unsafe-inline") || strings.Contains(policy, "unsafe-eval") {
			t.Fatal("unsafe script exception")
		}
		if response.Header.Get("X-Content-Type-Options") != "nosniff" || response.Header.Get("X-Frame-Options") != "DENY" {
			t.Fatal("missing browser protections")
		}
		private := response.Header.Get("Cache-Control") == "private, no-store"
		if private != (path != "/style.css") {
			t.Fatalf("%s: inappropriate caching policy", path)
		}
		if response.Header.Get("Strict-Transport-Security") != "" {
			t.Fatal("HTTP development must not advertise HSTS")
		}
	}
}
