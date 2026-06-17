package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ShipmentService struct {
	db *sql.DB
}

func NewShipmentService(db *sql.DB) *ShipmentService {
	return &ShipmentService{db: db}
}

type ShipmentResponse struct {
	ID             string                `json:"id"`
	OrderID        string                `json:"orderId"`
	TrackingNumber string                `json:"trackingNumber"`
	WarehouseID    string                `json:"warehouseId"`
	VendorID       *string               `json:"vendorId,omitempty"`
	CreatedByID    string                `json:"createdById"`
	CreatedAt      time.Time             `json:"createdAt"`
	UpdatedAt      time.Time             `json:"updatedAt"`
	Items          []ShipmentItemResponse `json:"items"`
}

type ShipmentItemResponse struct {
	ID          string `json:"id"`
	ShipmentID  string `json:"shipmentId"`
	OrderItemID string `json:"orderItemId"`
	ProductID   string `json:"productId"`
	Quantity    int    `json:"quantity"`
}

type CreateShipmentInput struct {
	OrderID        string              `json:"orderId"`
	TrackingNumber string              `json:"trackingNumber"`
	WarehouseID    string              `json:"warehouseId"`
	Items          []ShipmentItemInput `json:"items"`
}

type ShipmentItemInput struct {
	OrderItemID string `json:"orderItemId"`
	ProductID   string `json:"productId"`
	Quantity    int    `json:"quantity"`
}

func (s *ShipmentService) CreateShipment(ctx context.Context, input CreateShipmentInput, createdByID string) (*ShipmentResponse, error) {
	if input.OrderID == "" {
		return nil, fmt.Errorf("order id is required")
	}
	if input.WarehouseID == "" {
		return nil, fmt.Errorf("warehouse id is required")
	}
	if len(input.Items) == 0 {
		return nil, fmt.Errorf("at least one shipment item is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	var shipmentID string
	var createdAt, updatedAt time.Time
	err = tx.QueryRowContext(ctx, `
		INSERT INTO shipping_shipments (order_id, tracking_number, warehouse_id, created_by_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,NOW(),NOW()) RETURNING id, created_at, updated_at
	`, input.OrderID, input.TrackingNumber, input.WarehouseID, createdByID).Scan(&shipmentID, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	for _, item := range input.Items {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO shipping_shipment_items (shipment_id, order_item_id, product_id, quantity, created_at, updated_at)
			VALUES ($1,$2,$3,$4,NOW(),NOW())
		`, shipmentID, item.OrderItemID, item.ProductID, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to create shipment item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return s.GetShipment(ctx, shipmentID)
}

func (s *ShipmentService) GetShipment(ctx context.Context, id string) (*ShipmentResponse, error) {
	var vendorID sql.NullString
	var sh ShipmentResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, order_id, tracking_number, warehouse_id, vendor_id, created_by_id, created_at, updated_at
		FROM shipping_shipments WHERE id = $1
	`, id).Scan(&sh.ID, &sh.OrderID, &sh.TrackingNumber, &sh.WarehouseID, &vendorID, &sh.CreatedByID, &sh.CreatedAt, &sh.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shipment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}
	if vendorID.Valid {
		sh.VendorID = &vendorID.String
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, shipment_id, order_item_id, product_id, quantity
		FROM shipping_shipment_items WHERE shipment_id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item ShipmentItemResponse
		if err := rows.Scan(&item.ID, &item.ShipmentID, &item.OrderItemID, &item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		sh.Items = append(sh.Items, item)
	}
	if sh.Items == nil {
		sh.Items = []ShipmentItemResponse{}
	}
	return &sh, nil
}

func (s *ShipmentService) ListShipments(ctx context.Context, page, pageSize int) ([]ShipmentResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM shipping_shipments`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count shipments: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT id FROM shipping_shipments ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list shipments: %w", err)
	}
	defer rows.Close()

	var shipments []ShipmentResponse
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		sh, err := s.GetShipment(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		shipments = append(shipments, *sh)
	}
	if shipments == nil {
		shipments = []ShipmentResponse{}
	}
	return shipments, total, nil
}
