package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PageService struct {
	db *sql.DB
}

func NewPageService(db *sql.DB) *PageService {
	return &PageService{db: db}
}

type PageResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Body        string     `json:"body"`
	IsPublished bool       `json:"isPublished"`
	PublishedOn *time.Time `json:"publishedOn,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (s *PageService) List(ctx context.Context) ([]PageResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, slug, body, is_published, published_on, created_at, updated_at
		FROM cms_pages ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list pages: %w", err)
	}
	defer rows.Close()

	var pages []PageResponse
	for rows.Next() {
		var p PageResponse
		var publishedOn sql.NullTime
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Body, &p.IsPublished, &publishedOn, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if publishedOn.Valid {
			p.PublishedOn = &publishedOn.Time
		}
		pages = append(pages, p)
	}
	if pages == nil {
		pages = []PageResponse{}
	}
	return pages, nil
}

func (s *PageService) GetBySlug(ctx context.Context, slug string) (*PageResponse, error) {
	var p PageResponse
	var publishedOn sql.NullTime
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, slug, body, is_published, published_on, created_at, updated_at
		FROM cms_pages WHERE slug = $1
	`, slug).Scan(&p.ID, &p.Name, &p.Slug, &p.Body, &p.IsPublished, &publishedOn, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("page not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}
	if publishedOn.Valid {
		p.PublishedOn = &publishedOn.Time
	}
	return &p, nil
}

func (s *PageService) GetByID(ctx context.Context, id string) (*PageResponse, error) {
	var p PageResponse
	var publishedOn sql.NullTime
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, slug, body, is_published, published_on, created_at, updated_at
		FROM cms_pages WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &p.Slug, &p.Body, &p.IsPublished, &publishedOn, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("page not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}
	if publishedOn.Valid {
		p.PublishedOn = &publishedOn.Time
	}
	return &p, nil
}

func (s *PageService) Create(ctx context.Context, name, slug, body string, isPublished bool) (*PageResponse, error) {
	if slug == "" {
		slug = GenerateSlug(name)
	}
	var publishedOn *time.Time
	now := time.Now()
	if isPublished {
		publishedOn = &now
	}

	var p PageResponse
	var pubOn sql.NullTime
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO cms_pages (name, slug, body, is_published, published_on, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
		RETURNING id, name, slug, body, is_published, published_on, created_at, updated_at
	`, name, slug, body, isPublished, publishedOn).Scan(&p.ID, &p.Name, &p.Slug, &p.Body, &p.IsPublished, &pubOn, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	if pubOn.Valid {
		p.PublishedOn = &pubOn.Time
	}
	return &p, nil
}

func (s *PageService) Update(ctx context.Context, id string, name, slug, body string, isPublished bool) (*PageResponse, error) {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name == "" {
		name = existing.Name
	}
	if slug == "" {
		slug = existing.Slug
	}
	if body == "" {
		body = existing.Body
	}
	var publishedOn *time.Time
	if isPublished && existing.PublishedOn == nil {
		now := time.Now()
		publishedOn = &now
	} else if !isPublished {
		publishedOn = nil
	} else {
		publishedOn = existing.PublishedOn
	}

	var p PageResponse
	var pubOn sql.NullTime
	err = s.db.QueryRowContext(ctx, `
		UPDATE cms_pages SET name=$1, slug=$2, body=$3, is_published=$4, published_on=$5, updated_at=NOW()
		WHERE id=$6 RETURNING id, name, slug, body, is_published, published_on, created_at, updated_at
	`, name, slug, body, isPublished, publishedOn, id).Scan(&p.ID, &p.Name, &p.Slug, &p.Body, &p.IsPublished, &pubOn, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}
	if pubOn.Valid {
		p.PublishedOn = &pubOn.Time
	}
	return &p, nil
}

func (s *PageService) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM cms_pages WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete page: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("page not found")
	}
	return nil
}
