package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"

	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/handlers"
	"github.com/simplcommerce-go/services/catalog/internal/repositories"
	"github.com/simplcommerce-go/services/catalog/internal/services"
)

func main() {
	cfg := config.Load()

	dbInstance, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbInstance.Close()

	if err := dbInstance.AutoMigrate(
		&catalogmodels.Product{},
		&catalogmodels.Category{},
		&catalogmodels.Brand{},
		&catalogmodels.Media{},
		&catalogmodels.ProductAttribute{},
		&catalogmodels.ProductAttributeGroup{},
		&catalogmodels.ProductAttributeValue{},
		&catalogmodels.ProductOption{},
		&catalogmodels.ProductOptionValue{},
		&catalogmodels.ProductOptionCombination{},
		&catalogmodels.ProductLink{},
		&catalogmodels.ProductMedia{},
		&catalogmodels.ProductTemplate{},
		&catalogmodels.ProductTemplateProductAttribute{},
		&catalogmodels.ProductPriceHistory{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	db := dbInstance.GetDB()

	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	brandRepo := repositories.NewBrandRepository(db)
	attrRepo := repositories.NewProductAttributeRepository(db)
	attrGroupRepo := repositories.NewProductAttributeGroupRepository(db)
	optionRepo := repositories.NewProductOptionRepository(db)
	templateRepo := repositories.NewProductTemplateRepository(db)

	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	brandService := services.NewBrandService(brandRepo)
	attrService := services.NewProductAttributeService(attrRepo)
	attrGroupService := services.NewProductAttributeGroupService(attrGroupRepo)
	optionService := services.NewProductOptionService(optionRepo)
	templateService := services.NewProductTemplateService(templateRepo)

	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	brandHandler := handlers.NewBrandHandler(brandService)
	attrHandler := handlers.NewProductAttributeHandler(attrService)
	attrGroupHandler := handlers.NewProductAttributeGroupHandler(attrGroupService)
	optionHandler := handlers.NewProductOptionHandler(optionService)
	templateHandler := handlers.NewProductTemplateHandler(templateService)

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api")

	products := api.Group("/products")
	products.Get("/", productHandler.Index)
	products.Get("/featured", productHandler.Featured)
	products.Get("/search", productHandler.Search)
	products.Get("/:slug", productHandler.Show)
	products.Post("/", middleware.AdminRequired(), productHandler.Create)
	products.Put("/:id", middleware.AdminRequired(), productHandler.Update)
	products.Delete("/:id", middleware.AdminRequired(), productHandler.Delete)

	categories := api.Group("/categories")
	categories.Get("/", categoryHandler.Index)
	categories.Get("/:slug", categoryHandler.Show)
	categories.Post("/", middleware.AdminRequired(), categoryHandler.Create)
	categories.Put("/:id", middleware.AdminRequired(), categoryHandler.Update)
	categories.Delete("/:id", middleware.AdminRequired(), categoryHandler.Delete)

	brands := api.Group("/brands")
	brands.Get("/", brandHandler.Index)
	brands.Get("/:slug", brandHandler.Show)
	brands.Post("/", middleware.AdminRequired(), brandHandler.Create)
	brands.Put("/:id", middleware.AdminRequired(), brandHandler.Update)
	brands.Delete("/:id", middleware.AdminRequired(), brandHandler.Delete)

	productAttributes := api.Group("/product-attributes")
	productAttributes.Get("/", attrHandler.Index)
	productAttributes.Post("/", middleware.AdminRequired(), attrHandler.Create)

	productAttributeGroups := api.Group("/product-attribute-groups")
	productAttributeGroups.Get("/", attrGroupHandler.Index)
	productAttributeGroups.Post("/", middleware.AdminRequired(), attrGroupHandler.Create)

	productOptions := api.Group("/product-options")
	productOptions.Get("/", optionHandler.Index)
	productOptions.Post("/", middleware.AdminRequired(), optionHandler.Create)

	productTemplates := api.Group("/product-templates")
	productTemplates.Get("/", templateHandler.Index)
	productTemplates.Post("/", middleware.AdminRequired(), templateHandler.Create)

	log.Printf("Catalog service starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
