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

const (
	faviconPath = "./web/dist/img/favicon.png"
	faviconURL  = "/favicon.ico"
)

func NewApplication() *fiber.App {
	database.SetupDatabase()
	session.SetupSession()

	fiberApp := createFiberApp()
	registerMiddlewares(fiberApp)
	serveStaticFiles(fiberApp)
	configureHTTPS(fiberApp)
	registerRoutes(fiberApp)
	return fiberApp
}

func createFiberApp() *fiber.App {
	engine := setupTemplateEngine()

	return fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})
}

func setupTemplateEngine() *html.Engine {
	engine := html.New("./views", ".html")
	engine.AddFuncMap(sprig.FuncMap())
	engine.Reload(true)
	return engine
}

func registerMiddlewares(fiberApp *fiber.App) {
	fiberApp.Use(recover.New())
	fiberApp.Use(logger.New())
	fiberApp.Use(compress.New(compress.Config{Level: compress.LevelBestCompression}))
	fiberApp.Use(cors.New(cors.Config{AllowOrigins: "*"}))
	fiberApp.Use(favicon.New(favicon.Config{File: faviconPath, URL: faviconURL}))
	fiberApp.Use(minifier.New(minifier.Config{MinifyHTML: true, MinifyCSS: true, MinifyJS: true}))
}

func serveStaticFiles(fiberApp *fiber.App) {
	fiberApp.Static("/", "./web")
}

func configureHTTPS(fiberApp *fiber.App) {
	if !isDevMode() {
		fiberApp.Use(redirectToHTTPS)
	}
}

func registerRoutes(fiberApp *fiber.App) {
	fiberApp.Use(middlewares.Auth())
	routers.InstallRouters(fiberApp)
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
