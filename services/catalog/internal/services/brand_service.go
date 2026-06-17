package services

import (
	"context"

	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/repositories"
)

type BrandService interface {
	CreateBrand(ctx context.Context, brand *catalogmodels.Brand) error
	UpdateBrand(ctx context.Context, brand *catalogmodels.Brand) error
	DeleteBrand(ctx context.Context, id uint) error
	GetBrandByID(ctx context.Context, id uint) (*catalogmodels.Brand, error)
	GetBrandBySlug(ctx context.Context, slug string) (*catalogmodels.Brand, error)
	GetBrands(ctx context.Context, page, pageSize int) ([]catalogmodels.Brand, int64, error)
}

type brandService struct {
	brandRepo repositories.BrandRepository
}

func NewBrandService(brandRepo repositories.BrandRepository) BrandService {
	return &brandService{brandRepo: brandRepo}
}

func (s *brandService) CreateBrand(ctx context.Context, brand *catalogmodels.Brand) error {
	if brand.Slug == "" {
		brand.Slug = generateSlug(brand.Name)
	}
	return s.brandRepo.Create(ctx, brand)
}

func (s *brandService) UpdateBrand(ctx context.Context, brand *catalogmodels.Brand) error {
	if brand.Slug == "" {
		brand.Slug = generateSlug(brand.Name)
	}
	return s.brandRepo.Update(ctx, brand)
}

func (s *brandService) DeleteBrand(ctx context.Context, id uint) error {
	return s.brandRepo.Delete(ctx, id)
}

func (s *brandService) GetBrandByID(ctx context.Context, id uint) (*catalogmodels.Brand, error) {
	return s.brandRepo.FindByID(ctx, id)
}

func (s *brandService) GetBrandBySlug(ctx context.Context, slug string) (*catalogmodels.Brand, error) {
	return s.brandRepo.FindBySlug(ctx, slug)
}

func (s *brandService) GetBrands(ctx context.Context, page, pageSize int) ([]catalogmodels.Brand, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.brandRepo.Paginate(ctx, page, pageSize)
}
