package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

type UserListItem struct {
	ID          string    `json:"id"`
	UserGuid    string    `json:"userGuid"`
	FullName    string    `json:"fullName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Culture     string    `json:"culture"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	VendorID    *string   `json:"vendorId"`
	Roles       []string  `json:"roles"`
}

type UpdateUserRequest struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	PhoneNumber string `json:"phoneNumber"`
	Culture     string `json:"culture"`
}

func (s *UserService) FindByID(ctx context.Context, id string) (*UserResponse, error) {
	user, _, err := s.findUserWithRolesByID(ctx, s.db, id)
	return user, err
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*UserResponse, error) {
	var user UserResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, email, full_name FROM identity_users
		WHERE email = $1 AND is_deleted = false
	`, email).Scan(&user.ID, &user.Email, &user.FullName)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT r.name FROM identity_roles r
		JOIN identity_user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find roles: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			user.Roles = append(user.Roles, name)
		}
	}
	if user.Roles == nil {
		user.Roles = []string{}
	}
	return &user, nil
}

func (s *UserService) FindAll(ctx context.Context) ([]UserListItem, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_guid, full_name, email, COALESCE(phone_number, ''), COALESCE(culture, ''), created_at, updated_at, vendor_id
		FROM identity_users WHERE is_deleted = false ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []UserListItem
	for rows.Next() {
		var u UserListItem
		var vendorID sql.NullString
		err := rows.Scan(&u.ID, &u.UserGuid, &u.FullName, &u.Email, &u.PhoneNumber, &u.Culture, &u.CreatedAt, &u.UpdatedAt, &vendorID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		if vendorID.Valid {
			u.VendorID = &vendorID.String
		}
		users = append(users, u)
	}
	if users == nil {
		users = []UserListItem{}
	}
	return users, nil
}

func (s *UserService) Paginate(ctx context.Context, page, pageSize int) ([]UserListItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM identity_users WHERE is_deleted = false`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_guid, full_name, email, COALESCE(phone_number, ''), COALESCE(culture, ''), created_at, updated_at, vendor_id
		FROM identity_users WHERE is_deleted = false
		ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []UserListItem
	for rows.Next() {
		var u UserListItem
		var vendorID sql.NullString
		err := rows.Scan(&u.ID, &u.UserGuid, &u.FullName, &u.Email, &u.PhoneNumber, &u.Culture, &u.CreatedAt, &u.UpdatedAt, &vendorID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		if vendorID.Valid {
			u.VendorID = &vendorID.String
		}
		users = append(users, u)
	}
	if users == nil {
		users = []UserListItem{}
	}
	return users, total, nil
}

func (s *UserService) Update(ctx context.Context, id string, req UpdateUserRequest) (*UserResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if req.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		_, err = tx.ExecContext(ctx, `
			UPDATE identity_users SET full_name = $1, email = $2, password_hash = $3, phone_number = $4, culture = $5, updated_at = NOW()
			WHERE id = $6 AND is_deleted = false
		`, req.FullName, req.Email, string(passwordHash), req.PhoneNumber, req.Culture, id)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		_, err = tx.ExecContext(ctx, `
			UPDATE identity_users SET full_name = $1, email = $2, phone_number = $3, culture = $4, updated_at = NOW()
			WHERE id = $5 AND is_deleted = false
		`, req.FullName, req.Email, req.PhoneNumber, req.Culture, id)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return s.FindByID(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE identity_users SET is_deleted = true, updated_at = NOW() WHERE id = $1 AND is_deleted = false`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *UserService) findUserWithRolesByID(ctx context.Context, db *sql.DB, userID string) (*UserResponse, []string, error) {
	var user UserResponse
	err := db.QueryRowContext(ctx, `
		SELECT id, email, full_name FROM identity_users WHERE id = $1 AND is_deleted = false
	`, userID).Scan(&user.ID, &user.Email, &user.FullName)
	if err == sql.ErrNoRows {
		return nil, nil, ErrUserNotFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find user: %w", err)
	}

	rows, err := db.QueryContext(ctx, `
		SELECT r.name FROM identity_roles r
		JOIN identity_user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find roles: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			user.Roles = append(user.Roles, name)
		}
	}
	if user.Roles == nil {
		user.Roles = []string{}
	}

	return &user, user.Roles, nil
}
