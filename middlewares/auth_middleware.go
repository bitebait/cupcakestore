package middlewares

import (
	"errors"

	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
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
	return func(c *fiber.Ctx) error {
		profile, ok := c.Locals("Profile").(*models.Profile)
		if !ok || profile == nil {
			sess, err := session.Store.Get(c)
			if err != nil {
				return fiber.ErrInternalServerError
			}
			stored, authenticated := sess.Get("Profile").(*models.Profile)
			if authenticated && stored != nil && stored.UserID != 0 {
				current, err := lookup(stored.UserID)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return fiber.ErrInternalServerError
				}
				if err != nil || !current.User.IsActive {
					if err := sess.Destroy(); err != nil {
						return fiber.ErrInternalServerError
					}
					return c.Redirect("/auth/login")
				}
				current.User.Password = ""
				profile = &current
				c.Locals("Profile", profile)
			}
		}
		if profile == nil {
			if requireLogin {
				return c.Redirect("/auth/login")
			}
			return c.Next()
		}
		if requireStaff && !profile.User.IsStaff {
			return fiber.NewError(fiber.StatusForbidden, "acesso negado")
		}
		return c.Next()
	}
}
