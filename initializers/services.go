package initializers

import (
	"os"
	"strconv"

	"simplcommerce/services/implementations"
)

type Services struct {
	SQLDB *Database

	AuthService          *implementations.AuthService
	UserService          *implementations.UserService
	RoleService          *implementations.RoleService
	AddressService       *implementations.AddressService
	CountryService       *implementations.CountryService
	VendorService        *implementations.VendorService
	CustomerGroupService *implementations.CustomerGroupService
	MediaService         *implementations.MediaService

	ProductService              *implementations.ProductService
	CategoryService             *implementations.CategoryService
	BrandService                *implementations.BrandService
	ProductAttributeService     *implementations.ProductAttributeService
	ProductAttributeGroupService *implementations.ProductAttributeGroupService
	ProductOptionService        *implementations.ProductOptionService
	ProductTemplateService      *implementations.ProductTemplateService

	CartService     *implementations.CartService
	OrderService    *implementations.OrderService
	CheckoutService *implementations.CheckoutService

	PaymentService *implementations.PaymentService

	ShipmentService *implementations.ShipmentService

	StockService     *implementations.StockService
	WarehouseService *implementations.WarehouseService

	ReviewService *implementations.ReviewService

	PageService *implementations.PageService
	MenuService *implementations.MenuService

	TaxService *implementations.TaxService

	PricingService *implementations.PricingService

	ActivityService     *implementations.ActivityService
	NotificationService *implementations.NotificationService

	SearchService *implementations.SearchService
}

func InitServices(db *Database) *Services {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "simplcommerce-secret-key-change-in-production"
	}
	jwtExpiry := 24
	if v := os.Getenv("JWT_EXPIRY_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			jwtExpiry = n
		}
	}

	return &Services{
		SQLDB:                     db,
		AuthService:               implementations.NewAuthService(db.SQLDB, jwtSecret, jwtExpiry),
		UserService:               implementations.NewUserService(db.SQLDB),
		RoleService:               implementations.NewRoleService(db.SQLDB),
		AddressService:            implementations.NewAddressService(db.SQLDB),
		CountryService:            implementations.NewCountryService(db.SQLDB),
		VendorService:             implementations.NewVendorService(db.SQLDB),
		CustomerGroupService:      implementations.NewCustomerGroupService(db.SQLDB),
		MediaService:              implementations.NewMediaService(db.SQLDB),
		ProductService:            implementations.NewProductService(db.SQLDB),
		CategoryService:           implementations.NewCategoryService(db.SQLDB),
		BrandService:              implementations.NewBrandService(db.SQLDB),
		ProductAttributeService:   implementations.NewProductAttributeService(db.SQLDB),
		ProductAttributeGroupService: implementations.NewProductAttributeGroupService(db.SQLDB),
		ProductOptionService:      implementations.NewProductOptionService(db.SQLDB),
		ProductTemplateService:    implementations.NewProductTemplateService(db.SQLDB),
		CartService:               implementations.NewCartService(db.SQLDB),
		OrderService:              implementations.NewOrderService(db.SQLDB),
		CheckoutService:           implementations.NewCheckoutService(db.SQLDB),
		PaymentService:            implementations.NewPaymentService(db.SQLDB),
		ShipmentService:           implementations.NewShipmentService(db.SQLDB),
		StockService:              implementations.NewStockService(db.SQLDB),
		WarehouseService:          implementations.NewWarehouseService(db.SQLDB),
		ReviewService:             implementations.NewReviewService(db.SQLDB),
		PageService:               implementations.NewPageService(db.SQLDB),
		MenuService:               implementations.NewMenuService(db.SQLDB),
		TaxService:                implementations.NewTaxService(db.SQLDB),
		PricingService:            implementations.NewPricingService(db.SQLDB),
		ActivityService:           implementations.NewActivityService(db.SQLDB),
		NotificationService:       implementations.NewNotificationService(db.SQLDB),
		SearchService:             implementations.NewSearchService(db.SQLDB),
	}
}
