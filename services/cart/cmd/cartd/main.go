package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/auth"
	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"
	"github.com/simplcommerce-go/services/cart/internal/handlers"
	"github.com/simplcommerce-go/services/cart/internal/models"
	"github.com/simplcommerce-go/services/cart/internal/repositories"
	"github.com/simplcommerce-go/services/cart/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(&models.CartItem{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}

	jwtSvc := auth.NewJWTService(cfg.JWT)

	cartRepo := repositories.NewCartItemRepository(db.GetDB())
	cartSvc := services.NewCartService(cartRepo)
	cartHandler := handlers.NewCartHandler(cartSvc)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	app.Use(middleware.CORSMiddleware())

	api := app.Group("/api/cart")
	api.Use(middleware.AuthRequired(jwtSvc))
	api.Get("/", cartHandler.GetCart)
	api.Delete("/", cartHandler.ClearCart)
	api.Post("/items", cartHandler.AddItem)
	api.Put("/items/:id", cartHandler.UpdateQuantity)
	api.Delete("/items/:id", cartHandler.RemoveItem)

	log.Fatal(app.Listen(":" + cfg.Port))
}
