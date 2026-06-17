package services

import (
	"context"

	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
)

type VendorService interface {
	Create(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error)
	Update(ctx context.Context, id uint, vendor *models.Vendor) (*models.Vendor, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Vendor, error)
	FindAll(ctx context.Context) ([]models.Vendor, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Vendor, int64, error)
}

type vendorService struct {
	vendorRepo repositories.VendorRepository
}

func NewVendorService(vendorRepo repositories.VendorRepository) VendorService {
	return &vendorService{vendorRepo: vendorRepo}
}

func (s *vendorService) Create(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	if err := s.vendorRepo.Create(ctx, vendor); err != nil {
		return nil, err
	}
	return vendor, nil
}

func (s *vendorService) Update(ctx context.Context, id uint, vendor *models.Vendor) (*models.Vendor, error) {
	existing, err := s.vendorRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Name = vendor.Name
	existing.Slug = vendor.Slug
	existing.Description = vendor.Description
	existing.Email = vendor.Email
	existing.IsActive = vendor.IsActive

	if err := s.vendorRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *vendorService) Delete(ctx context.Context, id uint) error {
	return s.vendorRepo.Delete(ctx, id)
}

func (s *vendorService) FindByID(ctx context.Context, id uint) (*models.Vendor, error) {
	return s.vendorRepo.FindByID(ctx, id)
}

func (s *vendorService) FindAll(ctx context.Context) ([]models.Vendor, error) {
	return s.vendorRepo.FindAll(ctx)
}

func (s *vendorService) Paginate(ctx context.Context, page, pageSize int) ([]models.Vendor, int64, error) {
	return s.vendorRepo.Paginate(ctx, page, pageSize)
}
