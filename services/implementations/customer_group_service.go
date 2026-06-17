package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CustomerGroupService struct {
	db *sql.DB
}

func NewCustomerGroupService(db *sql.DB) *CustomerGroupService {
	return &CustomerGroupService{db: db}
}

type CustomerGroupResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UserCount   int       `json:"userCount"`
}

type CreateCustomerGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCustomerGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"isActive"`
}

func (s *CustomerGroupService) FindAll(ctx context.Context) ([]CustomerGroupResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT cg.id, cg.name, COALESCE(cg.description, ''), cg.is_active, cg.created_at, cg.updated_at,
			(SELECT COUNT(*) FROM identity_customer_group_users cgu WHERE cgu.customer_group_id = cg.id)
		FROM identity_customer_groups cg
		WHERE cg.is_deleted = false
		ORDER BY cg.name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list customer groups: %w", err)
	}
	defer rows.Close()

	var groups []CustomerGroupResponse
	for rows.Next() {
		var g CustomerGroupResponse
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.IsActive, &g.CreatedAt, &g.UpdatedAt, &g.UserCount); err != nil {
			return nil, fmt.Errorf("failed to scan customer group: %w", err)
		}
		groups = append(groups, g)
	}
	if groups == nil {
		groups = []CustomerGroupResponse{}
	}
	return groups, nil
}

func (s *CustomerGroupService) FindByID(ctx context.Context, id uint) (*CustomerGroupResponse, error) {
	var g CustomerGroupResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT cg.id, cg.name, COALESCE(cg.description, ''), cg.is_active, cg.created_at, cg.updated_at,
			(SELECT COUNT(*) FROM identity_customer_group_users cgu WHERE cgu.customer_group_id = cg.id)
		FROM identity_customer_groups cg
		WHERE cg.id = $1 AND cg.is_deleted = false
	`, id).Scan(&g.ID, &g.Name, &g.Description, &g.IsActive, &g.CreatedAt, &g.UpdatedAt, &g.UserCount)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer group not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find customer group: %w", err)
	}
	return &g, nil
}

func (s *CustomerGroupService) Create(ctx context.Context, req CreateCustomerGroupRequest) (*CustomerGroupResponse, error) {
	var g CustomerGroupResponse
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO identity_customer_groups (name, description, is_active, is_deleted, created_at, updated_at)
		VALUES ($1, $2, true, false, NOW(), NOW())
		RETURNING id, name, COALESCE(description, ''), is_active, created_at, updated_at
	`, req.Name, req.Description).Scan(
		&g.ID, &g.Name, &g.Description, &g.IsActive, &g.CreatedAt, &g.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer group: %w", err)
	}
	g.UserCount = 0
	return &g, nil
}

func (s *CustomerGroupService) Update(ctx context.Context, id uint, req UpdateCustomerGroupRequest) (*CustomerGroupResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	_, err := s.db.ExecContext(ctx, `
		UPDATE identity_customer_groups SET name = $1, description = $2, is_active = $3, updated_at = NOW()
		WHERE id = $4 AND is_deleted = false
	`, req.Name, req.Description, isActive, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update customer group: %w", err)
	}

	return s.FindByID(ctx, id)
}

func (s *CustomerGroupService) Delete(ctx context.Context, id uint) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM identity_customer_group_users WHERE customer_group_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to remove group users: %w", err)
	}

	result, err := s.db.ExecContext(ctx, `DELETE FROM identity_customer_groups WHERE id = $1 AND is_deleted = false`, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer group: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("customer group not found")
	}
	return nil
}

func (s *CustomerGroupService) Paginate(ctx context.Context, page, pageSize int) ([]CustomerGroupResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM identity_customer_groups WHERE is_deleted = false`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customer groups: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT cg.id, cg.name, COALESCE(cg.description, ''), cg.is_active, cg.created_at, cg.updated_at,
			(SELECT COUNT(*) FROM identity_customer_group_users cgu WHERE cgu.customer_group_id = cg.id)
		FROM identity_customer_groups cg
		WHERE cg.is_deleted = false
		ORDER BY cg.name ASC LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list customer groups: %w", err)
	}
	defer rows.Close()

	var groups []CustomerGroupResponse
	for rows.Next() {
		var g CustomerGroupResponse
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.IsActive, &g.CreatedAt, &g.UpdatedAt, &g.UserCount); err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer group: %w", err)
		}
		groups = append(groups, g)
	}
	if groups == nil {
		groups = []CustomerGroupResponse{}
	}
	return groups, total, nil
}
