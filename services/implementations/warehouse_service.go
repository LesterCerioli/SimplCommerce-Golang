package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type WarehouseService struct {
	db *sql.DB
}

func NewWarehouseService(db *sql.DB) *WarehouseService {
	return &WarehouseService{db: db}
}

type WarehouseResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	AddressID uint      `json:"addressId"`
	VendorID  *uint     `json:"vendorId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *WarehouseService) List(ctx context.Context) ([]WarehouseResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, address_id, vendor_id, created_at, updated_at
		FROM inventory_warehouses ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list warehouses: %w", err)
	}
	defer rows.Close()

	var warehouses []WarehouseResponse
	for rows.Next() {
		var w WarehouseResponse
		var vendorID sql.NullInt64
		if err := rows.Scan(&w.ID, &w.Name, &w.AddressID, &vendorID, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		if vendorID.Valid {
			v := uint(vendorID.Int64)
			w.VendorID = &v
		}
		warehouses = append(warehouses, w)
	}
	if warehouses == nil {
		warehouses = []WarehouseResponse{}
	}
	return warehouses, nil
}

func (s *WarehouseService) Create(ctx context.Context, name string, addressID uint) (*WarehouseResponse, error) {
	var w WarehouseResponse
	var vendorID sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO inventory_warehouses (name, address_id, created_at, updated_at)
		VALUES ($1,$2,NOW(),NOW()) RETURNING id, name, address_id, vendor_id, created_at, updated_at
	`, name, addressID).Scan(&w.ID, &w.Name, &w.AddressID, &vendorID, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create warehouse: %w", err)
	}
	if vendorID.Valid {
		v := uint(vendorID.Int64)
		w.VendorID = &v
	}
	return &w, nil
}

func (s *WarehouseService) Update(ctx context.Context, id uint, name string, addressID uint) (*WarehouseResponse, error) {
	var w WarehouseResponse
	var vendorID sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		UPDATE inventory_warehouses SET name = $1, address_id = $2, updated_at = NOW()
		WHERE id = $3 RETURNING id, name, address_id, vendor_id, created_at, updated_at
	`, name, addressID, id).Scan(&w.ID, &w.Name, &w.AddressID, &vendorID, &w.CreatedAt, &w.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("warehouse not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update warehouse: %w", err)
	}
	if vendorID.Valid {
		v := uint(vendorID.Int64)
		w.VendorID = &v
	}
	return &w, nil
}

func (s *WarehouseService) Delete(ctx context.Context, id uint) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM inventory_warehouses WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("warehouse not found")
	}
	return nil
}

func (s *WarehouseService) GetByID(ctx context.Context, id uint) (*WarehouseResponse, error) {
	var w WarehouseResponse
	var vendorID sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, address_id, vendor_id, created_at, updated_at
		FROM inventory_warehouses WHERE id = $1
	`, id).Scan(&w.ID, &w.Name, &w.AddressID, &vendorID, &w.CreatedAt, &w.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("warehouse not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse: %w", err)
	}
	if vendorID.Valid {
		v := uint(vendorID.Int64)
		w.VendorID = &v
	}
	return &w, nil
}
