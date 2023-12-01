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

	// Static files
	app.Static("/", "./web")

	// Redirect to HTTPS if not in DEV_MODE
	if config.GetEnv("DEV_MODE", "true") == "false" {
		app.Use(func(c *fiber.Ctx) error {
			if c.Protocol() == "http" {
				return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
			}
			return c.Next()
		})
	}

	// Auth middleware
	app.Use(middlewares.Auth())

	// Install routers
	routers.InstallRouters(app)

	return app
}
