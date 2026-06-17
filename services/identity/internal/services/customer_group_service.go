package services

import (
	"context"

	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
)

type CustomerGroupService interface {
	Create(ctx context.Context, group *models.CustomerGroup) (*models.CustomerGroup, error)
	Update(ctx context.Context, id uint, group *models.CustomerGroup) (*models.CustomerGroup, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.CustomerGroup, error)
	FindAll(ctx context.Context) ([]models.CustomerGroup, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.CustomerGroup, int64, error)
}

type customerGroupService struct {
	groupRepo repositories.CustomerGroupRepository
}

func NewCustomerGroupService(groupRepo repositories.CustomerGroupRepository) CustomerGroupService {
	return &customerGroupService{groupRepo: groupRepo}
}

func (s *customerGroupService) Create(ctx context.Context, group *models.CustomerGroup) (*models.CustomerGroup, error) {
	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *customerGroupService) Update(ctx context.Context, id uint, group *models.CustomerGroup) (*models.CustomerGroup, error) {
	existing, err := s.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Name = group.Name
	existing.Description = group.Description
	existing.IsActive = group.IsActive

	if err := s.groupRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *customerGroupService) Delete(ctx context.Context, id uint) error {
	return s.groupRepo.Delete(ctx, id)
}

func (s *customerGroupService) FindByID(ctx context.Context, id uint) (*models.CustomerGroup, error) {
	return s.groupRepo.FindByID(ctx, id)
}

func (s *customerGroupService) FindAll(ctx context.Context) ([]models.CustomerGroup, error) {
	return s.groupRepo.FindAll(ctx)
}

func (s *customerGroupService) Paginate(ctx context.Context, page, pageSize int) ([]models.CustomerGroup, int64, error) {
	return s.groupRepo.Paginate(ctx, page, pageSize)
}
