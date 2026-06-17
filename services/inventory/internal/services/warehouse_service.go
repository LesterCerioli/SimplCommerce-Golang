package services

import (
	"context"

	"github.com/simplcommerce-go/services/inventory/internal/models"
	"github.com/simplcommerce-go/services/inventory/internal/repositories"
)

type WarehouseService interface {
	Create(ctx context.Context, warehouse *models.Warehouse) error
	Update(ctx context.Context, warehouse *models.Warehouse) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Warehouse, error)
	FindAll(ctx context.Context) ([]models.Warehouse, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Warehouse, int64, error)
}

type warehouseService struct {
	repo repositories.WarehouseRepository
}

func NewWarehouseService(repo repositories.WarehouseRepository) WarehouseService {
	return &warehouseService{repo: repo}
}

func (s *warehouseService) Create(ctx context.Context, warehouse *models.Warehouse) error {
	return s.repo.Create(ctx, warehouse)
}

func (s *warehouseService) Update(ctx context.Context, warehouse *models.Warehouse) error {
	return s.repo.Update(ctx, warehouse)
}

func (s *warehouseService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *warehouseService) FindByID(ctx context.Context, id uint) (*models.Warehouse, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *warehouseService) FindAll(ctx context.Context) ([]models.Warehouse, error) {
	return s.repo.FindAll(ctx)
}

func (s *warehouseService) Paginate(ctx context.Context, page, pageSize int) ([]models.Warehouse, int64, error) {
	return s.repo.Paginate(ctx, page, pageSize)
}
