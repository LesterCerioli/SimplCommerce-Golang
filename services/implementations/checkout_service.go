package implementations

import (
	"context"
	"database/sql"
	"fmt"
)

type CheckoutService struct {
	db *sql.DB
}

func NewCheckoutService(db *sql.DB) *CheckoutService {
	return &CheckoutService{db: db}
}

type CheckoutItemRequest struct {
	ProductID      string  `json:"productId"`
	ProductName    string  `json:"productName"`
	ProductSKU     string  `json:"productSku"`
	ProductPrice   float64 `json:"productPrice"`
	Quantity       int     `json:"quantity"`
	DiscountAmount float64 `json:"discountAmount"`
	TaxAmount      float64 `json:"taxAmount"`
	TaxPercent     float64 `json:"taxPercent"`
}

type AddressRequest struct {
	ContactName       string  `json:"contactName"`
	Phone             string  `json:"phone"`
	AddressLine1      string  `json:"addressLine1"`
	AddressLine2      string  `json:"addressLine2"`
	City              string  `json:"city"`
	ZipCode           string  `json:"zipCode"`
	DistrictID        *string `json:"districtId"`
	StateOrProvinceID string  `json:"stateOrProvinceId"`
	CountryID         string  `json:"countryId"`
}

type CheckoutRequest struct {
	Items            []CheckoutItemRequest `json:"items"`
	ShippingAddress  AddressRequest        `json:"shippingAddress"`
	BillingAddress   AddressRequest        `json:"billingAddress"`
	CouponCode       string                `json:"couponCode"`
	ShippingMethod   string                `json:"shippingMethod"`
	PaymentMethod    string                `json:"paymentMethod"`
	OrderNote        string                `json:"orderNote"`
	ShippingFeeAmount float64               `json:"shippingFeeAmount"`
	TaxAmount        float64               `json:"taxAmount"`
	PaymentFeeAmount float64               `json:"paymentFeeAmount"`
	DiscountAmount   float64               `json:"discountAmount"`
}

type CheckoutInfo struct {
	CustomerID      string   `json:"customerId"`
	ShippingMethods []string `json:"shippingMethods"`
	PaymentMethods  []string `json:"paymentMethods"`
}

func (s *CheckoutService) GetCheckout(ctx context.Context, customerID string) (*CheckoutInfo, error) {
	return &CheckoutInfo{
		CustomerID:      customerID,
		ShippingMethods: []string{"Standard", "Express"},
		PaymentMethods:  []string{"Credit Card", "Bank Transfer", "COD"},
	}, nil
}

func (s *CheckoutService) ProcessCheckout(ctx context.Context, customerID string, req CheckoutRequest) (*OrderResponse, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("checkout must have at least one item")
	}

	var subTotal float64
	var orderItems []CreateOrderItemRequest
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("item quantity must be positive")
		}
		lineTotal := item.ProductPrice * float64(item.Quantity)
		subTotal += lineTotal
		orderItems = append(orderItems, CreateOrderItemRequest{
			ProductID:      item.ProductID,
			ProductName:    item.ProductName,
			ProductSKU:     item.ProductSKU,
			ProductPrice:   item.ProductPrice,
			Quantity:       item.Quantity,
			DiscountAmount: item.DiscountAmount,
			TaxAmount:      item.TaxAmount,
			TaxPercent:     item.TaxPercent,
		})
	}

	subTotalWithDiscount := subTotal - req.DiscountAmount

	createReq := CreateOrderRequest{
		Items:               orderItems,
		ShippingAddressID:   "",
		BillingAddressID:    "",
		CouponCode:          req.CouponCode,
		DiscountAmount:      req.DiscountAmount,
		SubTotal:            subTotal,
		SubTotalWithDiscount: subTotalWithDiscount,
		OrderNote:           req.OrderNote,
		ShippingMethod:      req.ShippingMethod,
		ShippingFeeAmount:   req.ShippingFeeAmount,
		TaxAmount:           req.TaxAmount,
		PaymentMethod:       req.PaymentMethod,
		PaymentFeeAmount:    req.PaymentFeeAmount,
	}

	orderSvc := &OrderService{db: s.db}
	return orderSvc.CreateOrder(ctx, customerID, createReq)
}
