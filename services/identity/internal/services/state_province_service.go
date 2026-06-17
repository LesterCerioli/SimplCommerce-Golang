package services

import (
	"context"

	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
)

type StateOrProvinceService interface {
	FindByCountryID(ctx context.Context, countryID string) ([]models.StateOrProvince, error)
	FindByID(ctx context.Context, id uint) (*models.StateOrProvince, error)
}

type stateOrProvinceService struct {
	stateRepo repositories.StateOrProvinceRepository
}

func NewStateOrProvinceService(stateRepo repositories.StateOrProvinceRepository) StateOrProvinceService {
	return &stateOrProvinceService{stateRepo: stateRepo}
}

func (s *stateOrProvinceService) FindByCountryID(ctx context.Context, countryID string) ([]models.StateOrProvince, error) {
	return s.stateRepo.FindByCountryID(ctx, countryID)
}

func (s *stateOrProvinceService) FindByID(ctx context.Context, id uint) (*models.StateOrProvince, error) {
	return s.stateRepo.FindByID(ctx, id)
}
