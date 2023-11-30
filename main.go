package main

import (
	"fmt"
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

		log.Fatal(app.ListenTLS(fmt.Sprintf("%s:%s", host, port), certFile, keyFile))
	}
}
