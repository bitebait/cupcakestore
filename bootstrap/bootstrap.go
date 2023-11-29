package bootstrap

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/routers"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

// NewApplication creates a new instance of the Fiber application.
func NewApplication() *fiber.App {
	// Config
	config.SetupEnvFile()

	// Database
	database.SetupDatabase()

	// Session
	session.SetupSession()

	// Fiber
	engine := html.New("./views", ".html")
	engine.AddFuncMap(sprig.FuncMap())

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Static("/", "./web")
	app.Use(middlewares.Auth())

	// Routes
	routers.InstallRouters(app)

	return app
}
