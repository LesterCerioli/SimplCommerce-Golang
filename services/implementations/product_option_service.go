package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v3"
)

type ProductOptionService struct {
	db *sql.DB
}

type ProductOptionResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateProductOptionRequest struct {
	Name string `json:"name"`
}

type UpdateProductOptionRequest struct {
	Name string `json:"name"`
}

type ProductTemplateService struct {
	db *sql.DB
}

type ProductTemplateResponse struct {
	ID           string                     `json:"id"`
	Name         string                     `json:"name"`
	CreatedAt    time.Time                  `json:"createdAt"`
	UpdatedAt    time.Time                  `json:"updatedAt"`
	AttributeIDs []string                   `json:"attributeIds,omitempty"`
	Attributes   []*ProductAttributeResponse `json:"attributes,omitempty"`
}

type CreateProductTemplateRequest struct {
	Name         string   `json:"name"`
	AttributeIDs []string `json:"attributeIds,omitempty"`
}

func NewProductOptionService(db *sql.DB) *ProductOptionService {
	return &ProductOptionService{db: db}
}

func NewProductTemplateService(db *sql.DB) *ProductTemplateService {
	return &ProductTemplateService{db: db}
}

func (s *ProductOptionService) GetOptions(ctx context.Context, page, pageSize int) ([]*ProductOptionResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM catalog_product_options").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, created_at, updated_at FROM catalog_product_options ORDER BY id ASC LIMIT $1 OFFSET $2",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var opts []*ProductOptionResponse
	for rows.Next() {
		var o ProductOptionResponse
		if err := rows.Scan(&o.ID, &o.Name, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		opts = append(opts, &o)
	}
	if opts == nil {
		opts = []*ProductOptionResponse{}
	}
	return opts, total, nil
}

func (s *ProductOptionService) CreateOption(ctx context.Context, req *CreateProductOptionRequest) (*ProductOptionResponse, error) {
	var o ProductOptionResponse
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO catalog_product_options (name, created_at, updated_at)
		 VALUES ($1, NOW(), NOW())
		 RETURNING id, name, created_at, updated_at`,
		req.Name,
	).Scan(&o.ID, &o.Name, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *ProductOptionService) UpdateOption(ctx context.Context, id string, req *UpdateProductOptionRequest) (*ProductOptionResponse, error) {
	var o ProductOptionResponse
	err := s.db.QueryRowContext(ctx,
		`UPDATE catalog_product_options SET name = $1, updated_at = NOW()
		 WHERE id = $2
		 RETURNING id, name, created_at, updated_at`,
		req.Name, id,
	).Scan(&o.ID, &o.Name, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fiber.NewError(fiber.StatusNotFound, "Option not found")
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *ProductOptionService) DeleteOption(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM catalog_product_options WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Option not found")
	}
	return nil
}

func (s *ProductTemplateService) GetTemplates(ctx context.Context, page, pageSize int) ([]*ProductTemplateResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM catalog_product_templates").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, created_at, updated_at FROM catalog_product_templates ORDER BY id ASC LIMIT $1 OFFSET $2",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*ProductTemplateResponse
	for rows.Next() {
		var t ProductTemplateResponse
		if err := rows.Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		templates = append(templates, &t)
	}
	if templates == nil {
		templates = []*ProductTemplateResponse{}
	}

	for _, t := range templates {
		attrRows, err := s.db.QueryContext(ctx,
			`SELECT a.id, a.name, a.group_id, a.created_at, a.updated_at
			 FROM catalog_product_template_product_attributes tpa
			 JOIN catalog_product_attributes a ON tpa.product_attribute_id = a.id
			 WHERE tpa.product_template_id = $1
			 ORDER BY a.id ASC`,
			t.ID,
		)
		if err != nil {
			return nil, 0, err
		}
		for attrRows.Next() {
			var a ProductAttributeResponse
			if err := attrRows.Scan(&a.ID, &a.Name, &a.GroupID, &a.CreatedAt, &a.UpdatedAt); err != nil {
				attrRows.Close()
				return nil, 0, err
			}
			t.AttributeIDs = append(t.AttributeIDs, a.ID)
			t.Attributes = append(t.Attributes, &a)
		}
		attrRows.Close()
		if t.Attributes == nil {
			t.Attributes = []*ProductAttributeResponse{}
		}
	}

	return templates, total, nil
}

func (s *ProductTemplateService) CreateTemplate(ctx context.Context, req *CreateProductTemplateRequest) (*ProductTemplateResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var t ProductTemplateResponse
	err = tx.QueryRowContext(ctx,
		`INSERT INTO catalog_product_templates (name, created_at, updated_at)
		 VALUES ($1, NOW(), NOW())
		 RETURNING id, name, created_at, updated_at`,
		req.Name,
	).Scan(&t.ID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	for _, attrID := range req.AttributeIDs {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO catalog_product_template_product_attributes (product_template_id, product_attribute_id)
			 VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			t.ID, attrID,
		)
		if err != nil {
			return nil, err
		}
		t.AttributeIDs = append(t.AttributeIDs, attrID)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	if t.Attributes == nil {
		t.Attributes = []*ProductAttributeResponse{}
	}
	return &t, nil
}
