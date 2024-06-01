package main

import (
	"fmt"
	"log"

	"github.com/bitebait/cupcakestore/bootstrap"
	"github.com/bitebait/cupcakestore/config"
)

func main() {
	app := bootstrap.NewApplication()
	cfg := config.Instance()

	host := cfg.GetEnvVar("APP_HOST", "localhost")
	port := cfg.GetEnvVar("APP_PORT", "4000")
	addr := fmt.Sprintf("%s:%s", host, port)

	if cfg.GetEnvVar("DEV_MODE", "true") == "true" {
		log.Fatal(app.Listen(addr))
		return
	}

	certFile := "/etc/letsencrypt/live/cupcakestore.schwaab.me/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/cupcakestore.schwaab.me/privkey.pem"
	log.Fatal(app.ListenTLS(addr, certFile, keyFile))
}
