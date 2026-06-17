package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/pkg/repository"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
)

type ProductTemplateRepository interface {
	repository.Repository[catalogmodels.ProductTemplate]
}

type productTemplateRepository struct {
	db *gorm.DB
}

func NewProductTemplateRepository(db *gorm.DB) ProductTemplateRepository {
	return &productTemplateRepository{db: db}
}

func (r *productTemplateRepository) Create(ctx context.Context, entity *catalogmodels.ProductTemplate) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *productTemplateRepository) Update(ctx context.Context, entity *catalogmodels.ProductTemplate) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *productTemplateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalogmodels.ProductTemplate{}, id).Error
}

func (r *productTemplateRepository) FindByID(ctx context.Context, id uint) (*catalogmodels.ProductTemplate, error) {
	var entity catalogmodels.ProductTemplate
	err := r.db.WithContext(ctx).Preload("Attributes").First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productTemplateRepository) FindAll(ctx context.Context) ([]catalogmodels.ProductTemplate, error) {
	var entities []catalogmodels.ProductTemplate
	err := r.db.WithContext(ctx).Preload("Attributes").Find(&entities).Error
	return entities, err
}

func (r *productTemplateRepository) FindWhere(ctx context.Context, condition map[string]interface{}) ([]catalogmodels.ProductTemplate, error) {
	var entities []catalogmodels.ProductTemplate
	err := r.db.WithContext(ctx).Preload("Attributes").Where(condition).Find(&entities).Error
	return entities, err
}

func (r *productTemplateRepository) FirstWhere(ctx context.Context, condition map[string]interface{}) (*catalogmodels.ProductTemplate, error) {
	var entity catalogmodels.ProductTemplate
	err := r.db.WithContext(ctx).Preload("Attributes").Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productTemplateRepository) Paginate(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductTemplate, int64, error) {
	var entities []catalogmodels.ProductTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&catalogmodels.ProductTemplate{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Attributes").
		Offset(offset).Limit(pageSize).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
