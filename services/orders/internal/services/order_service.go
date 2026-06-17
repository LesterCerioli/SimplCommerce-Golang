package services

import (
	"context"
	"errors"

	"github.com/simplcommerce-go/services/orders/internal/models"
	"github.com/simplcommerce-go/services/orders/internal/repositories"
)

type CreateOrderRequest struct {
	Items             []CreateOrderItemRequest `json:"items"`
	ShippingAddressID uint                     `json:"shippingAddressId"`
	BillingAddressID  uint                     `json:"billingAddressId"`
	CouponCode        string                   `json:"couponCode"`
	CouponRuleName    string                   `json:"couponRuleName"`
	DiscountAmount    float64                  `json:"discountAmount"`
	OrderNote         string                   `json:"orderNote"`
	ShippingMethod    string                   `json:"shippingMethod"`
	ShippingFeeAmount float64                  `json:"shippingFeeAmount"`
	TaxAmount         float64                  `json:"taxAmount"`
	PaymentMethod     string                   `json:"paymentMethod"`
	PaymentFeeAmount  float64                  `json:"paymentFeeAmount"`
}

type CreateOrderItemRequest struct {
	ProductID    uint    `json:"productId"`
	ProductName  string  `json:"productName"`
	ProductSKU   string  `json:"productSku"`
	ProductPrice float64 `json:"productPrice"`
	Quantity     int     `json:"quantity"`
	DiscountAmount float64 `json:"discountAmount"`
	TaxAmount    float64 `json:"taxAmount"`
	TaxPercent   float64 `json:"taxPercent"`
}

type OrderService interface {
	CreateOrder(ctx context.Context, customerID uint, req CreateOrderRequest) (*models.Order, error)
	UpdateStatus(ctx context.Context, id uint, status string, updatedByID uint) error
	GetOrderByID(ctx context.Context, id uint, customerID uint) (*models.Order, error)
	GetCustomerOrders(ctx context.Context, customerID uint, page, pageSize int) ([]models.Order, int64, error)
	GetAllOrders(ctx context.Context, page, pageSize int) ([]models.Order, int64, error)
}

type orderService struct {
	orderRepo        repositories.OrderRepository
	orderItemRepo    repositories.OrderItemRepository
	orderHistoryRepo repositories.OrderHistoryRepository
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	orderItemRepo repositories.OrderItemRepository,
	orderHistoryRepo repositories.OrderHistoryRepository,
) OrderService {
	return &orderService{
		orderRepo:        orderRepo,
		orderItemRepo:    orderItemRepo,
		orderHistoryRepo: orderHistoryRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, customerID uint, req CreateOrderRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	var items []models.OrderItem
	var subTotal float64

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be positive")
		}
		lineTotal := item.ProductPrice * float64(item.Quantity)
		subTotal += lineTotal

		items = append(items, models.OrderItem{
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
	orderTotal := subTotalWithDiscount + req.ShippingFeeAmount + req.TaxAmount + req.PaymentFeeAmount
	if orderTotal < 0 {
		orderTotal = 0
	}

	order := &models.Order{
		CustomerID:           customerID,
		CreatedByID:          customerID,
		CouponCode:           req.CouponCode,
		CouponRuleName:       req.CouponRuleName,
		DiscountAmount:       req.DiscountAmount,
		SubTotal:             subTotal,
		SubTotalWithDiscount: subTotalWithDiscount,
		ShippingAddressID:    req.ShippingAddressID,
		BillingAddressID:     req.BillingAddressID,
		OrderStatus:          "New",
		OrderNote:            req.OrderNote,
		ShippingMethod:       req.ShippingMethod,
		ShippingFeeAmount:    req.ShippingFeeAmount,
		TaxAmount:            req.TaxAmount,
		OrderTotal:           orderTotal,
		PaymentMethod:        req.PaymentMethod,
		PaymentFeeAmount:     req.PaymentFeeAmount,
		Items:                items,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	history := &models.OrderHistory{
		OrderID:     order.ID,
		NewStatus:   "New",
		CreatedByID: customerID,
	}
	if err := s.orderHistoryRepo.Create(ctx, history); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) UpdateStatus(ctx context.Context, id uint, status string, updatedByID uint) error {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	oldStatus := order.OrderStatus

	if err := s.orderRepo.UpdateStatus(ctx, id, status, updatedByID); err != nil {
		return err
	}

	history := &models.OrderHistory{
		OrderID:     id,
		OldStatus:   oldStatus,
		NewStatus:   status,
		CreatedByID: updatedByID,
	}
	return s.orderHistoryRepo.Create(ctx, history)
}

func (s *orderService) GetOrderByID(ctx context.Context, id uint, customerID uint) (*models.Order, error) {
	order, err := s.orderRepo.FindWithItems(ctx, id)
	if err != nil {
		return nil, err
	}
	if order.CustomerID != customerID {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *orderService) GetCustomerOrders(ctx context.Context, customerID uint, page, pageSize int) ([]models.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return s.orderRepo.FindByCustomerID(ctx, customerID, page, pageSize)
}

func (s *orderService) GetAllOrders(ctx context.Context, page, pageSize int) ([]models.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return s.orderRepo.FindAllPaginated(ctx, page, pageSize)
}
