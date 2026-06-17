package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/simplcommerce-go/pkg/auth"
	"github.com/simplcommerce-go/pkg/config"
	"github.com/simplcommerce-go/pkg/database"
	"github.com/simplcommerce-go/pkg/middleware"
	"github.com/simplcommerce-go/services/identity/internal/handlers"
	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
	"github.com/simplcommerce-go/services/identity/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Address{},
		&models.Country{},
		&models.StateOrProvince{},
		&models.District{},
		&models.Vendor{},
		&models.CustomerGroup{},
		&models.Media{},
		&models.UserAddress{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	gormDB := db.GetDB()

	userRepo := repositories.NewUserRepository(gormDB)
	roleRepo := repositories.NewRoleRepository(gormDB)
	addressRepo := repositories.NewAddressRepository(gormDB)
	userAddrRepo := repositories.NewUserAddressRepository(gormDB)
	countryRepo := repositories.NewCountryRepository(gormDB)
	stateRepo := repositories.NewStateOrProvinceRepository(gormDB)
	districtRepo := repositories.NewDistrictRepository(gormDB)
	vendorRepo := repositories.NewVendorRepository(gormDB)
	groupRepo := repositories.NewCustomerGroupRepository(gormDB)
	mediaRepo := repositories.NewMediaRepository(gormDB)

	jwtService := auth.NewJWTService(cfg.JWT)

	authService := services.NewAuthService(userRepo, roleRepo, jwtService)
	userService := services.NewUserService(userRepo)
	roleService := services.NewRoleService(roleRepo)
	addressService := services.NewAddressService(addressRepo, userAddrRepo)
	countryService := services.NewCountryService(countryRepo)
	stateService := services.NewStateOrProvinceService(stateRepo)
	districtService := services.NewDistrictService(districtRepo)
	vendorService := services.NewVendorService(vendorRepo)
	groupService := services.NewCustomerGroupService(groupRepo)
	mediaService := services.NewMediaService(mediaRepo)

	authHandler := handlers.NewAuthHandler(authService, userService)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	addressHandler := handlers.NewAddressHandler(addressService)
	countryHandler := handlers.NewCountryHandler(countryService)
	stateHandler := handlers.NewStateHandler(stateService)
	districtHandler := handlers.NewDistrictHandler(districtService)
	vendorHandler := handlers.NewVendorHandler(vendorService)
	groupHandler := handlers.NewCustomerGroupHandler(groupService)
	mediaHandler := handlers.NewMediaHandler(mediaService)

	app := fiber.New()

	app.Use(cors.New())

	authRequired := middleware.AuthRequired(jwtService)
	adminRequired := middleware.AdminRequired()

	auth := app.Group("/api/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Get("/me", authRequired, authHandler.Me)

	users := app.Group("/api/users", authRequired, adminRequired)
	users.Get("/", userHandler.Index)
	users.Post("/", userHandler.Create)
	users.Get("/:id", userHandler.Show)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	roles := app.Group("/api/roles", authRequired, adminRequired)
	roles.Get("/", roleHandler.Index)
	roles.Post("/", roleHandler.Create)

	countries := app.Group("/api/countries")
	countries.Get("/", countryHandler.Index)
	countries.Get("/:id/states", stateHandler.Index)

	app.Get("/api/states/:id/districts", districtHandler.Index)

	addresses := app.Group("/api/addresses", authRequired)
	addresses.Get("/", addressHandler.Index)
	addresses.Post("/", addressHandler.Create)
	addresses.Put("/:id", addressHandler.Update)
	addresses.Delete("/:id", addressHandler.Delete)

	vendors := app.Group("/api/vendors")
	vendors.Get("/", vendorHandler.Index)
	vendors.Post("/", authRequired, adminRequired, vendorHandler.Create)
	vendors.Put("/:id", authRequired, adminRequired, vendorHandler.Update)

	groups := app.Group("/api/customer-groups", authRequired, adminRequired)
	groups.Get("/", groupHandler.Index)
	groups.Post("/", groupHandler.Create)

	media := app.Group("/api/media", authRequired, adminRequired)
	media.Get("/", mediaHandler.Index)
	media.Get("/:id", mediaHandler.Show)
	media.Post("/", mediaHandler.Create)
	media.Put("/:id", mediaHandler.Update)
	media.Delete("/:id", mediaHandler.Delete)

	log.Printf("Identity service starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
