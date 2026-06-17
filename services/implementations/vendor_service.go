package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type VendorService struct {
	db *sql.DB
}

func NewVendorService(db *sql.DB) *VendorService {
	return &VendorService{db: db}
}

type VendorResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Email       string    `json:"email"`
	IsActive    bool      `json:"isActive"`
	IsDeleted   bool      `json:"isDeleted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateVendorRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Email       string `json:"email"`
}

type UpdateVendorRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Email       string `json:"email"`
	IsActive    *bool  `json:"isActive"`
}

func (s *VendorService) FindAll(ctx context.Context) ([]VendorResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(slug, ''), COALESCE(description, ''), COALESCE(email, ''), is_active, is_deleted, created_at, updated_at
		FROM identity_vendors WHERE is_deleted = false ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list vendors: %w", err)
	}
	defer rows.Close()

	var vendors []VendorResponse
	for rows.Next() {
		var v VendorResponse
		if err := rows.Scan(&v.ID, &v.Name, &v.Slug, &v.Description, &v.Email, &v.IsActive, &v.IsDeleted, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan vendor: %w", err)
		}
		vendors = append(vendors, v)
	}
	if vendors == nil {
		vendors = []VendorResponse{}
	}
	return vendors, nil
}

func (s *VendorService) FindByID(ctx context.Context, id uint) (*VendorResponse, error) {
	var v VendorResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, COALESCE(slug, ''), COALESCE(description, ''), COALESCE(email, ''), is_active, is_deleted, created_at, updated_at
		FROM identity_vendors WHERE id = $1 AND is_deleted = false
	`, id).Scan(&v.ID, &v.Name, &v.Slug, &v.Description, &v.Email, &v.IsActive, &v.IsDeleted, &v.CreatedAt, &v.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("vendor not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	return &v, nil
}

func (s *VendorService) Create(ctx context.Context, req CreateVendorRequest) (*VendorResponse, error) {
	var v VendorResponse
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO identity_vendors (name, slug, description, email, is_active, is_deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, true, false, NOW(), NOW())
		RETURNING id, name, COALESCE(slug, ''), COALESCE(description, ''), COALESCE(email, ''), is_active, is_deleted, created_at, updated_at
	`, req.Name, req.Slug, req.Description, req.Email).Scan(
		&v.ID, &v.Name, &v.Slug, &v.Description, &v.Email, &v.IsActive, &v.IsDeleted, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create vendor: %w", err)
	}
	return &v, nil
}

func (s *VendorService) Update(ctx context.Context, id uint, req UpdateVendorRequest) (*VendorResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	_, err := s.db.ExecContext(ctx, `
		UPDATE identity_vendors SET name = $1, slug = $2, description = $3, email = $4, is_active = $5, updated_at = NOW()
		WHERE id = $6 AND is_deleted = false
	`, req.Name, req.Slug, req.Description, req.Email, isActive, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update vendor: %w", err)
	}

	return s.FindByID(ctx, id)
}

func (s *VendorService) Delete(ctx context.Context, id uint) error {
	result, err := s.db.ExecContext(ctx, `UPDATE identity_vendors SET is_deleted = true, updated_at = NOW() WHERE id = $1 AND is_deleted = false`, id)
	if err != nil {
		return fmt.Errorf("failed to delete vendor: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("vendor not found")
	}
	return nil
}

func (s *VendorService) Paginate(ctx context.Context, page, pageSize int) ([]VendorResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM identity_vendors WHERE is_deleted = false`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vendors: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(slug, ''), COALESCE(description, ''), COALESCE(email, ''), is_active, is_deleted, created_at, updated_at
		FROM identity_vendors WHERE is_deleted = false
		ORDER BY name ASC LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list vendors: %w", err)
	}
	defer rows.Close()

	var vendors []VendorResponse
	for rows.Next() {
		var v VendorResponse
		if err := rows.Scan(&v.ID, &v.Name, &v.Slug, &v.Description, &v.Email, &v.IsActive, &v.IsDeleted, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan vendor: %w", err)
		}
		vendors = append(vendors, v)
	}
	if vendors == nil {
		vendors = []VendorResponse{}
	}
	return vendors, total, nil
}
