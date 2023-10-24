package bootstrap

import (
	"github.com/Masterminds/sprig"
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/models"
	"github.com/bitebait/cupcakestore/routers"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"log"
)

func NewApplication() *fiber.App {
	// Config
	config.SetupEnvFile()

	// Database
	database.SetupDatabase()
	err := database.DB.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Product{},
	)
	if err != nil {
		log.Panic("error migrate models")
	}

	// Session
	session.SetupSession()

	// Fiber
	engine := html.New("./views", ".html")
	engine.AddFuncMap(sprig.FuncMap())
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	app.Static("/", "./web")
	// Routes
	routers.InstallRouters(app)

	return app
}
