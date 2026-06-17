package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/identity/internal/models"
)

type CountryRepository interface {
	FindAll(ctx context.Context) ([]models.Country, error)
	FindByID(ctx context.Context, id string) (*models.Country, error)
}

type GormCountryRepository struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) CountryRepository {
	return &GormCountryRepository{db: db}
}

func (r *GormCountryRepository) FindAll(ctx context.Context) ([]models.Country, error) {
	var countries []models.Country
	err := r.db.WithContext(ctx).Order("name ASC").Find(&countries).Error
	return countries, err
}

func (r *GormCountryRepository) FindByID(ctx context.Context, id string) (*models.Country, error) {
	var country models.Country
	err := r.db.WithContext(ctx).First(&country, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &country, nil
}
