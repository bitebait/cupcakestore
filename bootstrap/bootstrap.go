package bootstrap

import (
	"github.com/Masterminds/sprig/v3"
	minifier "github.com/beyer-stefan/gofiber-minifier"
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/routers"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func NewApplication() *fiber.App {
	setupDatabase()
	setupSession()

	app := createFiberApp()
	registerMiddlewares(app)
	serveStaticFiles(app)
	configureHTTPS(app)
	registerRoutes(app)

	return app
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
	app.Use(favicon.New(favicon.Config{
		File: "./web/dist/img/favicon.png",
		URL:  "/favicon.ico",
	}))
	app.Use(minifier.New(minifier.Config{
		MinifyHTML: true,
		MinifyCSS:  true,
		MinifyJS:   true,
	}))
}

func serveStaticFiles(app *fiber.App) {
	app.Static("/", "./web")
}

func configureHTTPS(app *fiber.App) {
	if !isDevMode() {
		app.Use(redirectToHTTPS)
	}
}

func registerRoutes(app *fiber.App) {
	app.Use(middlewares.Auth())
	routers.InstallRouters(app)
}

func isDevMode() bool {
	return config.Instance().GetEnvVar("DEV_MODE", "true") == "true"
}

func redirectToHTTPS(c *fiber.Ctx) error {
	if c.Protocol() == "http" {
		return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
	}
	return c.Next()
}
