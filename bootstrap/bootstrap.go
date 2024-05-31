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

func NewApplication() *fiber.App {
	// Setup configurations and dependencies
	setupConfigurations()
	setupDatabase()
	setupSession()

	// Create Fiber application with HTML template engine
	app := createFiberApp()

	// Middlewares
	registerMiddlewares(app)

	// Serve static files
	app.Static("/", "./web")

	// Redirect to HTTPS if not in DEV_MODE
	if !isDevMode() {
		app.Use(redirectToHTTPS)
	}

	// Auth middleware
	app.Use(middlewares.Auth())

	// Install routers
	routers.InstallRouters(app)

	return app
}

func setupConfigurations() {
	config.SetupEnvFile()
}

func setupDatabase() {
	database.SetupDatabase()
}

func setupSession() {
	session.SetupSession()
}

func createFiberApp() *fiber.App {
	engine := html.New("./views", ".html")
	engine.AddFuncMap(sprig.FuncMap())
	engine.Reload(true)

	return fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})
}

func registerMiddlewares(app *fiber.App) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
}

func isDevMode() bool {
	return config.GetEnv("DEV_MODE", "true") == "true"
}

func redirectToHTTPS(c *fiber.Ctx) error {
	if c.Protocol() == "http" {
		return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
	}
	return c.Next()
}
