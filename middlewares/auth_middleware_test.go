package middlewares

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func TestSessionAuthorizationUsesCurrentAccountAndStopsDeniedRequests(t *testing.T) {
	cases := []struct {
		name          string
		active, staff bool
		lookupErr     error
		status        int
		allowed       bool
	}{
		{"demoted administrator", true, false, nil, 403, false},
		{"inactive administrator", false, true, nil, 302, false},
		{"deleted account", false, false, gorm.ErrRecordNotFound, 302, false},
		{"database unavailable", false, false, errors.New("offline"), 500, false},
		{"active administrator", true, true, nil, 200, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			session.SetupSession()
			app := fiber.New()
			app.Use(session.Middleware)
			app.Get("/session", func(c fiber.Ctx) error {
				sess := session.FromContext(c)
				sess.Set(session.UserIDKey, uint(7))
				return c.SendStatus(200)
			})
			called := false
			lookups := 0
			app.Use(Auth(func(id uint) (models.Profile, error) {
				lookups++
				if id != 7 {
					t.Errorf("lookup ID = %d", id)
				}
				return models.Profile{UserID: 7, User: models.User{IsActive: tc.active, IsStaff: tc.staff}}, tc.lookupErr
			}))
			app.Get("/admin", LoginRequired(), LoginAndStaffRequired(), func(c fiber.Ctx) error { called = true; return c.SendStatus(200) })
			login, err := app.Test(httptest.NewRequest("GET", "/session", nil))
			if err != nil {
				t.Fatal(err)
			}
			defer login.Body.Close()
			request := httptest.NewRequest("GET", "/admin", nil)
			for _, cookie := range login.Cookies() {
				request.AddCookie(cookie)
			}
			response, err := app.Test(request)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if lookups != 1 {
				t.Fatalf("account lookups = %d, want one per request", lookups)
			}
			if response.StatusCode != tc.status || called != tc.allowed {
				t.Fatalf("status=%d handlerCalled=%v", response.StatusCode, called)
			}
		})
	}
}
