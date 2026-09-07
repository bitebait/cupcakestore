package bootstrap

import (
	"encoding/json"
	"github.com/Masterminds/sprig/v3"
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/routers"
	"github.com/bitebait/cupcakestore/session"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"time"
)

const (
	faviconPath = "./web/dist/img/favicon.png"
	faviconURL  = "/favicon.ico"
)

func NewApplication() *fiber.App {
	app, err := NewApplicationWithError()
	if err != nil {
		panic(err)
	}
	return app
}

func NewApplicationWithError() (*fiber.App, error) {
	if err := config.Initialize(); err != nil {
		return nil, err
	}
	db, err := database.Open(config.Get())
	if err != nil {
		return nil, err
	}
	database.DB = db
	session.SetupSession()

	fiberApp := createFiberApp()
	registerMiddlewares(fiberApp)
	configureHTTPS(fiberApp)
	serveStaticFiles(fiberApp)
	registerRoutes(fiberApp)
	return fiberApp, nil
}

func createFiberApp() *fiber.App {
	engine := setupTemplateEngine()

	return fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	})
}

func setupTemplateEngine() *html.Engine {
	engine := html.New("./views", ".html")
	engine.AddFuncMap(sprig.FuncMap())
	engine.Reload(config.Get().DevMode)
	return engine
}

func registerMiddlewares(fiberApp *fiber.App) {
	fiberApp.Use(logger.New())
	fiberApp.Use(recover.New())
	fiberApp.Use(csrf.New(csrf.Config{
		CookieHTTPOnly:    true,
		CookieSecure:      !config.Get().DevMode,
		CookieSameSite:    "Lax",
		CookiePath:        "/",
		Expiration:        time.Hour,
		KeyLookup:         "form:_csrf",
		ContextKey:        "CSRFToken",
		Session:           session.Store,
		SessionKey:        "fiber.csrf.token",
		HandlerContextKey: "fiber.csrf.handler",
	}))
	fiberApp.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	fiberApp.Use(favicon.New(favicon.Config{File: faviconPath, URL: faviconURL}))
	fiberApp.Use(middlewares.Message())
}

func serveStaticFiles(fiberApp *fiber.App) {
	fiberApp.Static("/", "./web")
}

func configureHTTPS(fiberApp *fiber.App) {
	if !config.Get().DevMode {
		fiberApp.Use(redirectToHTTPS)
	}
}

func registerRoutes(fiberApp *fiber.App) {
	fiberApp.Use(middlewares.Auth())
	routers.InstallRouters(fiberApp)
}

func redirectToHTTPS(c *fiber.Ctx) error {
	if c.Protocol() == "http" {
		return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
	}
	return c.Next()
}
