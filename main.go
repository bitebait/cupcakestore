package main

import (
	"fmt"
	"github.com/bitebait/cupcakestore/bootstrap"
	"github.com/bitebait/cupcakestore/config"
	"log"
)

const (
	envCertFilePath = "CERT_FILE_PATH"
	envKeyFilePath  = "KEY_FILE_PATH"
	envAppHost      = "APP_HOST"
	envAppPort      = "APP_PORT"
	envDevMode      = "DEV_MODE"
	defaultHost     = "localhost"
	defaultPort     = "4000"
	defaultDevMode  = "true"
)

// Extracted function to get configuration values for better readability
func getConfigValue(cfg *config.Config, key, def string) string {
	return cfg.GetEnvVar(key, def)
}

func main() {
	app := bootstrap.NewApplication()
	cfg := config.Instance()

	certFilePath := getConfigValue(cfg, envCertFilePath, "")
	keyFilePath := getConfigValue(cfg, envKeyFilePath, "")
	host := getConfigValue(cfg, envAppHost, defaultHost)
	port := getConfigValue(cfg, envAppPort, defaultPort)
	addr := fmt.Sprintf("%s:%s", host, port)

	if isDevelopmentMode(cfg) {
		log.Fatal(app.Listen(addr))
	} else {
		log.Fatal(app.ListenTLS(addr, certFilePath, keyFilePath))
	}
}

func isDevelopmentMode(cfg *config.Config) bool {
	return getConfigValue(cfg, envDevMode, defaultDevMode) == "true"
}
