package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/simplcommerce-go/pkg/auth"
	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"fullName" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
	UserID       uint   `json:"userId"`
}

type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*models.User, error)
	Login(ctx context.Context, email, password string) (*TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
}

type authService struct {
	userRepo   repositories.UserRepository
	roleRepo   repositories.RoleRepository
	jwtService auth.JWTService
	password   PasswordService
}

func NewAuthService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, jwtService auth.JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		jwtService: jwtService,
		password:   NewPasswordService(),
	}
}

func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := s.password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: passwordHash,
		UserGuid:     generateUUID(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	customerRole, err := s.roleRepo.FindByName(ctx, "customer")
	if err != nil {
		if err := s.roleRepo.Create(ctx, &models.Role{
			Name:           "customer",
			NormalizedName: "CUSTOMER",
			ConcurrencyStamp: generateUUID(),
		}); err != nil {
			return nil, err
		}
		customerRole, _ = s.roleRepo.FindByName(ctx, "customer")
	}

	if customerRole != nil {
		_ = s.roleRepo.AssignToUser(ctx, user.ID, customerRole.ID)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*TokenResponse, error) {
	user, err := s.userRepo.FindByEmailWithRoles(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := s.password.Verify(user.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	accessToken, err := s.jwtService.GenerateToken(user.ID, user.Email, user.FullName, roleNames)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshHash, err := s.password.Hash(refreshToken)
	if err != nil {
		return nil, err
	}

	user.RefreshTokenHash = refreshHash
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(72 * time.Hour).Unix(),
		UserID:       user.ID,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	var foundUser *models.User

	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	for i := range users {
		if users[i].RefreshTokenHash == "" {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(users[i].RefreshTokenHash), []byte(refreshToken)); err == nil {
			foundUser = &users[i]
			break
		}
	}

	if foundUser == nil {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(ctx, foundUser.ID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	accessToken, err := s.jwtService.GenerateToken(user.ID, user.Email, user.FullName, roleNames)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshHash, err := s.password.Hash(newRefreshToken)
	if err != nil {
		return nil, err
	}

	user.RefreshTokenHash = refreshHash
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(72 * time.Hour).Unix(),
		UserID:       user.ID,
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := s.password.Verify(user.PasswordHash, oldPassword); err != nil {
		return ErrInvalidCredentials
	}

	newHash, err := s.password.Hash(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	return s.userRepo.Update(ctx, user)
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
