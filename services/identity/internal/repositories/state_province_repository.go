package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/identity/internal/models"
)

type StateOrProvinceRepository interface {
	FindByCountryID(ctx context.Context, countryID string) ([]models.StateOrProvince, error)
	FindByID(ctx context.Context, id uint) (*models.StateOrProvince, error)
}

type GormStateOrProvinceRepository struct {
	db *gorm.DB
}

func NewStateOrProvinceRepository(db *gorm.DB) StateOrProvinceRepository {
	return &GormStateOrProvinceRepository{db: db}
}

func (r *GormStateOrProvinceRepository) FindByCountryID(ctx context.Context, countryID string) ([]models.StateOrProvince, error) {
	var states []models.StateOrProvince
	err := r.db.WithContext(ctx).Where("country_id = ?", countryID).Preload("Country").Order("name ASC").Find(&states).Error
	return states, err
}

func (r *GormStateOrProvinceRepository) FindByID(ctx context.Context, id uint) (*models.StateOrProvince, error) {
	var state models.StateOrProvince
	err := r.db.WithContext(ctx).Preload("Country").First(&state, id).Error
	if err != nil {
		return nil, err
	}
	return &state, nil
}
