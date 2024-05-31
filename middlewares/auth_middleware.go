package middlewares

import (
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/services"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		profileRespository := repositories.NewProfileRepository(database.DB)
		profileService := services.NewProfileService(profileRespository)
		sess, err := session.Store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
		}

		// Check if user is authenticated
		if profile := sess.Get("Profile"); profile != nil {
			p, err := profileService.FindByUserId(profile.(*models.Profile).UserID)
			if err != nil {
				return err
			}

			// Check if user is active
			if !p.User.IsActive {
				if err = sess.Destroy(); err != nil {
					return err
				}
			}

			c.Locals("Profile", &p)
			return c.Next()
		}

		// Check allowed paths for unauthenticated users
		allowedPaths := []string{"/auth/login", "/store", "/auth/register", "/"}
		for _, path := range allowedPaths {
			if c.Path() == path {
				return c.Next()
			}
		}

		// Redirect to login page for other paths
		return c.Redirect("/auth/login")
	}
}

func LoginAndStaffRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		profile, ok := c.Locals("Profile").(*models.Profile)
		if !(ok && profile != nil && profile.User.IsStaff && profile.User.IsActive) {
			return c.Redirect("/auth/logout")
		}

		return c.Next()
	}
}
