package services

import (
	"context"

	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
)

type UserService interface {
	Create(ctx context.Context, req *RegisterRequest) (*models.User, error)
	Update(ctx context.Context, id uint, req *RegisterRequest) (*models.User, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindAll(ctx context.Context) ([]models.User, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.User, int64, error)
}

type userService struct {
	userRepo repositories.UserRepository
	password PasswordService
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
		password: NewPasswordService(),
	}
}

func (s *userService) Create(ctx context.Context, req *RegisterRequest) (*models.User, error) {
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

	return user, nil
}

func (s *userService) Update(ctx context.Context, id uint, req *RegisterRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = req.FullName
	user.Email = req.Email

	if req.Password != "" {
		passwordHash, err := s.password.Hash(req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = passwordHash
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) FindByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) FindAll(ctx context.Context) ([]models.User, error) {
	return s.userRepo.FindAll(ctx)
}

func (s *userService) Paginate(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.Paginate(ctx, page, pageSize)
}
