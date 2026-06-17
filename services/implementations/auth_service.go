package implementations

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUserNotFound        = errors.New("user not found")
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
	UserID       string `json:"userId"`
}

type UserResponse struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	FullName string   `json:"fullName"`
	Roles    []string `json:"roles"`
}

type jwtClaims struct {
	UserID   string   `json:"userId"`
	Email    string   `json:"email"`
	FullName string   `json:"fullName"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

type AuthService struct {
	db        *sql.DB
	jwtSecret string
	jwtExpiry int
}

func NewAuthService(db *sql.DB, jwtSecret string, jwtExpiry int) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM identity_users WHERE email = $1 AND is_deleted = false)`, req.Email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var userID string
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO identity_users (user_guid, full_name, email, password_hash, is_deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, false, NOW(), NOW())
		RETURNING id
	`, generateUUID(), req.FullName, req.Email, string(passwordHash)).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	var roleID uint
	err = s.db.QueryRowContext(ctx, `SELECT id FROM identity_roles WHERE normalized_name = 'CUSTOMER'`).Scan(&roleID)
	if err == sql.ErrNoRows {
		concurrencyStamp := generateUUID()
		err = s.db.QueryRowContext(ctx, `
			INSERT INTO identity_roles (name, normalized_name, concurrency_stamp, created_at, updated_at)
			VALUES ('customer', 'CUSTOMER', $1, NOW(), NOW())
			RETURNING id
		`, concurrencyStamp).Scan(&roleID)
		if err != nil {
			return nil, fmt.Errorf("failed to create customer role: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to find customer role: %w", err)
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO identity_user_roles (user_id, role_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, userID, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	return s.generateTokens(ctx, userID)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	var userID string
	var passwordHash string
	err := s.db.QueryRowContext(ctx, `
		SELECT id, password_hash FROM identity_users
		WHERE email = $1 AND is_deleted = false
	`, email).Scan(&userID, &passwordHash)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokens(ctx, userID)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, refresh_token_hash FROM identity_users WHERE is_deleted = false`)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}
	defer rows.Close()

	var foundID string
	found := false
	for rows.Next() {
		var id string
		var hash sql.NullString
		if err := rows.Scan(&id, &hash); err != nil {
			continue
		}
		if !hash.Valid || hash.String == "" {
			continue
		}
		if bcrypt.CompareHashAndPassword([]byte(hash.String), []byte(refreshToken)) == nil {
			foundID = id
			found = true
			break
		}
	}
	if !found {
		return nil, ErrInvalidRefreshToken
	}

	return s.generateTokens(ctx, foundID)
}

func (s *AuthService) ValidateToken(tokenString string) (string, string, []string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return "", "", nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return "", "", nil, errors.New("invalid token claims")
	}

	return claims.UserID, claims.Email, claims.Roles, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (*UserResponse, error) {
	user, _, err := s.findUserWithRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) generateTokens(ctx context.Context, userID string) (*AuthResponse, error) {
	var email, fullName string
	err := s.db.QueryRowContext(ctx, `
		SELECT email, full_name FROM identity_users WHERE id = $1 AND is_deleted = false
	`, userID).Scan(&email, &fullName)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	roleNames := []string{}
	roleRows, err := s.db.QueryContext(ctx, `
		SELECT r.name FROM identity_roles r
		JOIN identity_user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`, userID)
	if err == nil {
		defer roleRows.Close()
		for roleRows.Next() {
			var name string
			if err := roleRows.Scan(&name); err == nil {
				roleNames = append(roleNames, name)
			}
		}
	}

	expiresAt := time.Now().Add(time.Duration(s.jwtExpiry) * time.Hour)
	claims := &jwtClaims{
		UserID:   userID,
		Email:    email,
		FullName: fullName,
		Roles:    roleNames,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	_, err = s.db.ExecContext(ctx, `UPDATE identity_users SET refresh_token_hash = $1, updated_at = NOW() WHERE id = $2`, string(refreshHash), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
		UserID:       userID,
	}, nil
}

func (s *AuthService) findUserWithRoles(ctx context.Context, userID string) (*UserResponse, []string, error) {
	var user UserResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, email, full_name FROM identity_users WHERE id = $1 AND is_deleted = false
	`, userID).Scan(&user.ID, &user.Email, &user.FullName)
	if err == sql.ErrNoRows {
		return nil, nil, ErrUserNotFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find user: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, `
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

func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return hex.EncodeToString(b[:4]) + "-" + hex.EncodeToString(b[4:6]) + "-" + hex.EncodeToString(b[6:8]) + "-" + hex.EncodeToString(b[8:10]) + "-" + hex.EncodeToString(b[10:])
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
