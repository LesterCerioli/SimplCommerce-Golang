package services

import (
	"context"

	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
)

type DistrictService interface {
	FindByStateOrProvinceID(ctx context.Context, stateID uint) ([]models.District, error)
	FindByID(ctx context.Context, id uint) (*models.District, error)
}

type districtService struct {
	districtRepo repositories.DistrictRepository
}

func NewDistrictService(districtRepo repositories.DistrictRepository) DistrictService {
	return &districtService{districtRepo: districtRepo}
}

func (s *districtService) FindByStateOrProvinceID(ctx context.Context, stateID uint) ([]models.District, error) {
	return s.districtRepo.FindByStateOrProvinceID(ctx, stateID)
}

func (s *districtService) FindByID(ctx context.Context, id uint) (*models.District, error) {
	return s.districtRepo.FindByID(ctx, id)
}
