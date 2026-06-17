package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v3"
)

type ProductAttributeService struct {
	db *sql.DB
}

type ProductAttributeResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	GroupID   uint      `json:"groupId"`
	GroupName string    `json:"groupName,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateProductAttributeRequest struct {
	Name    string `json:"name"`
	GroupID uint   `json:"groupId"`
}

type ProductAttributeGroupService struct {
	db *sql.DB
}

type ProductAttributeGroupResponse struct {
	ID         uint                       `json:"id"`
	Name       string                     `json:"name"`
	CreatedAt  time.Time                  `json:"createdAt"`
	UpdatedAt  time.Time                  `json:"updatedAt"`
	Attributes []*ProductAttributeResponse `json:"attributes,omitempty"`
}

type CreateProductAttributeGroupRequest struct {
	Name string `json:"name"`
}

type UpdateProductAttributeGroupRequest struct {
	Name string `json:"name"`
}

func NewProductAttributeService(db *sql.DB) *ProductAttributeService {
	return &ProductAttributeService{db: db}
}

func NewProductAttributeGroupService(db *sql.DB) *ProductAttributeGroupService {
	return &ProductAttributeGroupService{db: db}
}

func (s *ProductAttributeService) GetAttributes(ctx context.Context, page, pageSize int) ([]*ProductAttributeResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM catalog_product_attributes").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx,
		`SELECT a.id, a.name, a.group_id, COALESCE(g.name, '') as group_name, a.created_at, a.updated_at
		 FROM catalog_product_attributes a
		 LEFT JOIN catalog_product_attribute_groups g ON a.group_id = g.id
		 ORDER BY a.id ASC LIMIT $1 OFFSET $2`,
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var attrs []*ProductAttributeResponse
	for rows.Next() {
		var a ProductAttributeResponse
		if err := rows.Scan(&a.ID, &a.Name, &a.GroupID, &a.GroupName, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, err
		}
		attrs = append(attrs, &a)
	}
	if attrs == nil {
		attrs = []*ProductAttributeResponse{}
	}
	return attrs, total, nil
}

func (s *ProductAttributeService) CreateAttribute(ctx context.Context, req *CreateProductAttributeRequest) (*ProductAttributeResponse, error) {
	var a ProductAttributeResponse
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO catalog_product_attributes (name, group_id, created_at, updated_at)
		 VALUES ($1, $2, NOW(), NOW())
		 RETURNING id, name, group_id, created_at, updated_at`,
		req.Name, req.GroupID,
	).Scan(&a.ID, &a.Name, &a.GroupID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if a.GroupID > 0 {
		_ = s.db.QueryRowContext(ctx, "SELECT name FROM catalog_product_attribute_groups WHERE id = $1", a.GroupID).Scan(&a.GroupName)
	}
	return &a, nil
}

func (s *ProductAttributeService) DeleteAttribute(ctx context.Context, id uint) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM catalog_product_attributes WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Attribute not found")
	}
	return nil
}

func (s *ProductAttributeGroupService) GetGroups(ctx context.Context, page, pageSize int) ([]*ProductAttributeGroupResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM catalog_product_attribute_groups").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, created_at, updated_at FROM catalog_product_attribute_groups ORDER BY id ASC LIMIT $1 OFFSET $2",
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var groups []*ProductAttributeGroupResponse
	for rows.Next() {
		var g ProductAttributeGroupResponse
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, 0, err
		}
		groups = append(groups, &g)
	}
	if groups == nil {
		groups = []*ProductAttributeGroupResponse{}
	}

	for _, g := range groups {
		attrRows, err := s.db.QueryContext(ctx,
			"SELECT id, name, group_id, created_at, updated_at FROM catalog_product_attributes WHERE group_id = $1 ORDER BY id ASC",
			g.ID,
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
			a.GroupName = g.Name
			g.Attributes = append(g.Attributes, &a)
		}
		attrRows.Close()
		if g.Attributes == nil {
			g.Attributes = []*ProductAttributeResponse{}
		}
	}

	return groups, total, nil
}

func (s *ProductAttributeGroupService) CreateGroup(ctx context.Context, req *CreateProductAttributeGroupRequest) (*ProductAttributeGroupResponse, error) {
	var g ProductAttributeGroupResponse
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO catalog_product_attribute_groups (name, created_at, updated_at)
		 VALUES ($1, NOW(), NOW())
		 RETURNING id, name, created_at, updated_at`,
		req.Name,
	).Scan(&g.ID, &g.Name, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	g.Attributes = []*ProductAttributeResponse{}
	return &g, nil
}

func (s *ProductAttributeGroupService) UpdateGroup(ctx context.Context, id uint, req *UpdateProductAttributeGroupRequest) (*ProductAttributeGroupResponse, error) {
	var g ProductAttributeGroupResponse
	err := s.db.QueryRowContext(ctx,
		`UPDATE catalog_product_attribute_groups SET name = $1, updated_at = NOW()
		 WHERE id = $2
		 RETURNING id, name, created_at, updated_at`,
		req.Name, id,
	).Scan(&g.ID, &g.Name, &g.CreatedAt, &g.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fiber.NewError(fiber.StatusNotFound, "Attribute group not found")
	}
	if err != nil {
		return nil, err
	}
	g.Attributes = []*ProductAttributeResponse{}
	return &g, nil
}
