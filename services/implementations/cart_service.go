package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CartService struct {
	db *sql.DB
}

func NewCartService(db *sql.DB) *CartService {
	return &CartService{db: db}
}

type CartItemResponse struct {
	ID         string    `json:"id"`
	ProductID  string    `json:"productId"`
	CustomerID string    `json:"customerId"`
	Quantity   int       `json:"quantity"`
	VendorID   *string   `json:"vendorId,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (s *CartService) GetCart(ctx context.Context, customerID string) ([]CartItemResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, product_id, customer_id, quantity, vendor_id, created_at, updated_at
		FROM cart_cart_items WHERE customer_id = $1
	`, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	defer rows.Close()

	var items []CartItemResponse
	for rows.Next() {
		var item CartItemResponse
		var vendorID sql.NullString
		if err := rows.Scan(&item.ID, &item.ProductID, &item.CustomerID, &item.Quantity, &vendorID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}
		if vendorID.Valid {
			item.VendorID = &vendorID.String
		}
		items = append(items, item)
	}
	if items == nil {
		items = []CartItemResponse{}
	}
	return items, nil
}

func (s *CartService) AddItem(ctx context.Context, customerID, productID string, quantity int) error {
	existing, err := s.findItemByCustomerAndProduct(ctx, customerID, productID)
	if err == nil && existing != nil {
		_, err = s.db.ExecContext(ctx, `
			UPDATE cart_cart_items SET quantity = quantity + $1, updated_at = NOW()
			WHERE id = $2 AND customer_id = $3
		`, quantity, existing.ID, customerID)
		if err != nil {
			return fmt.Errorf("failed to update cart item: %w", err)
		}
		return nil
	}

	var maxDisplayOrder int
	s.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(display_order),0) FROM cart_cart_items WHERE customer_id = $1`, customerID).Scan(&maxDisplayOrder)

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO cart_cart_items (product_id, customer_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`, productID, customerID, quantity)
	if err != nil {
		return fmt.Errorf("failed to add cart item: %w", err)
	}
	return nil
}

func (s *CartService) UpdateQuantity(ctx context.Context, customerID, itemID string, quantity int) error {
	if quantity <= 0 {
		return s.RemoveItem(ctx, customerID, itemID)
	}
	result, err := s.db.ExecContext(ctx, `
		UPDATE cart_cart_items SET quantity = $1, updated_at = NOW()
		WHERE id = $2 AND customer_id = $3
	`, quantity, itemID, customerID)
	if err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}
	return nil
}

func (s *CartService) RemoveItem(ctx context.Context, customerID, itemID string) error {
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM cart_cart_items WHERE id = $1 AND customer_id = $2
	`, itemID, customerID)
	if err != nil {
		return fmt.Errorf("failed to remove cart item: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}
	return nil
}

func (s *CartService) ClearCart(ctx context.Context, customerID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM cart_cart_items WHERE customer_id = $1`, customerID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}

func (s *CartService) findItemByCustomerAndProduct(ctx context.Context, customerID, productID string) (*CartItemResponse, error) {
	var item CartItemResponse
	var vendorID sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT id, product_id, customer_id, quantity, vendor_id, created_at, updated_at
		FROM cart_cart_items WHERE customer_id = $1 AND product_id = $2
	`, customerID, productID).Scan(&item.ID, &item.ProductID, &item.CustomerID, &item.Quantity, &vendorID, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if vendorID.Valid {
		item.VendorID = &vendorID.String
	}
	return &item, nil
}
