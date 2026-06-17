package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/identity/internal/models"
)

type VendorRepository interface {
	Create(ctx context.Context, vendor *models.Vendor) error
	Update(ctx context.Context, vendor *models.Vendor) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Vendor, error)
	FindAll(ctx context.Context) ([]models.Vendor, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Vendor, int64, error)
}

type GormVendorRepository struct {
	db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) VendorRepository {
	return &GormVendorRepository{db: db}
}

func (r *GormVendorRepository) Create(ctx context.Context, vendor *models.Vendor) error {
	return r.db.WithContext(ctx).Create(vendor).Error
}

func (r *GormVendorRepository) Update(ctx context.Context, vendor *models.Vendor) error {
	return r.db.WithContext(ctx).Save(vendor).Error
}

func (r *GormVendorRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Vendor{}, id).Error
}

func (r *GormVendorRepository) FindByID(ctx context.Context, id uint) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.WithContext(ctx).Preload("Users").First(&vendor, id).Error
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

func (r *GormVendorRepository) FindAll(ctx context.Context) ([]models.Vendor, error) {
	var vendors []models.Vendor
	err := r.db.WithContext(ctx).Find(&vendors).Error
	return vendors, err
}

func (r *GormVendorRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.Vendor, int64, error) {
	var vendors []models.Vendor
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Vendor{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&vendors).Error
	if err != nil {
		return nil, 0, err
	}

	return vendors, total, nil
}
