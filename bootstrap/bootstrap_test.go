package bootstrap

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v3"
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
	app.Get("/form", func(c fiber.Ctx) error {
		return c.SendString(c.Locals("CSRFToken").(string))
	})
	app.Post("/form", func(c fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

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

func TestStaticCompressionDoesNotWritePublicAssets(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	if err := os.Mkdir(filepath.Join(directory, "web"), 0755); err != nil {
		t.Fatal(err)
	}
	content := strings.Repeat("body { color: black; }\n", 200)
	if err := os.WriteFile(filepath.Join(directory, "web", "style.css"), []byte(content), 0444); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(directory); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	app := fiber.New()
	serveStaticFiles(app)
	req := httptest.NewRequest(http.MethodGet, "/style.css", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	response, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK || response.Header.Get("Content-Encoding") != "gzip" {
		t.Fatalf("compressed asset: status=%d encoding=%q", response.StatusCode, response.Header.Get("Content-Encoding"))
	}
	reader, err := gzip.NewReader(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil || string(data) != content {
		t.Fatalf("decompressed asset differs: %v", err)
	}
	files, err := os.ReadDir(filepath.Join(directory, "web"))
	if err != nil || len(files) != 1 {
		t.Fatalf("compression must not create files: count=%d error=%v", len(files), err)
	}
}

func TestStaticMiddlewareServesAssetsAndLeavesDynamicPagesUncompressed(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	session.SetupSession()
	app := fiber.New()
	registerMiddlewares(app)
	serveStaticFiles(app)
	app.Get("/page", func(c fiber.Ctx) error {
		c.Type("html")
		return c.SendString(strings.Repeat("page content ", 200) + c.Locals("CSRFToken").(string))
	})
	for _, method := range []string{http.MethodGet, http.MethodHead} {
		res, err := app.Test(httptest.NewRequest(method, "/dist/img/logo.png", nil))
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusOK || !strings.HasPrefix(res.Header.Get("Content-Type"), "image/png") {
			t.Fatalf("%s asset: status %d, content type %q", method, res.StatusCode, res.Header.Get("Content-Type"))
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/page", nil)
	req.Header.Set("Accept-Encoding", "gzip, br")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusOK || res.Header.Get("Content-Encoding") != "" {
		t.Fatalf("dynamic page must remain uncompressed: status %d, encoding %q", res.StatusCode, res.Header.Get("Content-Encoding"))
	}
}

func TestHTTPSRedirectUsesSchemeAndPreservesHostPort(t *testing.T) {
	app := fiber.New()
	app.Use(redirectToHTTPS)
	app.Get("/store", func(c fiber.Ctx) error { return c.SendStatus(http.StatusOK) })
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8443/store?page=2", nil)
	// Forwarded headers from an untrusted client cannot bypass HTTPS handling.
	req.Header.Set("X-Forwarded-Proto", "https")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusMovedPermanently || res.Header.Get("Location") != "https://localhost:8443/store?page=2" {
		t.Fatalf("unexpected HTTPS redirect: status %d, location %q", res.StatusCode, res.Header.Get("Location"))
	}
}
