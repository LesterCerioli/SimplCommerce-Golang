package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"
	"github.com/simplcommerce-go/pkg/auth"

	"github.com/simplcommerce-go/services/notifications/internal/handlers"
	"github.com/simplcommerce-go/services/notifications/internal/models"
	"github.com/simplcommerce-go/services/notifications/internal/repositories"
	"github.com/simplcommerce-go/services/notifications/internal/services"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = "notifications"
	cfg.Database.DBName = "simplcommerce_notifications"

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Notification{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	notificationRepo := repositories.NewNotificationRepository(db.GetDB())

	notificationService := services.NewNotificationService(notificationRepo)

	notificationHandler := handlers.NewNotificationHandler(notificationService)

	jwtService := auth.NewJWTService(cfg.JWT)
	authMiddleware := middleware.AuthRequired(jwtService)

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api")

	api.Get("/notifications", notificationHandler.GetUserNotifications, authMiddleware)
	api.Put("/notifications/:id/read", notificationHandler.MarkAsRead, authMiddleware)
	api.Get("/notifications/unread-count", notificationHandler.UnreadCount, authMiddleware)

	log.Printf("Notifications service starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
