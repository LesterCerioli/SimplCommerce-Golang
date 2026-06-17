package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/identity/internal/models"
)

type DistrictRepository interface {
	FindByStateOrProvinceID(ctx context.Context, stateID uint) ([]models.District, error)
	FindByID(ctx context.Context, id uint) (*models.District, error)
}

type GormDistrictRepository struct {
	db *gorm.DB
}

func NewDistrictRepository(db *gorm.DB) DistrictRepository {
	return &GormDistrictRepository{db: db}
}

func (r *GormDistrictRepository) FindByStateOrProvinceID(ctx context.Context, stateID uint) ([]models.District, error) {
	var districts []models.District
	err := r.db.WithContext(ctx).Where("state_or_province_id = ?", stateID).Preload("StateOrProvince").Order("name ASC").Find(&districts).Error
	return districts, err
}

func (r *GormDistrictRepository) FindByID(ctx context.Context, id uint) (*models.District, error) {
	var district models.District
	err := r.db.WithContext(ctx).Preload("StateOrProvince").First(&district, id).Error
	if err != nil {
		return nil, err
	}
	return &district, nil
}
