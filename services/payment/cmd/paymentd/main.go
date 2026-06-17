package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"

	"github.com/simplcommerce-go/services/payment/internal/handlers"
	"github.com/simplcommerce-go/services/payment/internal/models"
	"github.com/simplcommerce-go/services/payment/internal/repositories"
	"github.com/simplcommerce-go/services/payment/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(
		&models.Payment{},
		&models.PaymentProvider{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	gormDB := db.GetDB()

	paymentRepo := repositories.NewPaymentRepository(gormDB)
	providerRepo := repositories.NewPaymentProviderRepository(gormDB)

	paymentSvc := services.NewPaymentService(paymentRepo)
	providerSvc := services.NewPaymentProviderService(providerRepo)

	paymentHandler := handlers.NewPaymentHandler(paymentSvc)
	providerHandler := handlers.NewPaymentProviderHandler(providerSvc)

	app := fiber.New()

	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.LoggerMiddleware())

	api := app.Group("/api")

	paymentRoutes := api.Group("/payments")
	paymentRoutes.Get("/", paymentHandler.GetAll, middleware.AdminRequired())
	paymentRoutes.Get("/:id", paymentHandler.GetByID, middleware.AdminRequired())
	paymentRoutes.Post("/", paymentHandler.Create, middleware.AdminRequired())

	providerRoutes := api.Group("/payment-providers")
	providerRoutes.Get("/", providerHandler.GetAll)
	providerRoutes.Put("/:id", providerHandler.Update, middleware.AdminRequired())

	go func() {
		port := cfg.Port
		if port == "" {
			port = "8082"
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
