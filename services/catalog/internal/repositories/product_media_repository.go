package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/pkg/repository"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
)

type ProductMediaRepository interface {
	repository.Repository[catalogmodels.ProductMedia]
}

type productMediaRepository struct {
	db *gorm.DB
}

func NewProductMediaRepository(db *gorm.DB) ProductMediaRepository {
	return &productMediaRepository{db: db}
}

func (r *productMediaRepository) Create(ctx context.Context, entity *catalogmodels.ProductMedia) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *productMediaRepository) Update(ctx context.Context, entity *catalogmodels.ProductMedia) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *productMediaRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalogmodels.ProductMedia{}, id).Error
}

func (r *productMediaRepository) FindByID(ctx context.Context, id uint) (*catalogmodels.ProductMedia, error) {
	var entity catalogmodels.ProductMedia
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Media").
		First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productMediaRepository) FindAll(ctx context.Context) ([]catalogmodels.ProductMedia, error) {
	var entities []catalogmodels.ProductMedia
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Media").
		Order("display_order ASC").
		Find(&entities).Error
	return entities, err
}

func (r *productMediaRepository) FindWhere(ctx context.Context, condition map[string]interface{}) ([]catalogmodels.ProductMedia, error) {
	var entities []catalogmodels.ProductMedia
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Media").
		Where(condition).
		Order("display_order ASC").
		Find(&entities).Error
	return entities, err
}

func (r *productMediaRepository) FirstWhere(ctx context.Context, condition map[string]interface{}) (*catalogmodels.ProductMedia, error) {
	var entity catalogmodels.ProductMedia
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Media").
		Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productMediaRepository) Paginate(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductMedia, int64, error) {
	var entities []catalogmodels.ProductMedia
	var total int64

	query := r.db.WithContext(ctx).Model(&catalogmodels.ProductMedia{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Product").
		Preload("Media").
		Order("display_order ASC").
		Offset(offset).Limit(pageSize).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
