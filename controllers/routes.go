package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/initializers"
)

func SetupRoutes(app *fiber.App, svcs *initializers.Services) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	authCtrl := NewAuthController(svcs.AuthService)
	auth.Post("/register", authCtrl.Register)
	auth.Post("/login", authCtrl.Login)
	auth.Post("/refresh", authCtrl.Refresh)

	protected := api.Group("", AuthRequired(svcs.AuthService))

	protected.Get("/auth/me", authCtrl.Me)

	userCtrl := NewUserController(svcs.UserService)
	protected.Get("/users", AdminRequired(), userCtrl.List)
	protected.Get("/users/:id", AdminRequired(), userCtrl.Get)
	protected.Put("/users/:id", AdminRequired(), userCtrl.Update)
	protected.Delete("/users/:id", AdminRequired(), userCtrl.Delete)

	roleCtrl := NewRoleController(svcs.RoleService)
	protected.Get("/roles", AdminRequired(), roleCtrl.List)
	protected.Post("/roles", AdminRequired(), roleCtrl.Create)
	protected.Delete("/roles/:id", AdminRequired(), roleCtrl.Delete)

	addrCtrl := NewAddressController(svcs.AddressService)
	protected.Get("/addresses", addrCtrl.List)
	protected.Post("/addresses", addrCtrl.Create)
	protected.Put("/addresses/:id", addrCtrl.Update)
	protected.Delete("/addresses/:id", addrCtrl.Delete)

	countryCtrl := NewCountryController(svcs.CountryService)
	api.Get("/countries", countryCtrl.List)
	api.Get("/countries/:id/states", countryCtrl.ListStates)
	api.Get("/states/:id/districts", countryCtrl.ListDistricts)

	vendorCtrl := NewVendorController(svcs.VendorService)
	api.Get("/vendors", vendorCtrl.List)
	protected.Post("/vendors", AdminRequired(), vendorCtrl.Create)
	protected.Put("/vendors/:id", AdminRequired(), vendorCtrl.Update)

	cgCtrl := NewCustomerGroupController(svcs.CustomerGroupService)
	protected.Get("/customer-groups", AdminRequired(), cgCtrl.List)
	protected.Post("/customer-groups", AdminRequired(), cgCtrl.Create)
	protected.Put("/customer-groups/:id", AdminRequired(), cgCtrl.Update)
	protected.Delete("/customer-groups/:id", AdminRequired(), cgCtrl.Delete)

	mediaCtrl := NewMediaController(svcs.MediaService)
	protected.Get("/media", AdminRequired(), mediaCtrl.List)
	protected.Post("/media/upload", AdminRequired(), mediaCtrl.Upload)
	protected.Delete("/media/:id", AdminRequired(), mediaCtrl.Delete)

	productCtrl := NewProductController(svcs.ProductService, svcs.CategoryService, svcs.BrandService)
	api.Get("/products/featured", productCtrl.Featured)
	api.Get("/products/search", productCtrl.Search)
	api.Get("/products", productCtrl.Index)
	api.Get("/products/:slug", productCtrl.Show)
	protected.Post("/products", AdminRequired(), productCtrl.Create)
	protected.Put("/products/:id", AdminRequired(), productCtrl.Update)
	protected.Delete("/products/:id", AdminRequired(), productCtrl.Delete)

	categoryCtrl := NewCategoryController(svcs.CategoryService)
	api.Get("/categories", categoryCtrl.Index)
	api.Get("/categories/:slug", categoryCtrl.Show)
	protected.Post("/categories", AdminRequired(), categoryCtrl.Create)
	protected.Put("/categories/:id", AdminRequired(), categoryCtrl.Update)
	protected.Delete("/categories/:id", AdminRequired(), categoryCtrl.Delete)

	brandCtrl := NewBrandController(svcs.BrandService)
	api.Get("/brands", brandCtrl.Index)
	api.Get("/brands/:slug", brandCtrl.Show)
	protected.Post("/brands", AdminRequired(), brandCtrl.Create)
	protected.Put("/brands/:id", AdminRequired(), brandCtrl.Update)
	protected.Delete("/brands/:id", AdminRequired(), brandCtrl.Delete)

	attrCtrl := NewProductAttributeController(svcs.ProductAttributeService)
	api.Get("/product-attributes", attrCtrl.Index)
	protected.Post("/product-attributes", AdminRequired(), attrCtrl.Create)
	protected.Delete("/product-attributes/:id", AdminRequired(), attrCtrl.Delete)

	attrGroupCtrl := NewProductAttributeGroupController(svcs.ProductAttributeGroupService)
	api.Get("/product-attribute-groups", attrGroupCtrl.Index)
	protected.Post("/product-attribute-groups", AdminRequired(), attrGroupCtrl.Create)
	protected.Put("/product-attribute-groups/:id", AdminRequired(), attrGroupCtrl.Update)

	optionCtrl := NewProductOptionController(svcs.ProductOptionService)
	api.Get("/product-options", optionCtrl.Index)
	protected.Post("/product-options", AdminRequired(), optionCtrl.Create)
	protected.Put("/product-options/:id", AdminRequired(), optionCtrl.Update)
	protected.Delete("/product-options/:id", AdminRequired(), optionCtrl.Delete)

	templateCtrl := NewProductTemplateController(svcs.ProductTemplateService)
	api.Get("/product-templates", templateCtrl.Index)
	protected.Post("/product-templates", AdminRequired(), templateCtrl.Create)

	cartCtrl := NewCartController(svcs)
	protected.Get("/cart", cartCtrl.GetCart)
	protected.Post("/cart/items", cartCtrl.AddItem)
	protected.Put("/cart/items/:id", cartCtrl.UpdateQuantity)
	protected.Delete("/cart/items/:id", cartCtrl.RemoveItem)
	protected.Delete("/cart", cartCtrl.ClearCart)

	orderCtrl := NewOrderController(svcs)
	protected.Get("/orders", orderCtrl.ListOrders)
	api.Get("/orders/:id", orderCtrl.GetOrder)
	protected.Post("/orders", orderCtrl.CreateOrder)
	protected.Put("/orders/:id/status", orderCtrl.UpdateStatus)

	checkoutCtrl := NewCheckoutController(svcs)
	api.Get("/checkout/:cartId", checkoutCtrl.GetCheckout)
	protected.Post("/checkout/:cartId", checkoutCtrl.ProcessCheckout)

	paymentCtrl := NewPaymentController(svcs)
	protected.Get("/payments", paymentCtrl.ListPayments)
	protected.Get("/payments/:id", paymentCtrl.GetPayment)
	protected.Post("/payments", paymentCtrl.CreatePayment)
	protected.Get("/payment-providers", paymentCtrl.ListProviders)

	shipmentCtrl := NewShipmentController(svcs)
	protected.Get("/shipments", shipmentCtrl.ListShipments)
	protected.Post("/shipments", AdminRequired(), shipmentCtrl.CreateShipment)
	protected.Get("/shipments/:id", shipmentCtrl.GetShipment)

	stockCtrl := NewStockController(svcs)
	protected.Get("/stock/:productId", stockCtrl.GetStock)
	protected.Put("/stock/:productId", stockCtrl.UpdateStock)
	protected.Get("/stock/:productId/history", stockCtrl.GetStockHistory)

	warehouseCtrl := NewWarehouseController(svcs)
	protected.Get("/warehouses", warehouseCtrl.List)
	protected.Post("/warehouses", AdminRequired(), warehouseCtrl.Create)
	protected.Put("/warehouses/:id", AdminRequired(), warehouseCtrl.Update)
	protected.Delete("/warehouses/:id", AdminRequired(), warehouseCtrl.Delete)

	reviewCtrl := NewReviewController(svcs)
	api.Get("/reviews", reviewCtrl.GetProductReviews)
	protected.Post("/reviews", reviewCtrl.CreateReview)
	protected.Put("/reviews/:id/status", AdminRequired(), reviewCtrl.UpdateStatus)
	protected.Delete("/reviews/:id", AdminRequired(), reviewCtrl.DeleteReview)

	pageCtrl := NewPageController(svcs)
	api.Get("/pages", pageCtrl.ListPages)
	api.Get("/pages/:slug", pageCtrl.GetPage)
	protected.Post("/pages", AdminRequired(), pageCtrl.CreatePage)
	protected.Put("/pages/:id", AdminRequired(), pageCtrl.UpdatePage)
	protected.Delete("/pages/:id", AdminRequired(), pageCtrl.DeletePage)

	menuCtrl := NewMenuController(svcs)
	api.Get("/menus", menuCtrl.ListMenus)
	api.Get("/menus/:id", menuCtrl.GetMenu)
	protected.Post("/menus", AdminRequired(), menuCtrl.CreateMenu)
	protected.Put("/menus/:id", AdminRequired(), menuCtrl.UpdateMenu)
	protected.Delete("/menus/:id", AdminRequired(), menuCtrl.DeleteMenu)
	protected.Post("/menus/:id/items", AdminRequired(), menuCtrl.AddMenuItem)
	protected.Put("/menu-items/:id", AdminRequired(), menuCtrl.UpdateMenuItem)
	protected.Delete("/menu-items/:id", AdminRequired(), menuCtrl.DeleteMenuItem)

	taxCtrl := NewTaxController(svcs)
	protected.Get("/tax-classes", AdminRequired(), taxCtrl.ListTaxClasses)
	protected.Post("/tax-classes", AdminRequired(), taxCtrl.CreateTaxClass)
	protected.Put("/tax-classes/:id", AdminRequired(), taxCtrl.UpdateTaxClass)
	protected.Delete("/tax-classes/:id", AdminRequired(), taxCtrl.DeleteTaxClass)
	protected.Get("/tax-rates", AdminRequired(), taxCtrl.ListTaxRates)
	protected.Post("/tax-rates", AdminRequired(), taxCtrl.CreateTaxRate)
	protected.Put("/tax-rates/:id", AdminRequired(), taxCtrl.UpdateTaxRate)
	protected.Delete("/tax-rates/:id", AdminRequired(), taxCtrl.DeleteTaxRate)

	pricingCtrl := NewPricingController(svcs)
	protected.Get("/cart-rules", AdminRequired(), pricingCtrl.ListCartRules)
	protected.Get("/cart-rules/:id", AdminRequired(), pricingCtrl.GetCartRule)
	protected.Post("/cart-rules", AdminRequired(), pricingCtrl.CreateCartRule)
	protected.Put("/cart-rules/:id", AdminRequired(), pricingCtrl.UpdateCartRule)
	protected.Delete("/cart-rules/:id", AdminRequired(), pricingCtrl.DeleteCartRule)
	protected.Post("/coupons/validate", pricingCtrl.ValidateCoupon)
	protected.Get("/catalog-rules", AdminRequired(), pricingCtrl.ListCatalogRules)
	protected.Post("/catalog-rules", AdminRequired(), pricingCtrl.CreateCatalogRule)
	protected.Put("/catalog-rules/:id", AdminRequired(), pricingCtrl.UpdateCatalogRule)
	protected.Delete("/catalog-rules/:id", AdminRequired(), pricingCtrl.DeleteCatalogRule)

	activityCtrl := NewActivityController(svcs)
	protected.Get("/activities", AdminRequired(), activityCtrl.ListActivities)

	notificationCtrl := NewNotificationController(svcs)
	protected.Get("/notifications", notificationCtrl.ListNotifications)
	protected.Put("/notifications/:id/read", notificationCtrl.MarkAsRead)
	protected.Get("/notifications/unread-count", notificationCtrl.UnreadCount)

	searchCtrl := NewSearchController(svcs)
	protected.Post("/search/log", searchCtrl.LogSearch)
	protected.Get("/search/popular", searchCtrl.PopularSearches)
}
