package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/identity/internal/models"
)

type UserAddressRepository interface {
	Create(ctx context.Context, userAddress *models.UserAddress) error
	Delete(ctx context.Context, id uint) error
	FindByUserID(ctx context.Context, userID uint) ([]models.UserAddress, error)
}

type GormUserAddressRepository struct {
	db *gorm.DB
}

func NewUserAddressRepository(db *gorm.DB) UserAddressRepository {
	return &GormUserAddressRepository{db: db}
}

func (r *GormUserAddressRepository) Create(ctx context.Context, userAddress *models.UserAddress) error {
	return r.db.WithContext(ctx).Create(userAddress).Error
}

func (r *GormUserAddressRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.UserAddress{}, id).Error
}

func (r *GormUserAddressRepository) FindByUserID(ctx context.Context, userID uint) ([]models.UserAddress, error) {
	var userAddresses []models.UserAddress
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Address").Preload("Address.StateOrProvince").Preload("Address.Country").Preload("Address.District").Find(&userAddresses).Error
	return userAddresses, err
}
