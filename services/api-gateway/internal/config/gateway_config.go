package config

import "os"

type GatewayConfig struct {
	Port            string
	IdentityURL     string
	CatalogURL      string
	CartURL         string
	OrdersURL       string
	PaymentURL      string
	ShippingURL     string
	InventoryURL    string
	PricingURL      string
	ReviewsURL      string
	CMSURL          string
	TaxURL          string
	ActivityLogURL  string
	NotificationsURL string
	SearchURL       string
}

func Load() *GatewayConfig {
	return &GatewayConfig{
		Port:             getEnv("GATEWAY_PORT", "8080"),
		IdentityURL:      getEnv("IDENTITY_SERVICE_URL", "http://localhost:8081"),
		CatalogURL:       getEnv("CATALOG_SERVICE_URL", "http://localhost:8082"),
		CartURL:          getEnv("CART_SERVICE_URL", "http://localhost:8083"),
		OrdersURL:        getEnv("ORDERS_SERVICE_URL", "http://localhost:8084"),
		PaymentURL:       getEnv("PAYMENT_SERVICE_URL", "http://localhost:8085"),
		ShippingURL:      getEnv("SHIPPING_SERVICE_URL", "http://localhost:8086"),
		InventoryURL:     getEnv("INVENTORY_SERVICE_URL", "http://localhost:8087"),
		PricingURL:       getEnv("PRICING_SERVICE_URL", "http://localhost:8088"),
		ReviewsURL:       getEnv("REVIEWS_SERVICE_URL", "http://localhost:8089"),
		CMSURL:           getEnv("CMS_SERVICE_URL", "http://localhost:8090"),
		TaxURL:           getEnv("TAX_SERVICE_URL", "http://localhost:8091"),
		ActivityLogURL:   getEnv("ACTIVITYLOG_SERVICE_URL", "http://localhost:8092"),
		NotificationsURL: getEnv("NOTIFICATIONS_SERVICE_URL", "http://localhost:8093"),
		SearchURL:        getEnv("SEARCH_SERVICE_URL", "http://localhost:8094"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
