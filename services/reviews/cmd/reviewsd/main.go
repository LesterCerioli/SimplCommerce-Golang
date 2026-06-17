package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"
	"github.com/simplcommerce-go/pkg/auth"

	"github.com/simplcommerce-go/services/reviews/internal/handlers"
	"github.com/simplcommerce-go/services/reviews/internal/models"
	"github.com/simplcommerce-go/services/reviews/internal/repositories"
	"github.com/simplcommerce-go/services/reviews/internal/services"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = "reviews"
	cfg.Database.DBName = "simplcommerce_reviews"

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Review{},
		&models.Reply{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	reviewRepo := repositories.NewReviewRepository(db.GetDB())
	replyRepo := repositories.NewReplyRepository(db.GetDB())

	reviewService := services.NewReviewService(reviewRepo)
	replyService := services.NewReplyService(replyRepo)

	reviewHandler := handlers.NewReviewHandler(reviewService)
	replyHandler := handlers.NewReplyHandler(replyService)

	jwtService := auth.NewJWTService(cfg.JWT)
	authMiddleware := middleware.AuthRequired(jwtService)
	adminMiddleware := middleware.AdminRequired()

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api")

	api.Get("/reviews", reviewHandler.GetProductReviews)
	api.Post("/reviews", reviewHandler.Create, authMiddleware)
	api.Put("/reviews/:id/status", reviewHandler.UpdateStatus, authMiddleware, adminMiddleware)
	api.Delete("/reviews/:id", reviewHandler.Delete, authMiddleware, adminMiddleware)

	api.Post("/reviews/:id/replies", replyHandler.Create, authMiddleware)

	log.Printf("Reviews service starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
