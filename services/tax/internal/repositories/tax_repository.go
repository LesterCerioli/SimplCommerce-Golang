package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/tax/internal/models"
)

type TaxClassRepository interface {
	Create(ctx context.Context, entity *models.TaxClass) error
	Update(ctx context.Context, entity *models.TaxClass) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.TaxClass, error)
	FindAll(ctx context.Context) ([]models.TaxClass, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.TaxClass, int64, error)
}

type taxClassRepository struct {
	db *gorm.DB
}

func NewTaxClassRepository(db *gorm.DB) TaxClassRepository {
	return &taxClassRepository{db: db}
}

func (r *taxClassRepository) Create(ctx context.Context, entity *models.TaxClass) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *taxClassRepository) Update(ctx context.Context, entity *models.TaxClass) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *taxClassRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TaxClass{}, id).Error
}

func (r *taxClassRepository) FindByID(ctx context.Context, id uint) (*models.TaxClass, error) {
	var entity models.TaxClass
	err := r.db.WithContext(ctx).Preload("Rates").First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *taxClassRepository) FindAll(ctx context.Context) ([]models.TaxClass, error) {
	var entities []models.TaxClass
	err := r.db.WithContext(ctx).Preload("Rates").Find(&entities).Error
	return entities, err
}

func (r *taxClassRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.TaxClass, int64, error) {
	var entities []models.TaxClass
	var total int64

	query := r.db.WithContext(ctx).Model(&models.TaxClass{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Rates").Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

type TaxRateRepository interface {
	Create(ctx context.Context, entity *models.TaxRate) error
	Update(ctx context.Context, entity *models.TaxRate) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.TaxRate, error)
	FindAll(ctx context.Context) ([]models.TaxRate, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.TaxRate, int64, error)
	FindByTaxClassID(ctx context.Context, taxClassID uint) ([]models.TaxRate, error)
	FindByLocation(ctx context.Context, countryID string, stateOrProvinceID *uint, zipCode string) ([]models.TaxRate, error)
}

type taxRateRepository struct {
	db *gorm.DB
}

func NewTaxRateRepository(db *gorm.DB) TaxRateRepository {
	return &taxRateRepository{db: db}
}

func (r *taxRateRepository) Create(ctx context.Context, entity *models.TaxRate) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *taxRateRepository) Update(ctx context.Context, entity *models.TaxRate) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *taxRateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TaxRate{}, id).Error
}

func (r *taxRateRepository) FindByID(ctx context.Context, id uint) (*models.TaxRate, error) {
	var entity models.TaxRate
	err := r.db.WithContext(ctx).Preload("TaxClass").First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *taxRateRepository) FindAll(ctx context.Context) ([]models.TaxRate, error) {
	var entities []models.TaxRate
	err := r.db.WithContext(ctx).Preload("TaxClass").Find(&entities).Error
	return entities, err
}

func (r *taxRateRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.TaxRate, int64, error) {
	var entities []models.TaxRate
	var total int64

	query := r.db.WithContext(ctx).Model(&models.TaxRate{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("TaxClass").Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *taxRateRepository) FindByTaxClassID(ctx context.Context, taxClassID uint) ([]models.TaxRate, error) {
	var entities []models.TaxRate
	err := r.db.WithContext(ctx).Where("tax_class_id = ?", taxClassID).Find(&entities).Error
	return entities, err
}

func (r *taxRateRepository) FindByLocation(ctx context.Context, countryID string, stateOrProvinceID *uint, zipCode string) ([]models.TaxRate, error) {
	var entities []models.TaxRate

	query := r.db.WithContext(ctx).Where("country_id = ?", countryID)

	if stateOrProvinceID != nil {
		query = query.Where("(state_or_province_id IS NULL OR state_or_province_id = ?)", *stateOrProvinceID)
	} else {
		query = query.Where("state_or_province_id IS NULL")
	}

	if zipCode != "" {
		query = query.Where("(zip_code = '' OR zip_code = ?)", zipCode)
	} else {
		query = query.Where("zip_code = ''")
	}

	err := query.Find(&entities).Error
	return entities, err
}
