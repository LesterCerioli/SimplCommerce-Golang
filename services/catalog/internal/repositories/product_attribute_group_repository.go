package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/pkg/repository"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
)

type ProductAttributeGroupRepository interface {
	repository.Repository[catalogmodels.ProductAttributeGroup]
}

type productAttributeGroupRepository struct {
	db *gorm.DB
}

func NewProductAttributeGroupRepository(db *gorm.DB) ProductAttributeGroupRepository {
	return &productAttributeGroupRepository{db: db}
}

func (r *productAttributeGroupRepository) Create(ctx context.Context, entity *catalogmodels.ProductAttributeGroup) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *productAttributeGroupRepository) Update(ctx context.Context, entity *catalogmodels.ProductAttributeGroup) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *productAttributeGroupRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalogmodels.ProductAttributeGroup{}, id).Error
}

func (r *productAttributeGroupRepository) FindByID(ctx context.Context, id uint) (*catalogmodels.ProductAttributeGroup, error) {
	var entity catalogmodels.ProductAttributeGroup
	err := r.db.WithContext(ctx).Preload("Attributes").First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productAttributeGroupRepository) FindAll(ctx context.Context) ([]catalogmodels.ProductAttributeGroup, error) {
	var entities []catalogmodels.ProductAttributeGroup
	err := r.db.WithContext(ctx).Preload("Attributes").Find(&entities).Error
	return entities, err
}

func (r *productAttributeGroupRepository) FindWhere(ctx context.Context, condition map[string]interface{}) ([]catalogmodels.ProductAttributeGroup, error) {
	var entities []catalogmodels.ProductAttributeGroup
	err := r.db.WithContext(ctx).Preload("Attributes").Where(condition).Find(&entities).Error
	return entities, err
}

func (r *productAttributeGroupRepository) FirstWhere(ctx context.Context, condition map[string]interface{}) (*catalogmodels.ProductAttributeGroup, error) {
	var entity catalogmodels.ProductAttributeGroup
	err := r.db.WithContext(ctx).Preload("Attributes").Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productAttributeGroupRepository) Paginate(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductAttributeGroup, int64, error) {
	var entities []catalogmodels.ProductAttributeGroup
	var total int64

	query := r.db.WithContext(ctx).Model(&catalogmodels.ProductAttributeGroup{})
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
