package session

import (
	"time"

	"github.com/bitebait/cupcakestore/models"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const SessionExpiration = 1 * time.Hour

var Store *session.Store

func SetupSession() {
	sessConfig := session.Config{
		Expiration: SessionExpiration,
	}
	Store = session.New(sessConfig)
	Store.RegisterType(&models.Profile{})
}
