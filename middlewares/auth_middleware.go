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
	return sessionHandler(false, false)
}

func LoginRequired() fiber.Handler {
	return sessionHandler(true, false)
}

func LoginAndStaffRequired() fiber.Handler {
	return sessionHandler(true, true)
}

func sessionHandler(loginRequired, staffRequired bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := session.Store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
		}

		profile, ok := sess.Get("Profile").(*models.Profile)
		if ok {
			err := processAuthenticatedUser(c, profile)
			if err != nil {
				return err
			}
			if loginRequired {
				if staffRequired && !profile.User.IsStaff {
					return c.Redirect("/auth/logout")
				}
				return c.Next()
			}
		}

		if loginRequired {
			return c.Redirect("/auth/login")
		}

		return c.Next()
	}
}

func processAuthenticatedUser(c *fiber.Ctx, profile *models.Profile) error {
	profileService := getProfileService()
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

func getProfileService() services.ProfileService {
	profileRepository := repositories.NewProfileRepository(database.DB)
	return services.NewProfileService(profileRepository)
}
