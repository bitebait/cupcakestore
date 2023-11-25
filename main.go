package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"

	"github.com/bitebait/cupcakestore/bootstrap"
	"github.com/bitebait/cupcakestore/config"
)

func main() {
	app := bootstrap.NewApplication()

	host := config.GetEnv("APP_HOST", "localhost")
	port := config.GetEnv("APP_PORT", "4000")

	if config.GetEnv("DEV_MODE", "true") == "true" {
		log.Fatal(app.Listen(fmt.Sprintf("%s:%s", host, port)))
	} else {
		certFile := "/etc/letsencrypt/live/cupcakestore.schwaab.me/fullchain.pem"
		keyFile := "/etc/letsencrypt/live/cupcakestore.schwaab.me/privkey.pem"

		app.Use(func(c *fiber.Ctx) error {
			if c.Protocol() == "http" {
				return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
			}
			return c.Next()
		})

		log.Fatal(app.ListenTLS(fmt.Sprintf("%s:%s", host, port), certFile, keyFile))
	}
}
