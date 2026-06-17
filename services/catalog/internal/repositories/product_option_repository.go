package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/pkg/repository"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
)

type ProductOptionRepository interface {
	repository.Repository[catalogmodels.ProductOption]
}

type productOptionRepository struct {
	db *gorm.DB
}

func NewProductOptionRepository(db *gorm.DB) ProductOptionRepository {
	return &productOptionRepository{db: db}
}

func (r *productOptionRepository) Create(ctx context.Context, entity *catalogmodels.ProductOption) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *productOptionRepository) Update(ctx context.Context, entity *catalogmodels.ProductOption) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *productOptionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalogmodels.ProductOption{}, id).Error
}

func (r *productOptionRepository) FindByID(ctx context.Context, id uint) (*catalogmodels.ProductOption, error) {
	var entity catalogmodels.ProductOption
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productOptionRepository) FindAll(ctx context.Context) ([]catalogmodels.ProductOption, error) {
	var entities []catalogmodels.ProductOption
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *productOptionRepository) FindWhere(ctx context.Context, condition map[string]interface{}) ([]catalogmodels.ProductOption, error) {
	var entities []catalogmodels.ProductOption
	err := r.db.WithContext(ctx).Where(condition).Find(&entities).Error
	return entities, err
}

func (r *productOptionRepository) FirstWhere(ctx context.Context, condition map[string]interface{}) (*catalogmodels.ProductOption, error) {
	var entity catalogmodels.ProductOption
	err := r.db.WithContext(ctx).Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productOptionRepository) Paginate(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductOption, int64, error) {
	var entities []catalogmodels.ProductOption
	var total int64

	query := r.db.WithContext(ctx).Model(&catalogmodels.ProductOption{})
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
