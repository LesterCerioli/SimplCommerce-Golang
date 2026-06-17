package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/auth"
	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"

	"github.com/simplcommerce-go/services/shipping/internal/handlers"
	"github.com/simplcommerce-go/services/shipping/internal/models"
	"github.com/simplcommerce-go/services/shipping/internal/repositories"
	"github.com/simplcommerce-go/services/shipping/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(
		&models.Shipment{},
		&models.ShipmentItem{},
		&models.ShippingProvider{},
		&models.PriceAndDestination{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	gormDB := db.GetDB()

	jwtSvc := auth.NewJWTService(cfg.JWT)

	shipmentRepo := repositories.NewShipmentRepository(gormDB)
	shipmentItemRepo := repositories.NewShipmentItemRepository(gormDB)
	providerRepo := repositories.NewShippingProviderRepository(gormDB)
	priceRepo := repositories.NewPriceAndDestinationRepository(gormDB)

	shipmentSvc := services.NewShipmentService(shipmentRepo, shipmentItemRepo, gormDB)
	providerSvc := services.NewShippingProviderService(providerRepo)
	rateSvc := services.NewShippingRateService(priceRepo)
	priceDestSvc := services.NewPriceDestinationService(priceRepo)

	shipmentHandler := handlers.NewShipmentHandler(shipmentSvc)
	providerHandler := handlers.NewShippingProviderHandler(providerSvc)
	rateHandler := handlers.NewShippingRateHandler(rateSvc)
	priceDestHandler := handlers.NewPriceDestinationHandler(priceDestSvc)

	app := fiber.New()

	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.LoggerMiddleware())

	api := app.Group("/api")

	shipmentRoutes := api.Group("/shipments")
	shipmentRoutes.Get("/", shipmentHandler.GetAll, middleware.AdminRequired())
	shipmentRoutes.Post("/", shipmentHandler.Create, middleware.AdminRequired())
	shipmentRoutes.Get("/:id", shipmentHandler.GetByID, middleware.AdminRequired())

	providerRoutes := api.Group("/shipping-providers")
	providerRoutes.Get("/", providerHandler.GetAll)
	providerRoutes.Put("/:id", providerHandler.Update, middleware.AdminRequired())

	api.Get("/shipping-rates", rateHandler.Calculate, middleware.AuthRequired(jwtSvc))

	priceDestRoutes := api.Group("/price-destinations")
	priceDestRoutes.Get("/", priceDestHandler.GetAll, middleware.AdminRequired())
	priceDestRoutes.Post("/", priceDestHandler.Create, middleware.AdminRequired())

	go func() {
		port := cfg.Port
		if port == "" {
			port = "8083"
		}
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
}
