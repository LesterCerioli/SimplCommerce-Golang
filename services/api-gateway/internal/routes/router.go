package routes

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/api-gateway/internal/config"
)

type ProxyRouter struct {
	gatewayConfig *config.GatewayConfig
}

func New(gc *config.GatewayConfig) *ProxyRouter {
	return &ProxyRouter{gatewayConfig: gc}
}

type routeMapping struct {
	prefixes []string
	target   string
}

func (r *ProxyRouter) Setup(app *fiber.App) {
	routes := []routeMapping{
		{
			prefixes: []string{"/api/auth", "/api/users", "/api/roles", "/api/countries", "/api/states", "/api/districts", "/api/addresses", "/api/vendors", "/api/customer-groups"},
			target:   r.gatewayConfig.IdentityURL,
		},
		{
			prefixes: []string{"/api/products", "/api/categories", "/api/brands", "/api/product-attributes", "/api/product-attribute-groups", "/api/product-options", "/api/product-templates"},
			target:   r.gatewayConfig.CatalogURL,
		},
		{
			prefixes: []string{"/api/cart"},
			target:   r.gatewayConfig.CartURL,
		},
		{
			prefixes: []string{"/api/orders", "/api/checkout"},
			target:   r.gatewayConfig.OrdersURL,
		},
		{
			prefixes: []string{"/api/payments", "/api/payment-providers"},
			target:   r.gatewayConfig.PaymentURL,
		},
		{
			prefixes: []string{"/api/shipments", "/api/shipping-providers", "/api/shipping-rates", "/api/price-destinations"},
			target:   r.gatewayConfig.ShippingURL,
		},
		{
			prefixes: []string{"/api/stock", "/api/warehouses", "/api/back-in-stock"},
			target:   r.gatewayConfig.InventoryURL,
		},
		{
			prefixes: []string{"/api/cart-rules", "/api/coupons", "/api/catalog-rules"},
			target:   r.gatewayConfig.PricingURL,
		},
		{
			prefixes: []string{"/api/reviews"},
			target:   r.gatewayConfig.ReviewsURL,
		},
		{
			prefixes: []string{"/api/pages", "/api/menus", "/api/menu-items"},
			target:   r.gatewayConfig.CMSURL,
		},
		{
			prefixes: []string{"/api/tax-classes", "/api/tax-rates"},
			target:   r.gatewayConfig.TaxURL,
		},
		{
			prefixes: []string{"/api/activities"},
			target:   r.gatewayConfig.ActivityLogURL,
		},
		{
			prefixes: []string{"/api/notifications"},
			target:   r.gatewayConfig.NotificationsURL,
		},
		{
			prefixes: []string{"/api/search"},
			target:   r.gatewayConfig.SearchURL,
		},
	}

	for _, route := range routes {
		for _, prefix := range route.prefixes {
			app.All(prefix+"/*", r.proxyHandler(route.target))
			app.All(prefix, r.proxyHandler(route.target))
		}
	}
}

func (r *ProxyRouter) proxyHandler(targetBaseURL string) fiber.Handler {
	return func(c fiber.Ctx) error {
		return forward(targetBaseURL, c)
	}
}

func forward(targetBaseURL string, c fiber.Ctx) error {
	fullURL := targetBaseURL + string(c.Request().URI().Path())

	query := string(c.Request().URI().QueryString())
	if query != "" {
		fullURL += "?" + query
	}

	body := c.Body()

	req, err := http.NewRequest(string(c.Method()), fullURL, bytes.NewReader(body))
	if err != nil {
		return response.Error(502, err.Error())
	}

	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response.Error(502, err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.Error(502, err.Error())
	}

	c.Status(resp.StatusCode)
	for k, v := range resp.Header {
		if len(v) > 0 {
			c.Set(k, v[0])
		}
	}
	return c.Send(respBody)
}

func isMatchingPrefix(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
