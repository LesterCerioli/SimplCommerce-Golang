package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/identity/internal/models"
)

type AddressRepository interface {
	Create(ctx context.Context, address *models.Address) error
	Update(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Address, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.Address, error)
	FindAll(ctx context.Context) ([]models.Address, error)
}

type GormAddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &GormAddressRepository{db: db}
}

func (r *GormAddressRepository) Create(ctx context.Context, address *models.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *GormAddressRepository) Update(ctx context.Context, address *models.Address) error {
	return r.db.WithContext(ctx).Save(address).Error
}

func (r *GormAddressRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Address{}, id).Error
}

func (r *GormAddressRepository) FindByID(ctx context.Context, id uint) (*models.Address, error) {
	var address models.Address
	err := r.db.WithContext(ctx).Preload("StateOrProvince").Preload("Country").Preload("District").First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *GormAddressRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.WithContext(ctx).
		Joins("JOIN identity_user_addresses ON identity_user_addresses.address_id = identity_addresses.id").
		Where("identity_user_addresses.user_id = ?", userID).
		Preload("StateOrProvince").Preload("Country").Preload("District").
		Find(&addresses).Error
	return addresses, err
}

func (r *GormAddressRepository) FindAll(ctx context.Context) ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.WithContext(ctx).Preload("StateOrProvince").Preload("Country").Preload("District").Find(&addresses).Error
	return addresses, err
}
