package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bitebait/cupcakestore/bootstrap"
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/database"
	"github.com/gofiber/fiber/v3"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	app, err := bootstrap.NewApplicationWithError()
	if err != nil {
		return err
	}
	sqlDB, err := database.DB.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	cfg := config.Get()
	addr := net.JoinHostPort(cfg.AppHost, cfg.AppPort)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	errors := make(chan error, 1)
	go func() {
		if cfg.DevMode {
			errors <- app.Listen(addr)
		} else {
			errors <- app.Listen(addr, fiber.ListenConfig{
				CertFile: cfg.CertFilePath, CertKeyFile: cfg.KeyFilePath,
			})
		}
	}()
	select {
	case err := <-errors:
		return err
	case <-ctx.Done():
		return app.ShutdownWithTimeout(10 * time.Second)
	}
}
