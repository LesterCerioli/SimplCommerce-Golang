package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"

	"github.com/simplcommerce-go/services/tax/internal/handlers"
	"github.com/simplcommerce-go/services/tax/internal/models"
	"github.com/simplcommerce-go/services/tax/internal/repositories"
	"github.com/simplcommerce-go/services/tax/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(
		&models.TaxClass{},
		&models.TaxRate{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	gormDB := db.GetDB()

	taxClassRepo := repositories.NewTaxClassRepository(gormDB)
	taxRateRepo := repositories.NewTaxRateRepository(gormDB)

	taxSvc := services.NewTaxService(taxClassRepo, taxRateRepo)

	taxClassHandler := handlers.NewTaxClassHandler(taxSvc)
	taxRateHandler := handlers.NewTaxRateHandler(taxSvc)

	app := fiber.New()

	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.LoggerMiddleware())

	api := app.Group("/api")

	taxClassRoutes := api.Group("/tax-classes")
	taxClassRoutes.Get("/", taxClassHandler.GetAll)
	taxClassRoutes.Post("/", taxClassHandler.Create, middleware.AdminRequired())
	taxClassRoutes.Put("/:id", taxClassHandler.Update, middleware.AdminRequired())
	taxClassRoutes.Delete("/:id", taxClassHandler.Delete, middleware.AdminRequired())

	taxRateRoutes := api.Group("/tax-rates")
	taxRateRoutes.Get("/", taxRateHandler.GetAll, middleware.AdminRequired())
	taxRateRoutes.Post("/", taxRateHandler.Create, middleware.AdminRequired())
	taxRateRoutes.Put("/:id", taxRateHandler.Update, middleware.AdminRequired())
	taxRateRoutes.Delete("/:id", taxRateHandler.Delete, middleware.AdminRequired())

	go func() {
		port := cfg.Port
		if port == "" {
			port = "8084"
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
