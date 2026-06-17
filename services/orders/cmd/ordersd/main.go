package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/auth"
	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"
	"github.com/simplcommerce-go/services/orders/internal/handlers"
	"github.com/simplcommerce-go/services/orders/internal/models"
	"github.com/simplcommerce-go/services/orders/internal/repositories"
	"github.com/simplcommerce-go/services/orders/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(
		&models.Order{},
		&models.OrderItem{},
		&models.OrderAddress{},
		&models.OrderHistory{},
	); err != nil {
		log.Fatal("failed to migrate:", err)
	}

	jwtSvc := auth.NewJWTService(cfg.JWT)

	orderRepo := repositories.NewOrderRepository(db.GetDB())
	orderItemRepo := repositories.NewOrderItemRepository(db.GetDB())
	orderHistoryRepo := repositories.NewOrderHistoryRepository(db.GetDB())

	orderSvc := services.NewOrderService(orderRepo, orderItemRepo, orderHistoryRepo)
	checkoutSvc := services.NewCheckoutService(orderSvc)

	orderHandler := handlers.NewOrderHandler(orderSvc)
	checkoutHandler := handlers.NewCheckoutHandler(checkoutSvc)

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

	orders := app.Group("/api/orders")
	orders.Use(middleware.AuthRequired(jwtSvc))
	orders.Get("/", orderHandler.ListOrders)
	orders.Post("/", orderHandler.CreateOrder)
	orders.Get("/:id", orderHandler.GetOrderByID)
	orders.Put("/:id/status", middleware.AdminRequired(), orderHandler.UpdateStatus)

	checkout := app.Group("/api/checkout")
	checkout.Use(middleware.AuthRequired(jwtSvc))
	checkout.Get("/", checkoutHandler.GetCheckout)
	checkout.Post("/", checkoutHandler.ProcessCheckout)

	log.Fatal(app.Listen(":" + cfg.Port))
}
