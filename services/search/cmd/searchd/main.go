package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"

	"github.com/simplcommerce-go/services/search/internal/handlers"
	"github.com/simplcommerce-go/services/search/internal/models"
	"github.com/simplcommerce-go/services/search/internal/services"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = "search"
	cfg.Database.DBName = "simplcommerce_search"

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Query{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	searchRepo := services.NewSearchRepository(db.GetDB())

	searchService := services.NewSearchService(searchRepo)

	searchHandler := handlers.NewSearchHandler(searchService)

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api")

	api.Post("/search", searchHandler.Search)
	api.Get("/search/popular", searchHandler.Popular)

	log.Printf("Search service starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
