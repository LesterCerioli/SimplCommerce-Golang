package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type StockService struct {
	db *sql.DB
}

func NewStockService(db *sql.DB) *StockService {
	return &StockService{db: db}
}

type StockResponse struct {
	ProductID        uint    `json:"productId"`
	TotalStock       int     `json:"totalStock"`
	Available        int     `json:"available"`
	ReservedQuantity int     `json:"reservedQuantity"`
}

type StockHistoryResponse struct {
	ID               uint      `json:"id"`
	ProductID        uint      `json:"productId"`
	WarehouseID      uint      `json:"warehouseId"`
	CreatedByID      uint      `json:"createdById"`
	AdjustedQuantity int       `json:"adjustedQuantity"`
	Note             string    `json:"note"`
	CreatedAt        time.Time `json:"createdAt"`
}

func (s *StockService) GetStock(ctx context.Context, productID uint) (*StockResponse, error) {
	var totalStock, reserved int
	err := s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(quantity),0), COALESCE(SUM(reserved_quantity),0)
		FROM inventory_stocks WHERE product_id = $1
	`, productID).Scan(&totalStock, &reserved)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}

	return &StockResponse{
		ProductID:        productID,
		TotalStock:       totalStock,
		Available:        totalStock - reserved,
		ReservedQuantity: reserved,
	}, nil
}

func (s *StockService) UpdateStock(ctx context.Context, productID, warehouseID uint, quantity int, createdByID uint, note string) (*StockResponse, error) {
	var existingID uint
	var existingQty int
	err := s.db.QueryRowContext(ctx, `
		SELECT id, quantity FROM inventory_stocks WHERE product_id = $1 AND warehouse_id = $2
	`, productID, warehouseID).Scan(&existingID, &existingQty)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	if err == sql.ErrNoRows {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO inventory_stocks (product_id, warehouse_id, quantity, reserved_quantity, created_at, updated_at)
			VALUES ($1,$2,$3,0,NOW(),NOW())
		`, productID, warehouseID, quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to create stock: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO inventory_stock_histories (product_id, warehouse_id, created_by_id, adjusted_quantity, note, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
		`, productID, warehouseID, createdByID, quantity, note)
		if err != nil {
			return nil, fmt.Errorf("failed to create history: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to find stock: %w", err)
	} else {
		adjustment := quantity - existingQty
		_, err = tx.ExecContext(ctx, `
			UPDATE inventory_stocks SET quantity = $1, updated_at = NOW()
			WHERE product_id = $2 AND warehouse_id = $3
		`, quantity, productID, warehouseID)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}

		if adjustment != 0 {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO inventory_stock_histories (product_id, warehouse_id, created_by_id, adjusted_quantity, note, created_at, updated_at)
				VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
			`, productID, warehouseID, createdByID, adjustment, note)
			if err != nil {
				return nil, fmt.Errorf("failed to create history: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return s.GetStock(ctx, productID)
}

func (s *StockService) AdjustStock(ctx context.Context, productID, warehouseID uint, adjustment int, createdByID uint, note string) (*StockResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	var existingID uint
	var existingQty int
	err = tx.QueryRowContext(ctx, `
		SELECT id, quantity FROM inventory_stocks WHERE product_id = $1 AND warehouse_id = $2
	`, productID, warehouseID).Scan(&existingID, &existingQty)

	if err == sql.ErrNoRows {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO inventory_stocks (product_id, warehouse_id, quantity, reserved_quantity, created_at, updated_at)
			VALUES ($1,$2,$3,0,NOW(),NOW())
		`, productID, warehouseID, adjustment)
		if err != nil {
			return nil, fmt.Errorf("failed to create stock: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to find stock: %w", err)
	} else {
		_, err = tx.ExecContext(ctx, `
			UPDATE inventory_stocks SET quantity = quantity + $1, updated_at = NOW()
			WHERE product_id = $2 AND warehouse_id = $3
		`, adjustment, productID, warehouseID)
		if err != nil {
			return nil, fmt.Errorf("failed to adjust stock: %w", err)
		}
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO inventory_stock_histories (product_id, warehouse_id, created_by_id, adjusted_quantity, note, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
	`, productID, warehouseID, createdByID, adjustment, note)
	if err != nil {
		return nil, fmt.Errorf("failed to create history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return s.GetStock(ctx, productID)
}

func (s *StockService) GetStockHistory(ctx context.Context, productID uint, page, pageSize int) ([]StockHistoryResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM inventory_stock_histories WHERE product_id = $1`, productID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stock history: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, product_id, warehouse_id, created_by_id, adjusted_quantity, COALESCE(note,''), created_at
		FROM inventory_stock_histories WHERE product_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, productID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list stock history: %w", err)
	}
	defer rows.Close()

	var history []StockHistoryResponse
	for rows.Next() {
		var h StockHistoryResponse
		if err := rows.Scan(&h.ID, &h.ProductID, &h.WarehouseID, &h.CreatedByID, &h.AdjustedQuantity, &h.Note, &h.CreatedAt); err != nil {
			return nil, 0, err
		}
		history = append(history, h)
	}
	if history == nil {
		history = []StockHistoryResponse{}
	}
	return history, total, nil
}
