package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/pkg/repository"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
)

type BrandRepository interface {
	repository.Repository[catalogmodels.Brand]
	FindBySlug(ctx context.Context, slug string) (*catalogmodels.Brand, error)
}

type brandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &brandRepository{db: db}
}

func (r *brandRepository) Create(ctx context.Context, entity *catalogmodels.Brand) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *brandRepository) Update(ctx context.Context, entity *catalogmodels.Brand) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *brandRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalogmodels.Brand{}, id).Error
}

func (r *brandRepository) FindByID(ctx context.Context, id uint) (*catalogmodels.Brand, error) {
	var entity catalogmodels.Brand
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *brandRepository) FindAll(ctx context.Context) ([]catalogmodels.Brand, error) {
	var entities []catalogmodels.Brand
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *brandRepository) FindWhere(ctx context.Context, condition map[string]interface{}) ([]catalogmodels.Brand, error) {
	var entities []catalogmodels.Brand
	err := r.db.WithContext(ctx).Where(condition).Find(&entities).Error
	return entities, err
}

func (r *brandRepository) FirstWhere(ctx context.Context, condition map[string]interface{}) (*catalogmodels.Brand, error) {
	var entity catalogmodels.Brand
	err := r.db.WithContext(ctx).Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *brandRepository) Paginate(ctx context.Context, page, pageSize int) ([]catalogmodels.Brand, int64, error) {
	var entities []catalogmodels.Brand
	var total int64

	query := r.db.WithContext(ctx).Model(&catalogmodels.Brand{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *brandRepository) FindBySlug(ctx context.Context, slug string) (*catalogmodels.Brand, error) {
	var entity catalogmodels.Brand
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
