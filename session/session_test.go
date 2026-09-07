package session

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

func TestLoginRotatesSessionAndRejectsPreviousCSRFToken(t *testing.T) {
	SetupSession()
	app := fiber.New()
	app.Use(Middleware)
	app.Use(csrf.New(csrf.Config{Session: Store, Extractor: extractors.FromForm("_csrf")}))
	app.Get("/form", func(c fiber.Ctx) error { return c.SendString(csrf.TokenFromContext(c)) })
	app.Get("/identity", func(c fiber.Ctx) error {
		if FromContext(c).Get(UserIDKey) != uint(7) {
			return fiber.ErrUnauthorized
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
	app.Post("/login", func(c fiber.Ctx) error {
		if err := Login(c, 7); err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
	app.Post("/protected", func(c fiber.Ctx) error {
		if FromContext(c).Get(UserIDKey) != uint(7) {
			return fiber.ErrUnauthorized
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
	cookies := map[string]*http.Cookie{}
	request := func(method, path, token string, want int) string {
		t.Helper()
		form := url.Values{"_csrf": {token}}.Encode()
		req := httptest.NewRequest(method, path, strings.NewReader(form))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != want {
			t.Fatalf("%s %s: status %d, want %d: %s", method, path, res.StatusCode, want, body)
		}
		for _, cookie := range res.Cookies() {
			cookies[cookie.Name] = cookie
		}
		return string(body)
	}
	oldToken := request(http.MethodGet, "/form", "", http.StatusOK)
	oldSession := *cookies["session_id"]
	request(http.MethodPost, "/login", oldToken, http.StatusNoContent)
	if cookies["session_id"].Value == oldSession.Value {
		t.Fatal("login did not rotate the session ID")
	}
	request(http.MethodPost, "/protected", oldToken, http.StatusForbidden)
	newToken := request(http.MethodGet, "/form", "", http.StatusOK)
	if newToken == oldToken {
		t.Fatal("login did not invalidate the previous CSRF token")
	}
	request(http.MethodPost, "/protected", newToken, http.StatusNoContent)
	request(http.MethodGet, "/identity", "", http.StatusNoContent)

	// Replaying the anonymous session ID must not reveal the logged-in identity.
	replay := httptest.NewRequest(http.MethodGet, "/identity", nil)
	replay.AddCookie(&oldSession)
	res, err := app.Test(replay)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("old session was authenticated: status %d", res.StatusCode)
	}
}
