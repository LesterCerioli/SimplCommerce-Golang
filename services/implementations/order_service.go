package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type OrderService struct {
	db *sql.DB
}

func NewOrderService(db *sql.DB) *OrderService {
	return &OrderService{db: db}
}

type OrderResponse struct {
	ID                  string              `json:"id"`
	CustomerID          string              `json:"customerId"`
	VendorID            *string             `json:"vendorId,omitempty"`
	CreatedByID         string              `json:"createdById"`
	UpdatedByID         *string             `json:"updatedById,omitempty"`
	CouponCode          string              `json:"couponCode"`
	CouponRuleName      string              `json:"couponRuleName"`
	DiscountAmount      float64             `json:"discountAmount"`
	SubTotal            float64             `json:"subTotal"`
	SubTotalWithDiscount float64            `json:"subTotalWithDiscount"`
	ShippingAddressID   string              `json:"shippingAddressId"`
	BillingAddressID    string              `json:"billingAddressId"`
	OrderStatus         string              `json:"orderStatus"`
	OrderNote           string              `json:"orderNote"`
	ParentID            *string             `json:"parentId,omitempty"`
	IsMasterOrder       bool                `json:"isMasterOrder"`
	ShippingMethod      string              `json:"shippingMethod"`
	ShippingFeeAmount   float64             `json:"shippingFeeAmount"`
	TaxAmount           float64             `json:"taxAmount"`
	OrderTotal          float64             `json:"orderTotal"`
	PaymentMethod       string              `json:"paymentMethod"`
	PaymentFeeAmount    float64             `json:"paymentFeeAmount"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
	Items               []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ID             string  `json:"id"`
	OrderID        string  `json:"orderId"`
	ProductID      string  `json:"productId"`
	ProductName    string  `json:"productName"`
	ProductSKU     string  `json:"productSku"`
	ProductPrice   float64 `json:"productPrice"`
	Quantity       int     `json:"quantity"`
	DiscountAmount float64 `json:"discountAmount"`
	TaxAmount      float64 `json:"taxAmount"`
	TaxPercent     float64 `json:"taxPercent"`
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, req CreateOrderRequest) (*OrderResponse, error) {
	orderTotal := req.SubTotalWithDiscount + req.ShippingFeeAmount + req.TaxAmount + req.PaymentFeeAmount
	if orderTotal < 0 {
		orderTotal = 0
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	var orderID string
	var createdAt, updatedAt time.Time
	err = tx.QueryRowContext(ctx, `
		INSERT INTO orders_orders (customer_id, created_by_id, coupon_code, coupon_rule_name, discount_amount,
			sub_total, sub_total_with_discount, shipping_address_id, billing_address_id, order_status, order_note,
			shipping_method, shipping_fee_amount, tax_amount, order_total, payment_method, payment_fee_amount, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,NOW(),NOW())
		RETURNING id, created_at, updated_at`,
		customerID, customerID, req.CouponCode, req.CouponRuleName, req.DiscountAmount,
		req.SubTotal, req.SubTotalWithDiscount, req.ShippingAddressID, req.BillingAddressID, "New", req.OrderNote,
		req.ShippingMethod, req.ShippingFeeAmount, req.TaxAmount, orderTotal, req.PaymentMethod, req.PaymentFeeAmount,
	).Scan(&orderID, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range req.Items {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO orders_order_items (order_id, product_id, product_name, product_sku, product_price,
				quantity, discount_amount, tax_amount, tax_percent, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW(),NOW())`,
			orderID, item.ProductID, item.ProductName, item.ProductSKU, item.ProductPrice,
			item.Quantity, item.DiscountAmount, item.TaxAmount, item.TaxPercent)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders_order_histories (order_id, new_status, created_by_id, created_at, updated_at)
		VALUES ($1, 'New', $2, NOW(), NOW())`, orderID, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create order history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return s.GetOrderByID(ctx, orderID, customerID)
}

func (s *OrderService) GetOrderByID(ctx context.Context, id, customerID string) (*OrderResponse, error) {
	var vendorID, updatedByID, parentID sql.NullString
	var o OrderResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, customer_id, vendor_id, created_by_id, updated_by_id, coupon_code, coupon_rule_name,
			discount_amount, sub_total, sub_total_with_discount, shipping_address_id, billing_address_id,
			order_status, order_note, parent_id, is_master_order, shipping_method, shipping_fee_amount,
			tax_amount, order_total, payment_method, payment_fee_amount, created_at, updated_at
		FROM orders_orders WHERE id = $1
	`, id).Scan(
		&o.ID, &o.CustomerID, &vendorID, &o.CreatedByID, &updatedByID,
		&o.CouponCode, &o.CouponRuleName, &o.DiscountAmount, &o.SubTotal, &o.SubTotalWithDiscount,
		&o.ShippingAddressID, &o.BillingAddressID, &o.OrderStatus, &o.OrderNote, &parentID,
		&o.IsMasterOrder, &o.ShippingMethod, &o.ShippingFeeAmount, &o.TaxAmount, &o.OrderTotal,
		&o.PaymentMethod, &o.PaymentFeeAmount, &o.CreatedAt, &o.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if o.CustomerID != customerID && customerID != "" {
		return nil, fmt.Errorf("order not found")
	}
	if vendorID.Valid {
		o.VendorID = &vendorID.String
	}
	if updatedByID.Valid {
		o.UpdatedByID = &updatedByID.String
	}
	if parentID.Valid {
		o.ParentID = &parentID.String
	}

	items, err := s.getOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return &o, nil
}

func (s *OrderService) GetCustomerOrders(ctx context.Context, customerID string, page, pageSize int) ([]OrderResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders_orders WHERE customer_id = $1`, customerID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT id FROM orders_orders WHERE customer_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, customerID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []OrderResponse
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		o, err := s.GetOrderByID(ctx, id, customerID)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, *o)
	}
	if orders == nil {
		orders = []OrderResponse{}
	}
	return orders, total, nil
}

func (s *OrderService) GetAllOrders(ctx context.Context, page, pageSize int) ([]OrderResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders_orders`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT id FROM orders_orders ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []OrderResponse
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		o, err := s.GetOrderByID(ctx, id, "")
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, *o)
	}
	if orders == nil {
		orders = []OrderResponse{}
	}
	return orders, total, nil
}

func (s *OrderService) UpdateStatus(ctx context.Context, id string, status string, updatedByID string) error {
	var oldStatus string
	err := s.db.QueryRowContext(ctx, `SELECT order_status FROM orders_orders WHERE id = $1`, id).Scan(&oldStatus)
	if err == sql.ErrNoRows {
		return fmt.Errorf("order not found")
	}
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `UPDATE orders_orders SET order_status = $1, updated_by_id = $2, updated_at = NOW() WHERE id = $3`,
		status, updatedByID, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders_order_histories (order_id, old_status, new_status, created_by_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())`, id, oldStatus, status, updatedByID)
	if err != nil {
		return fmt.Errorf("failed to create history: %w", err)
	}

	return tx.Commit()
}

func (s *OrderService) getOrderItems(ctx context.Context, orderID string) ([]OrderItemResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, order_id, product_id, product_name, product_sku, product_price,
			quantity, discount_amount, tax_amount, tax_percent
		FROM orders_order_items WHERE order_id = $1
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var items []OrderItemResponse
	for rows.Next() {
		var item OrderItemResponse
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName, &item.ProductSKU,
			&item.ProductPrice, &item.Quantity, &item.DiscountAmount, &item.TaxAmount, &item.TaxPercent); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []OrderItemResponse{}
	}
	return items, nil
}

type CreateOrderRequest struct {
	Items               []CreateOrderItemRequest
	ShippingAddressID   string
	BillingAddressID    string
	CouponCode          string
	CouponRuleName      string
	DiscountAmount      float64
	SubTotal            float64
	SubTotalWithDiscount float64
	OrderNote           string
	ShippingMethod      string
	ShippingFeeAmount   float64
	TaxAmount           float64
	PaymentMethod       string
	PaymentFeeAmount    float64
}

type CreateOrderItemRequest struct {
	ProductID      string
	ProductName    string
	ProductSKU     string
	ProductPrice   float64
	Quantity       int
	DiscountAmount float64
	TaxAmount      float64
	TaxPercent     float64
}
