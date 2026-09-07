package middlewares

import (
	"errors"

	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Auth() fiber.Handler                  { return createSessionHandler(false, false) }
func LoginRequired() fiber.Handler         { return createSessionHandler(true, false) }
func LoginAndStaffRequired() fiber.Handler { return createSessionHandler(true, true) }

func createSessionHandler(requireLogin, requireStaff bool) fiber.Handler {
	return sessionHandler(requireLogin, requireStaff, func(id uint) (models.Profile, error) {
		return repositories.NewProfileRepository(database.DB).FindByUserId(id)
	})
}

// Reload account permissions from the database; a session is only proof of identity.
func sessionHandler(requireLogin, requireStaff bool, lookup func(uint) (models.Profile, error)) fiber.Handler {
	return func(c fiber.Ctx) error {
		profile, ok := c.Locals("Profile").(*models.Profile)
		if !ok || profile == nil {
			sess := session.FromContext(c)
			if sess == nil {
				return fiber.ErrInternalServerError
			}
			userID, authenticated := sess.Get(session.UserIDKey).(uint)
			if authenticated && userID != 0 {
				current, err := lookup(userID)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return fiber.ErrInternalServerError
				}
				if err != nil || !current.User.IsActive {
					if err := sess.Reset(); err != nil {
						return fiber.ErrInternalServerError
					}
					return c.Redirect().Status(fiber.StatusFound).To("/auth/login")
				}
				current.User.Password = ""
				profile = &current
				c.Locals("Profile", profile)
			}
		}
		if profile == nil {
			if requireLogin {
				return c.Redirect().Status(fiber.StatusFound).To("/auth/login")
			}
			return c.Next()
		}
		if requireStaff && !profile.User.IsStaff {
			return fiber.NewError(fiber.StatusForbidden, "acesso negado")
		}
		return c.Next()
	}
}
