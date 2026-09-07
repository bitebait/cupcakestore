package bootstrap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
)

func TestCSRFProtectsFormsAndUsesHttpOnlySessionCookie(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })
	session.SetupSession()
	app := fiber.New()
	registerMiddlewares(app)
	app.Get("/form", func(c *fiber.Ctx) error {
		return c.SendString(c.Locals("CSRFToken").(string))
	})
	app.Post("/form", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/form", nil))
	if err != nil {
		t.Fatal(err)
	}
	token, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil || len(token) == 0 {
		t.Fatalf("missing CSRF token: %v", err)
	}
	cookies := response.Cookies()
	foundSession := false
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			foundSession = true
			if !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode || cookie.Path != "/" {
				t.Fatalf("insecure session cookie options: %v", cookie)
			}
		}
	}
	if !foundSession {
		t.Fatal("CSRF token should be tied to a session cookie")
	}
	for _, tc := range []struct {
		name  string
		token string
		want  int
	}{
		{"missing token", "", fiber.StatusForbidden},
		{"forged token", "forged", fiber.StatusForbidden},
		{"valid token", string(token), fiber.StatusNoContent},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := url.Values{"_csrf": {tc.token}}.Encode()
			req := httptest.NewRequest(http.MethodPost, "/form", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			for _, cookie := range cookies {
				req.AddCookie(cookie)
			}
			response, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			response.Body.Close()
			if response.StatusCode != tc.want {
				t.Fatalf("status = %d, want %d", response.StatusCode, tc.want)
			}
		})
	}
}
