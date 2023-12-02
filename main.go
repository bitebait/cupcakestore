package main

import (
	"fmt"
	"log"

	"github.com/bitebait/cupcakestore/bootstrap"
	"github.com/bitebait/cupcakestore/config"
)

func main() {
	app := bootstrap.NewApplication()
	addr := fmt.Sprintf("%s:%s", config.GetEnv("APP_HOST", "localhost"), config.GetEnv("APP_PORT", "4000"))

	if config.GetEnv("DEV_MODE", "true") == "true" {
		log.Fatal(app.Listen(addr))
		return
	}

	certFile := "/etc/letsencrypt/live/cupcakestore.schwaab.me/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/cupcakestore.schwaab.me/privkey.pem"
	log.Fatal(app.ListenTLS(addr, certFile, keyFile))
}
