package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type RoleService struct {
	db *sql.DB
}

func NewRoleService(db *sql.DB) *RoleService {
	return &RoleService{db: db}
}

type RoleResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	NormalizedName   string    `json:"normalizedName"`
	ConcurrencyStamp string    `json:"concurrencyStamp"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (s *RoleService) Create(ctx context.Context, name string) (*RoleResponse, error) {
	normalizedName := strings.ToUpper(name)
	concurrencyStamp := generateUUID()

	var role RoleResponse
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO identity_roles (name, normalized_name, concurrency_stamp, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, name, normalized_name, concurrency_stamp, created_at, updated_at
	`, name, normalizedName, concurrencyStamp).Scan(
		&role.ID, &role.Name, &role.NormalizedName, &role.ConcurrencyStamp, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &role, nil
}

func (s *RoleService) FindAll(ctx context.Context) ([]RoleResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, normalized_name, concurrency_stamp, created_at, updated_at
		FROM identity_roles ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var roles []RoleResponse
	for rows.Next() {
		var r RoleResponse
		if err := rows.Scan(&r.ID, &r.Name, &r.NormalizedName, &r.ConcurrencyStamp, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, r)
	}
	if roles == nil {
		roles = []RoleResponse{}
	}
	return roles, nil
}

func (s *RoleService) FindByID(ctx context.Context, id string) (*RoleResponse, error) {
	var r RoleResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, normalized_name, concurrency_stamp, created_at, updated_at
		FROM identity_roles WHERE id = $1
	`, id).Scan(&r.ID, &r.Name, &r.NormalizedName, &r.ConcurrencyStamp, &r.CreatedAt, &r.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("role not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	return &r, nil
}

func (s *RoleService) FindByName(ctx context.Context, name string) (*RoleResponse, error) {
	var r RoleResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, normalized_name, concurrency_stamp, created_at, updated_at
		FROM identity_roles WHERE name = $1
	`, name).Scan(&r.ID, &r.Name, &r.NormalizedName, &r.ConcurrencyStamp, &r.CreatedAt, &r.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("role not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find role: %w", err)
	}
	return &r, nil
}

func (s *RoleService) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM identity_user_roles WHERE role_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to remove role assignments: %w", err)
	}

	result, err := s.db.ExecContext(ctx, `DELETE FROM identity_roles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("role not found")
	}
	return nil
}

func (s *RoleService) AssignToUser(ctx context.Context, userID, roleID string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO identity_user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}
	return nil
}

func (s *RoleService) RemoveFromUser(ctx context.Context, userID, roleID string) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM identity_user_roles WHERE user_id = $1 AND role_id = $2
	`, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}
	return nil
}
