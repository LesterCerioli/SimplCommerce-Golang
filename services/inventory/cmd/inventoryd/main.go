package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"
	"github.com/simplcommerce-go/services/inventory/internal/handlers"
	"github.com/simplcommerce-go/services/inventory/internal/models"
	"github.com/simplcommerce-go/services/inventory/internal/repositories"
	"github.com/simplcommerce-go/services/inventory/internal/services"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = "inventory-service"

	dbInstance, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := dbInstance.AutoMigrate(
		&models.Stock{},
		&models.StockHistory{},
		&models.Warehouse{},
		&models.ProductBackInStockSubscription{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	db := dbInstance.GetDB()

	stockRepo := repositories.NewStockRepository(db)
	stockHistoryRepo := repositories.NewStockHistoryRepository(db)
	warehouseRepo := repositories.NewWarehouseRepository(db)
	subRepo := repositories.NewSubscriptionRepository(db)

	inventoryService := services.NewInventoryService(stockRepo, stockHistoryRepo, db)
	warehouseService := services.NewWarehouseService(warehouseRepo)

	stockHandler := handlers.NewStockHandler(inventoryService, stockRepo, stockHistoryRepo)
	warehouseHandler := handlers.NewWarehouseHandler(warehouseService)
	subHandler := handlers.NewSubscriptionHandler(subRepo)

	app := fiber.New()

	app.Use(middleware.CORSMiddleware())

	api := app.Group("/api")

	stockRoutes := api.Group("/stock")
	stockRoutes.Get("/:productId", stockHandler.GetStock)
	stockRoutes.Put("/:productId", stockHandler.UpdateStock)
	stockRoutes.Get("/:productId/history", stockHandler.GetStockHistory)

	warehouseRoutes := api.Group("/warehouses")
	warehouseRoutes.Get("/", warehouseHandler.List)
	warehouseRoutes.Post("/", warehouseHandler.Create)
	warehouseRoutes.Put("/:id", warehouseHandler.Update)
	warehouseRoutes.Delete("/:id", warehouseHandler.Delete)

	api.Post("/back-in-stock/subscribe", subHandler.Subscribe)

	log.Printf("Inventory service starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
