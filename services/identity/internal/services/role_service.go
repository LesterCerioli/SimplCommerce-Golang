package services

import (
	"context"
	"strings"

	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
)

type RoleService interface {
	Create(ctx context.Context, name string) (*models.Role, error)
	FindAll(ctx context.Context) ([]models.Role, error)
	AssignToUser(ctx context.Context, userID, roleID uint) error
	RemoveFromUser(ctx context.Context, userID, roleID uint) error
}

type roleService struct {
	roleRepo repositories.RoleRepository
}

func NewRoleService(roleRepo repositories.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

func (s *roleService) Create(ctx context.Context, name string) (*models.Role, error) {
	role := &models.Role{
		Name:           name,
		NormalizedName: strings.ToUpper(name),
		ConcurrencyStamp: generateUUID(),
	}
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *roleService) FindAll(ctx context.Context) ([]models.Role, error) {
	return s.roleRepo.FindAll(ctx)
}

func (s *roleService) AssignToUser(ctx context.Context, userID, roleID uint) error {
	return s.roleRepo.AssignToUser(ctx, userID, roleID)
}

func (s *roleService) RemoveFromUser(ctx context.Context, userID, roleID uint) error {
	return s.roleRepo.RemoveFromUser(ctx, userID, roleID)
}
