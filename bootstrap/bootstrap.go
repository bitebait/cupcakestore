package bootstrap

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/database"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/session"
	"github.com/bitebait/cupcakestore/views"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v2"
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
	registerHealthChecks(fiberApp, db)
	registerMiddlewares(fiberApp)
	configureHTTPS(fiberApp)
	serveStaticFiles(fiberApp)
	registerRoutes(fiberApp, db)
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
		StructValidator:   &structValidator{validate: validator.New()},
	})
}

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(value any) error {
	return v.validate.Struct(value)
}

func setupTemplateEngine() *html.Engine {
	engine := html.New("./views", ".html")
	engine.AddFuncMap(sprig.FuncMap())
	engine.AddFunc("money", views.Money)
	engine.Reload(config.Get().DevMode)
	return engine
}

func registerMiddlewares(fiberApp *fiber.App) {
	fiberApp.Use(logger.New())
	fiberApp.Use(recover.New())
	fiberApp.Use(middlewares.SecurityHeaders())
	fiberApp.Use(session.Middleware)
	fiberApp.Use(csrf.New(csrf.Config{
		CookieHTTPOnly: true,
		CookieSecure:   !config.Get().DevMode,
		CookieSameSite: "Lax",
		CookiePath:     "/",
		IdleTimeout:    time.Hour,
		Extractor:      extractors.FromForm("_csrf"),
		Session:        session.Store,
	}))
	fiberApp.Use(func(c fiber.Ctx) error {
		c.Locals("CSRFToken", csrf.TokenFromContext(c))
		return c.Next()
	})
	fiberApp.Use(favicon.New(favicon.Config{File: faviconPath, URL: faviconURL}))
	fiberApp.Use(middlewares.Message())
}

func serveStaticFiles(fiberApp *fiber.App) {
	// Compress only public assets; pages contain secrets such as CSRF tokens.
	// io/fs keeps the compression cache in memory, including read-only containers.
	fiberApp.Use("/", static.New("", static.Config{FS: os.DirFS("./web"), Compress: true}))
}

func configureHTTPS(fiberApp *fiber.App) {
	if !config.Get().DevMode {
		fiberApp.Use(redirectToHTTPS)
	}
}

func redirectToHTTPS(c fiber.Ctx) error {
	if c.Scheme() == "http" {
		return c.Redirect().Status(fiber.StatusMovedPermanently).To("https://" + c.Host() + c.OriginalURL())
	}
	return c.Next()
}
