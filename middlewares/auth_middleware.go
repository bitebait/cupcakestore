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
	return createSessionHandler(false, false)
}

func LoginRequired() fiber.Handler {
	return createSessionHandler(true, false)
}

func LoginAndStaffRequired() fiber.Handler {
	return createSessionHandler(true, true)
}

func createSessionHandler(requireLogin, requireStaff bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := session.Store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
		}
		profile, isAuthenticated := sess.Get("Profile").(*models.Profile)
		if isAuthenticated {
			err := handleAuthenticatedUser(c, profile, requireLogin, requireStaff)
			if err != nil {
				return err
			}
			return c.Next()
		}
		if requireLogin {
			return c.Redirect("/auth/login")
		}
		return c.Next()
	}
}

func handleAuthenticatedUser(c *fiber.Ctx, profile *models.Profile, requireLogin, requireStaff bool) error {
	err := updateProfileIfNeeded(c, profile)
	if err != nil {
		return err
	}
	if requireStaff && !profile.User.IsStaff {
		return c.Redirect("/auth/logout")
	}
	return nil
}

func updateProfileIfNeeded(c *fiber.Ctx, profile *models.Profile) error {
	profileService := fetchProfileService()
	updatedProfile, err := profileService.FindByUserId(profile.UserID)
	if err != nil {
		return err
	}
	if !updatedProfile.User.IsActive {
		sess, _ := session.Store.Get(c)
		if err := sess.Destroy(); err != nil {
			return err
		}
		return c.Redirect("/auth/login")
	}
	c.Locals("Profile", &updatedProfile)
	return nil
}

func fetchProfileService() services.ProfileService {
	profileRepo := repositories.NewProfileRepository(database.DB)
	return services.NewProfileService(profileRepo)
}
