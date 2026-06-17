package services

import (
	"context"
	"errors"

	"github.com/simplcommerce-go/services/orders/internal/models"
)

type CheckoutItemRequest struct {
	ProductID      uint    `json:"productId"`
	ProductName    string  `json:"productName"`
	ProductSKU     string  `json:"productSku"`
	ProductPrice   float64 `json:"productPrice"`
	Quantity       int     `json:"quantity"`
	DiscountAmount float64 `json:"discountAmount"`
	TaxAmount      float64 `json:"taxAmount"`
	TaxPercent     float64 `json:"taxPercent"`
}

type AddressRequest struct {
	ContactName       string `json:"contactName"`
	Phone             string `json:"phone"`
	AddressLine1      string `json:"addressLine1"`
	AddressLine2      string `json:"addressLine2"`
	City              string `json:"city"`
	ZipCode           string `json:"zipCode"`
	DistrictID        *uint  `json:"districtId"`
	StateOrProvinceID uint   `json:"stateOrProvinceId"`
	CountryID         string `json:"countryId"`
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
	CustomerID      uint   `json:"customerId"`
	DefaultAddress  *AddressRequest `json:"defaultAddress,omitempty"`
	ShippingMethods []string `json:"shippingMethods,omitempty"`
	PaymentMethods  []string `json:"paymentMethods,omitempty"`
}

type CheckoutService interface {
	GetCheckout(ctx context.Context, customerID uint) (*CheckoutInfo, error)
	ProcessCheckout(ctx context.Context, customerID uint, req CheckoutRequest) (*models.Order, error)
}

type checkoutService struct {
	orderService OrderService
	db           interface {
		Create(ctx context.Context, entity interface{}) error
	}
}

func NewCheckoutService(orderService OrderService) CheckoutService {
	return &checkoutService{
		orderService: orderService,
	}
}

func (s *checkoutService) GetCheckout(ctx context.Context, customerID uint) (*CheckoutInfo, error) {
	return &CheckoutInfo{
		CustomerID:      customerID,
		ShippingMethods: []string{"Standard", "Express"},
		PaymentMethods:  []string{"Credit Card", "Bank Transfer", "COD"},
	}, nil
}

func (s *checkoutService) ProcessCheckout(ctx context.Context, customerID uint, req CheckoutRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("checkout must have at least one item")
	}

	if req.ShippingAddress.ContactName == "" || req.ShippingAddress.AddressLine1 == "" {
		return nil, errors.New("shipping address is required")
	}

	var orderItems []CreateOrderItemRequest
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be positive")
		}
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

	createReq := CreateOrderRequest{
		Items:             orderItems,
		ShippingAddressID: 0,
		BillingAddressID:  0,
		CouponCode:        req.CouponCode,
		DiscountAmount:    req.DiscountAmount,
		OrderNote:         req.OrderNote,
		ShippingMethod:    req.ShippingMethod,
		ShippingFeeAmount: req.ShippingFeeAmount,
		TaxAmount:         req.TaxAmount,
		PaymentMethod:     req.PaymentMethod,
		PaymentFeeAmount:  req.PaymentFeeAmount,
	}

	order, err := s.orderService.CreateOrder(ctx, customerID, createReq)
	if err != nil {
		return nil, err
	}

	return order, nil
}
