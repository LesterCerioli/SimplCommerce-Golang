package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v3"
)

type BrandService struct {
	db *sql.DB
}

type BrandResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	IsPublished bool       `json:"isPublished"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type CreateBrandRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublished bool   `json:"isPublished,omitempty"`
}

type UpdateBrandRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	IsPublished bool   `json:"isPublished,omitempty"`
}

func NewBrandService(db *sql.DB) *BrandService {
	return &BrandService{db: db}
}

func (s *BrandService) GetBrands(ctx context.Context, page, pageSize int) ([]*BrandResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM catalog_brands WHERE is_deleted = false").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, slug, description, is_published, created_at, updated_at FROM catalog_brands WHERE is_deleted = false ORDER BY id ASC LIMIT $1 OFFSET $2",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var brands []*BrandResponse
	for rows.Next() {
		var b BrandResponse
		if err := rows.Scan(&b.ID, &b.Name, &b.Slug, &b.Description, &b.IsPublished, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, 0, err
		}
		brands = append(brands, &b)
	}
	if brands == nil {
		brands = []*BrandResponse{}
	}
	return brands, total, nil
}

func (s *BrandService) GetBrandBySlug(ctx context.Context, slug string) (*BrandResponse, error) {
	var b BrandResponse
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, slug, description, is_published, created_at, updated_at FROM catalog_brands WHERE slug = $1 AND is_deleted = false",
		slug,
	).Scan(&b.ID, &b.Name, &b.Slug, &b.Description, &b.IsPublished, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fiber.NewError(fiber.StatusNotFound, "Brand not found")
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *BrandService) GetBrandByID(ctx context.Context, id uint) (*BrandResponse, error) {
	var b BrandResponse
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, slug, description, is_published, created_at, updated_at FROM catalog_brands WHERE id = $1 AND is_deleted = false",
		id,
	).Scan(&b.ID, &b.Name, &b.Slug, &b.Description, &b.IsPublished, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fiber.NewError(fiber.StatusNotFound, "Brand not found")
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *BrandService) CreateBrand(ctx context.Context, req *CreateBrandRequest) (*BrandResponse, error) {
	slug := req.Slug
	if slug == "" {
		slug = GenerateSlug(req.Name)
	}

	var b BrandResponse
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO catalog_brands (name, slug, description, is_published, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())
		 RETURNING id, name, slug, description, is_published, created_at, updated_at`,
		req.Name, slug, req.Description, req.IsPublished,
	).Scan(&b.ID, &b.Name, &b.Slug, &b.Description, &b.IsPublished, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *BrandService) UpdateBrand(ctx context.Context, id uint, req *UpdateBrandRequest) (*BrandResponse, error) {
	existing, err := s.GetBrandByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := existing.Name
	if req.Name != "" {
		name = req.Name
	}
	slug := req.Slug
	if slug == "" && req.Name != "" {
		slug = GenerateSlug(name)
	} else if slug == "" {
		slug = existing.Slug
	}
	desc := req.Description
	isPub := req.IsPublished

	var b BrandResponse
	err = s.db.QueryRowContext(ctx,
		`UPDATE catalog_brands SET name = $1, slug = $2, description = $3, is_published = $4, updated_at = NOW()
		 WHERE id = $5 AND is_deleted = false
		 RETURNING id, name, slug, description, is_published, created_at, updated_at`,
		name, slug, desc, isPub, id,
	).Scan(&b.ID, &b.Name, &b.Slug, &b.Description, &b.IsPublished, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fiber.NewError(fiber.StatusNotFound, "Brand not found")
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *BrandService) DeleteBrand(ctx context.Context, id uint) error {
	result, err := s.db.ExecContext(ctx, "UPDATE catalog_brands SET is_deleted = true, updated_at = NOW() WHERE id = $1 AND is_deleted = false", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Brand not found")
	}
	return nil
}
