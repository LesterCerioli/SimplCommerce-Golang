package initializers

import (
	"log"
	"simplcommerce/domain/entities"
)

func RunMigrations(db *Database) {
	err := db.GormDB.AutoMigrate(
		// Identity
		&entities.User{},
		&entities.Role{},
		&entities.Address{},
		&entities.UserAddress{},
		&entities.Country{},
		&entities.StateOrProvince{},
		&entities.District{},
		&entities.Vendor{},
		&entities.CustomerGroup{},
		&entities.Media{},

		// Catalog
		&entities.Product{},
		&entities.Category{},
		&entities.Brand{},
		&entities.ProductCategory{},
		&entities.ProductAttribute{},
		&entities.ProductAttributeGroup{},
		&entities.ProductAttributeValue{},
		&entities.ProductOption{},
		&entities.ProductOptionValue{},
		&entities.ProductOptionCombination{},
		&entities.ProductLink{},
		&entities.ProductMedia{},
		&entities.ProductTemplate{},
		&entities.ProductTemplateProductAttribute{},
		&entities.ProductPriceHistory{},

		// Cart
		&entities.CartItem{},

		// Orders
		&entities.Order{},
		&entities.OrderItem{},
		&entities.OrderAddress{},
		&entities.OrderHistory{},

		// Payment
		&entities.Payment{},
		&entities.PaymentProvider{},

		// Shipping
		&entities.Shipment{},
		&entities.ShipmentItem{},
		&entities.ShippingProvider{},
		&entities.PriceAndDestination{},

		// Inventory
		&entities.Stock{},
		&entities.StockHistory{},
		&entities.Warehouse{},
		&entities.ProductBackInStockSubscription{},

		// Pricing
		&entities.CartRule{},
		&entities.Coupon{},
		&entities.CartRuleCategory{},
		&entities.CartRuleProduct{},
		&entities.CartRuleCustomerGroup{},
		&entities.CartRuleUsage{},
		&entities.CatalogRule{},
		&entities.CatalogRuleCustomerGroup{},

		// Tax
		&entities.TaxClass{},
		&entities.TaxRate{},

		// Reviews
		&entities.Review{},
		&entities.Reply{},

		// CMS
		&entities.Page{},
		&entities.Menu{},
		&entities.MenuItem{},

		// Activity Log
		&entities.Activity{},
		&entities.ActivityType{},

		// Notifications
		&entities.Notification{},

		// Search
		&entities.SearchQuery{},
	)
	if err != nil {
		log.Printf("Warning: Auto-migration warning: %v", err)
	} else {
		log.Println("Database migration completed")
	}
}
