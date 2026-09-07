package bootstrap

import (
	"github.com/bitebait/cupcakestore/config"
	"github.com/bitebait/cupcakestore/controllers"
	"github.com/bitebait/cupcakestore/middlewares"
	"github.com/bitebait/cupcakestore/repositories"
	"github.com/bitebait/cupcakestore/routers"
	"github.com/bitebait/cupcakestore/services"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Assemble shared dependencies once. Route modules only map HTTP endpoints to handlers.
func registerRoutes(app *fiber.App, db *gorm.DB) {
	profileRepository := repositories.NewProfileRepository(db)
	profileService := services.NewProfileService(profileRepository)
	userService := services.NewUserService(repositories.NewUserRepository(db))
	productService := services.NewProductService(repositories.NewProductRepository(db))
	storeConfigService := services.NewStoreConfigService(repositories.NewStoreConfigRepository(db))
	stockService := services.NewStockService(repositories.NewStockRepository(db))
	cartItemService := services.NewShoppingCartItemService(repositories.NewShoppingCartItemRepository(db))
	cartService := services.NewShoppingCartService(repositories.NewShoppingCartRepository(db), cartItemService)
	orderService := services.NewOrderService(repositories.NewOrderRepository(db), storeConfigService)
	dashboardService := services.NewDashboardService(repositories.NewDashboardRepository(db))
	authService := services.NewAuthService(userService, profileService)

	app.Use(middlewares.Auth(profileRepository.FindByUserId))
	routers.RegisterAuthRoutes(app, controllers.NewAuthController(authService))
	routers.RegisterUserRoutes(app, controllers.NewUserController(userService))
	routers.RegisterProfileRoutes(app, controllers.NewProfileController(profileService))
	routers.RegisterProductRoutes(app, controllers.NewProductController(productService))
	routers.RegisterStockRoutes(app, controllers.NewStockController(stockService))
	routers.RegisterStoreConfigRoutes(app, controllers.NewStoreConfigController(storeConfigService))
	routers.RegisterStoreRoutes(app, controllers.NewStoreController(productService))
	routers.RegisterHomeRoutes(app, config.Get().RedirectAfterLogin)
	routers.RegisterShoppingCartRoutes(app, controllers.NewShoppingCartController(cartService))
	routers.RegisterOrderRoutes(app, controllers.NewOrderController(orderService, storeConfigService))
	routers.RegisterDashboardRoutes(app, controllers.NewDashboardController(dashboardService))
}
