package main

import (
	"fmt"
	"log"

	"github.com/bitebait/cupcakestore/bootstrap"
	"github.com/bitebait/cupcakestore/config"
)

func main() {
	app := bootstrap.NewApplication()
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", config.GetEnv("APP_HOST", "localhost"), config.GetEnv("APP_PORT", "4000"))))
}
