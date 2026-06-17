package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"
	"github.com/simplcommerce-go/pkg/auth"

	"github.com/simplcommerce-go/services/activitylog/internal/handlers"
	"github.com/simplcommerce-go/services/activitylog/internal/models"
	"github.com/simplcommerce-go/services/activitylog/internal/repositories"
	"github.com/simplcommerce-go/services/activitylog/internal/services"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = "activitylog"
	cfg.Database.DBName = "simplcommerce_activitylog"

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.ActivityType{},
		&models.Activity{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	activityRepo := repositories.NewActivityRepository(db.GetDB())
	activityTypeRepo := repositories.NewActivityTypeRepository(db.GetDB())

	activityService := services.NewActivityService(activityRepo, activityTypeRepo)

	activityHandler := handlers.NewActivityHandler(activityService)

	jwtService := auth.NewJWTService(cfg.JWT)
	authMiddleware := middleware.AuthRequired(jwtService)
	adminMiddleware := middleware.AdminRequired()

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api")

	api.Get("/activities", activityHandler.List, authMiddleware, adminMiddleware)
	api.Get("/activities/most-viewed", activityHandler.MostViewed)
	api.Post("/activities", activityHandler.LogActivity, authMiddleware)

	log.Printf("ActivityLog service starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
