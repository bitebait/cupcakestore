package session

import (
	"errors"
	"time"

	"github.com/bitebait/cupcakestore/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	fibersession "github.com/gofiber/fiber/v3/middleware/session"
)

const SessionExpiration = time.Hour
const UserIDKey = "user_id"

var (
	Store      *fibersession.Store
	Middleware fiber.Handler
)

func SetupSession() {
	Middleware, Store = fibersession.NewWithStore(fibersession.Config{
		IdleTimeout:     SessionExpiration,
		AbsoluteTimeout: 24 * time.Hour,
		CookieHTTPOnly:  true,
		CookieSecure:    !config.Get().DevMode,
		CookieSameSite:  "Lax",
		CookiePath:      "/",
	})
}

// FromContext returns the request session, whose save and release lifecycle is
// managed by Middleware. Callers must not manually save or release this session.
func FromContext(ctx fiber.Ctx) *fibersession.Middleware {
	return fibersession.FromContext(ctx)
}

// Login rotates the session and CSRF token before attaching an authenticated
// identity. Authorization remains a database lookup in the auth middleware.
func Login(ctx fiber.Ctx, userID uint) error {
	sess := FromContext(ctx)
	if sess == nil || userID == 0 {
		return errors.New("sessão ou usuário indisponível")
	}
	if err := sess.Regenerate(); err != nil {
		return err
	}
	if handler := csrf.HandlerFromContext(ctx); handler != nil {
		if err := handler.DeleteToken(ctx); err != nil {
			return err
		}
	}
	sess.Set(UserIDKey, userID)
	return nil
}
