package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"
	"github.com/simplcommerce-go/pkg/auth"

	"github.com/simplcommerce-go/services/cms/internal/handlers"
	"github.com/simplcommerce-go/services/cms/internal/models"
	"github.com/simplcommerce-go/services/cms/internal/repositories"
	"github.com/simplcommerce-go/services/cms/internal/services"
)

func main() {
	cfg := config.Load()
	cfg.ServiceName = "cms"
	cfg.Database.DBName = "simplcommerce_cms"

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Page{},
		&models.Menu{},
		&models.MenuItem{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	pageRepo := repositories.NewPageRepository(db.GetDB())
	menuRepo := repositories.NewMenuRepository(db.GetDB())
	menuItemRepo := repositories.NewMenuItemRepository(db.GetDB())

	pageService := services.NewPageService(pageRepo)
	menuService := services.NewMenuService(menuRepo, menuItemRepo)

	pageHandler := handlers.NewPageHandler(pageService)
	menuHandler := handlers.NewMenuHandler(menuService)

	jwtService := auth.NewJWTService(cfg.JWT)
	authMiddleware := middleware.AuthRequired(jwtService)
	adminMiddleware := middleware.AdminRequired()

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api")

	api.Get("/pages", pageHandler.List)
	api.Get("/pages/:slug", pageHandler.GetBySlug)
	api.Post("/pages", pageHandler.Create, authMiddleware, adminMiddleware)
	api.Put("/pages/:id", pageHandler.Update, authMiddleware, adminMiddleware)
	api.Delete("/pages/:id", pageHandler.Delete, authMiddleware, adminMiddleware)

	api.Get("/menus", menuHandler.List)
	api.Get("/menus/:id", menuHandler.GetByID)
	api.Post("/menus", menuHandler.Create, authMiddleware, adminMiddleware)
	api.Put("/menus/:id", menuHandler.Update, authMiddleware, adminMiddleware)
	api.Delete("/menus/:id", menuHandler.Delete, authMiddleware, adminMiddleware)

	api.Post("/menus/:id/items", menuHandler.AddItem, authMiddleware, adminMiddleware)
	api.Put("/menu-items/:id", menuHandler.UpdateItem, authMiddleware, adminMiddleware)
	api.Delete("/menu-items/:id", menuHandler.DeleteItem, authMiddleware, adminMiddleware)

	log.Printf("CMS service starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
