package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"
	"github.com/simplcommerce-go/services/pricing/internal/handlers"
	"github.com/simplcommerce-go/services/pricing/internal/models"
	"github.com/simplcommerce-go/services/pricing/internal/repositories"
	"github.com/simplcommerce-go/services/pricing/internal/services"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = "pricing-service"

	dbInstance, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := dbInstance.AutoMigrate(
		&models.CartRule{},
		&models.Coupon{},
		&models.CartRuleCategory{},
		&models.CartRuleProduct{},
		&models.CartRuleCustomerGroup{},
		&models.CartRuleUsage{},
		&models.CatalogRule{},
		&models.CatalogRuleCustomerGroup{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	db := dbInstance.GetDB()

	cartRuleRepo := repositories.NewCartRuleRepository(db)
	couponRepo := repositories.NewCouponRepository(db)
	catalogRuleRepo := repositories.NewCatalogRuleRepository(db)
	usageRepo := repositories.NewCartRuleUsageRepository(db)

	pricingService := services.NewPricingService(cartRuleRepo, couponRepo, catalogRuleRepo, usageRepo, db)

	cartRuleHandler := handlers.NewCartRuleHandler(cartRuleRepo)
	couponHandler := handlers.NewCouponHandler(couponRepo, cartRuleRepo, pricingService)
	catalogRuleHandler := handlers.NewCatalogRuleHandler(catalogRuleRepo)

	app := fiber.New()

	app.Use(middleware.CORSMiddleware())

	api := app.Group("/api")

	cartRuleRoutes := api.Group("/cart-rules")
	cartRuleRoutes.Get("/", cartRuleHandler.List)
	cartRuleRoutes.Post("/", cartRuleHandler.Create)
	cartRuleRoutes.Put("/:id", cartRuleHandler.Update)
	cartRuleRoutes.Delete("/:id", cartRuleHandler.Delete)

	cartRuleRoutes.Get("/:id/coupons", couponHandler.ListByCartRule)
	cartRuleRoutes.Post("/:id/coupons", couponHandler.Create)

	api.Post("/coupons/validate", couponHandler.Validate)

	catalogRuleRoutes := api.Group("/catalog-rules")
	catalogRuleRoutes.Get("/", catalogRuleHandler.List)
	catalogRuleRoutes.Post("/", catalogRuleHandler.Create)
	catalogRuleRoutes.Put("/:id", catalogRuleHandler.Update)
	catalogRuleRoutes.Delete("/:id", catalogRuleHandler.Delete)

	log.Printf("Pricing service starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
