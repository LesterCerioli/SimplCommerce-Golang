package services

import (
	"context"

	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
)

type AddressService interface {
	Create(ctx context.Context, userID uint, address *models.Address) (*models.Address, error)
	Update(ctx context.Context, id uint, address *models.Address) (*models.Address, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Address, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.Address, error)
}

type addressService struct {
	addressRepo   repositories.AddressRepository
	userAddrRepo  repositories.UserAddressRepository
}

func NewAddressService(addressRepo repositories.AddressRepository, userAddrRepo repositories.UserAddressRepository) AddressService {
	return &addressService{
		addressRepo:  addressRepo,
		userAddrRepo: userAddrRepo,
	}
}

func (s *addressService) Create(ctx context.Context, userID uint, address *models.Address) (*models.Address, error) {
	if err := s.addressRepo.Create(ctx, address); err != nil {
		return nil, err
	}

	userAddress := &models.UserAddress{
		UserID:    userID,
		AddressID: address.ID,
	}
	if err := s.userAddrRepo.Create(ctx, userAddress); err != nil {
		return nil, err
	}

	return address, nil
}

func (s *addressService) Update(ctx context.Context, id uint, address *models.Address) (*models.Address, error) {
	existing, err := s.addressRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.ContactName = address.ContactName
	existing.Phone = address.Phone
	existing.AddressLine1 = address.AddressLine1
	existing.AddressLine2 = address.AddressLine2
	existing.City = address.City
	existing.ZipCode = address.ZipCode
	existing.DistrictID = address.DistrictID
	existing.StateOrProvinceID = address.StateOrProvinceID
	existing.CountryID = address.CountryID

	if err := s.addressRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *addressService) Delete(ctx context.Context, id uint) error {
	return s.addressRepo.Delete(ctx, id)
}

func (s *addressService) FindByID(ctx context.Context, id uint) (*models.Address, error) {
	return s.addressRepo.FindByID(ctx, id)
}

func (s *addressService) FindByUserID(ctx context.Context, userID uint) ([]models.Address, error) {
	return s.addressRepo.FindByUserID(ctx, userID)
}
