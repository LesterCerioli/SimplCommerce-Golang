package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/simplcommerce-go/pkg/config"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   uint
	Email    string
	FullName string
	Roles    []string
}

type JWTService interface {
	GenerateToken(userID uint, email, fullName string, roles []string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type jwtService struct {
	config config.JWTConfig
}

func NewJWTService(cfg config.JWTConfig) JWTService {
	return &jwtService{config: cfg}
}

func (s *jwtService) GenerateToken(userID uint, email, fullName string, roles []string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.Expiration) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:   userID,
		Email:    email,
		FullName: fullName,
		Roles:    roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
