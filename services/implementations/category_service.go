package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v3"
)

type CategoryService struct {
	db *sql.DB
}

type CategoryResponse struct {
	ID              uint               `json:"id"`
	Name            string             `json:"name"`
	Slug            string             `json:"slug"`
	Description     string             `json:"description"`
	DisplayOrder    int                `json:"displayOrder"`
	IsPublished     bool               `json:"isPublished"`
	IncludeInMenu   bool               `json:"includeInMenu"`
	ParentID        *uint              `json:"parentId,omitempty"`
	ThumbnailImageID *uint             `json:"thumbnailImageId,omitempty"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
	Children        []*CategoryResponse `json:"children,omitempty"`
}

type CreateCategoryRequest struct {
	Name            string `json:"name"`
	Slug            string `json:"slug,omitempty"`
	Description     string `json:"description,omitempty"`
	DisplayOrder    int    `json:"displayOrder,omitempty"`
	IsPublished     bool   `json:"isPublished,omitempty"`
	IncludeInMenu   bool   `json:"includeInMenu,omitempty"`
	ParentID        *uint  `json:"parentId,omitempty"`
	ThumbnailImageID *uint `json:"thumbnailImageId,omitempty"`
}

type UpdateCategoryRequest struct {
	Name            string `json:"name"`
	Slug            string `json:"slug,omitempty"`
	Description     string `json:"description,omitempty"`
	DisplayOrder    int    `json:"displayOrder,omitempty"`
	IsPublished     bool   `json:"isPublished,omitempty"`
	IncludeInMenu   bool   `json:"includeInMenu,omitempty"`
	ParentID        *uint  `json:"parentId,omitempty"`
	ThumbnailImageID *uint `json:"thumbnailImageId,omitempty"`
}

func NewCategoryService(db *sql.DB) *CategoryService {
	return &CategoryService{db: db}
}

func scanCategory(row scannable) (*CategoryResponse, error) {
	var c CategoryResponse
	var parentID sql.NullInt64
	var thumbID sql.NullInt64

	err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.DisplayOrder, &c.IsPublished, &c.IncludeInMenu, &parentID, &thumbID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		pid := uint(parentID.Int64)
		c.ParentID = &pid
	}
	if thumbID.Valid {
		tid := uint(thumbID.Int64)
		c.ThumbnailImageID = &tid
	}
	return &c, nil
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]*CategoryResponse, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, slug, description, display_order, is_published, include_in_menu, parent_id, thumbnail_image_id, created_at, updated_at
		 FROM catalog_categories WHERE is_deleted = false ORDER BY display_order ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []*CategoryResponse
	for rows.Next() {
		c, err := scanCategory(rows)
		if err != nil {
			return nil, err
		}
		all = append(all, c)
	}
	if all == nil {
		all = []*CategoryResponse{}
	}
	return all, nil
}

func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*CategoryResponse, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, slug, description, display_order, is_published, include_in_menu, parent_id, thumbnail_image_id, created_at, updated_at
		 FROM catalog_categories WHERE slug = $1 AND is_deleted = false`,
		slug,
	)
	c, err := scanCategory(row)
	if err == sql.ErrNoRows {
		return nil, fiber.NewError(fiber.StatusNotFound, "Category not found")
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id uint) (*CategoryResponse, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, slug, description, display_order, is_published, include_in_menu, parent_id, thumbnail_image_id, created_at, updated_at
		 FROM catalog_categories WHERE id = $1 AND is_deleted = false`,
		id,
	)
	c, err := scanCategory(row)
	if err == sql.ErrNoRows {
		return nil, fiber.NewError(fiber.StatusNotFound, "Category not found")
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error) {
	slug := req.Slug
	if slug == "" {
		slug = GenerateSlug(req.Name)
	}

	row := s.db.QueryRowContext(ctx,
		`INSERT INTO catalog_categories (name, slug, description, display_order, is_published, include_in_menu, parent_id, thumbnail_image_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		 RETURNING id, name, slug, description, display_order, is_published, include_in_menu, parent_id, thumbnail_image_id, created_at, updated_at`,
		req.Name, slug, req.Description, req.DisplayOrder, req.IsPublished, req.IncludeInMenu, req.ParentID, req.ThumbnailImageID,
	)
	c, err := scanCategory(row)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id uint, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	existing, err := s.GetCategoryByID(ctx, id)
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
	dispOrder := req.DisplayOrder
	isPub := req.IsPublished
	incMenu := req.IncludeInMenu
	parentID := req.ParentID
	thumbID := req.ThumbnailImageID

	row := s.db.QueryRowContext(ctx,
		`UPDATE catalog_categories SET name = $1, slug = $2, description = $3, display_order = $4, is_published = $5, include_in_menu = $6, parent_id = $7, thumbnail_image_id = $8, updated_at = NOW()
		 WHERE id = $9 AND is_deleted = false
		 RETURNING id, name, slug, description, display_order, is_published, include_in_menu, parent_id, thumbnail_image_id, created_at, updated_at`,
		name, slug, desc, dispOrder, isPub, incMenu, parentID, thumbID, id,
	)
	c, err := scanCategory(row)
	if err == sql.ErrNoRows {
		return nil, fiber.NewError(fiber.StatusNotFound, "Category not found")
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uint) error {
	var childCount int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM catalog_categories WHERE parent_id = $1 AND is_deleted = false", id).Scan(&childCount)
	if err != nil {
		return err
	}
	if childCount > 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot delete category with children")
	}

	result, err := s.db.ExecContext(ctx, "UPDATE catalog_categories SET is_deleted = true, updated_at = NOW() WHERE id = $1 AND is_deleted = false", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Category not found")
	}
	return nil
}

func BuildCategoryTree(categories []*CategoryResponse) []*CategoryResponse {
	lookup := make(map[uint]*CategoryResponse)
	for _, c := range categories {
		c.Children = []*CategoryResponse{}
		lookup[c.ID] = c
	}

	var roots []*CategoryResponse
	for _, c := range categories {
		if c.ParentID == nil {
			roots = append(roots, c)
		} else {
			if parent, ok := lookup[*c.ParentID]; ok {
				parent.Children = append(parent.Children, c)
			}
		}
	}
	if roots == nil {
		roots = []*CategoryResponse{}
	}
	return roots
}
