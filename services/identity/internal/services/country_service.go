package services

import (
	"context"

	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
)

type CountryService interface {
	FindAll(ctx context.Context) ([]models.Country, error)
	FindByID(ctx context.Context, id string) (*models.Country, error)
}

type countryService struct {
	countryRepo repositories.CountryRepository
}

func NewCountryService(countryRepo repositories.CountryRepository) CountryService {
	return &countryService{countryRepo: countryRepo}
}

func (s *countryService) FindAll(ctx context.Context) ([]models.Country, error) {
	return s.countryRepo.FindAll(ctx)
}

func (s *countryService) FindByID(ctx context.Context, id string) (*models.Country, error) {
	return s.countryRepo.FindByID(ctx, id)
}
