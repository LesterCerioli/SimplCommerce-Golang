package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type MediaService struct {
	db *sql.DB
}

func NewMediaService(db *sql.DB) *MediaService {
	return &MediaService{db: db}
}

type MediaResponse struct {
	ID        string    `json:"id"`
	Caption   string    `json:"caption"`
	FileSize  int64     `json:"fileSize"`
	FileName  string    `json:"fileName"`
	MediaType int       `json:"mediaType"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateMediaRequest struct {
	Caption   string `json:"caption"`
	FileSize  int64  `json:"fileSize"`
	FileName  string `json:"fileName"`
	MediaType int    `json:"mediaType"`
}

func (s *MediaService) FindAll(ctx context.Context) ([]MediaResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, COALESCE(caption, ''), file_size, COALESCE(file_name, ''), media_type, created_at, updated_at
		FROM identity_media ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list media: %w", err)
	}
	defer rows.Close()

	var media []MediaResponse
	for rows.Next() {
		var m MediaResponse
		if err := rows.Scan(&m.ID, &m.Caption, &m.FileSize, &m.FileName, &m.MediaType, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan media: %w", err)
		}
		media = append(media, m)
	}
	if media == nil {
		media = []MediaResponse{}
	}
	return media, nil
}

func (s *MediaService) FindByID(ctx context.Context, id string) (*MediaResponse, error) {
	var m MediaResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, COALESCE(caption, ''), file_size, COALESCE(file_name, ''), media_type, created_at, updated_at
		FROM identity_media WHERE id = $1
	`, id).Scan(&m.ID, &m.Caption, &m.FileSize, &m.FileName, &m.MediaType, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("media not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}
	return &m, nil
}

func (s *MediaService) Create(ctx context.Context, req CreateMediaRequest) (*MediaResponse, error) {
	var m MediaResponse
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO identity_media (caption, file_size, file_name, media_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, COALESCE(caption, ''), file_size, COALESCE(file_name, ''), media_type, created_at, updated_at
	`, req.Caption, req.FileSize, req.FileName, req.MediaType).Scan(
		&m.ID, &m.Caption, &m.FileSize, &m.FileName, &m.MediaType, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create media: %w", err)
	}
	return &m, nil
}

func (s *MediaService) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM identity_media WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("media not found")
	}
	return nil
}

func (s *MediaService) Paginate(ctx context.Context, page, pageSize int) ([]MediaResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM identity_media`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count media: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, COALESCE(caption, ''), file_size, COALESCE(file_name, ''), media_type, created_at, updated_at
		FROM identity_media ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list media: %w", err)
	}
	defer rows.Close()

	var media []MediaResponse
	for rows.Next() {
		var m MediaResponse
		if err := rows.Scan(&m.ID, &m.Caption, &m.FileSize, &m.FileName, &m.MediaType, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan media: %w", err)
		}
		media = append(media, m)
	}
	if media == nil {
		media = []MediaResponse{}
	}
	return media, total, nil
}
