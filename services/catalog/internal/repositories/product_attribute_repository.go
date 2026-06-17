package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/pkg/repository"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
)

type ProductAttributeRepository interface {
	repository.Repository[catalogmodels.ProductAttribute]
}

type productAttributeRepository struct {
	db *gorm.DB
}

func NewProductAttributeRepository(db *gorm.DB) ProductAttributeRepository {
	return &productAttributeRepository{db: db}
}

func (r *productAttributeRepository) Create(ctx context.Context, entity *catalogmodels.ProductAttribute) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *productAttributeRepository) Update(ctx context.Context, entity *catalogmodels.ProductAttribute) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *productAttributeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalogmodels.ProductAttribute{}, id).Error
}

func (r *productAttributeRepository) FindByID(ctx context.Context, id uint) (*catalogmodels.ProductAttribute, error) {
	var entity catalogmodels.ProductAttribute
	err := r.db.WithContext(ctx).Preload("Group").First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productAttributeRepository) FindAll(ctx context.Context) ([]catalogmodels.ProductAttribute, error) {
	var entities []catalogmodels.ProductAttribute
	err := r.db.WithContext(ctx).Preload("Group").Find(&entities).Error
	return entities, err
}

func (r *productAttributeRepository) FindWhere(ctx context.Context, condition map[string]interface{}) ([]catalogmodels.ProductAttribute, error) {
	var entities []catalogmodels.ProductAttribute
	err := r.db.WithContext(ctx).Preload("Group").Where(condition).Find(&entities).Error
	return entities, err
}

func (r *productAttributeRepository) FirstWhere(ctx context.Context, condition map[string]interface{}) (*catalogmodels.ProductAttribute, error) {
	var entity catalogmodels.ProductAttribute
	err := r.db.WithContext(ctx).Preload("Group").Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productAttributeRepository) Paginate(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductAttribute, int64, error) {
	var entities []catalogmodels.ProductAttribute
	var total int64

	query := r.db.WithContext(ctx).Model(&catalogmodels.ProductAttribute{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Group").
		Offset(offset).Limit(pageSize).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
