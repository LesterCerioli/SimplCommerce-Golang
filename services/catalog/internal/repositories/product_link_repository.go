package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/pkg/repository"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
)

type ProductLinkRepository interface {
	repository.Repository[catalogmodels.ProductLink]
}

type productLinkRepository struct {
	db *gorm.DB
}

func NewProductLinkRepository(db *gorm.DB) ProductLinkRepository {
	return &productLinkRepository{db: db}
}

func (r *productLinkRepository) Create(ctx context.Context, entity *catalogmodels.ProductLink) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *productLinkRepository) Update(ctx context.Context, entity *catalogmodels.ProductLink) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *productLinkRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalogmodels.ProductLink{}, id).Error
}

func (r *productLinkRepository) FindByID(ctx context.Context, id uint) (*catalogmodels.ProductLink, error) {
	var entity catalogmodels.ProductLink
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("LinkedProduct").
		First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productLinkRepository) FindAll(ctx context.Context) ([]catalogmodels.ProductLink, error) {
	var entities []catalogmodels.ProductLink
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("LinkedProduct").
		Find(&entities).Error
	return entities, err
}

func (r *productLinkRepository) FindWhere(ctx context.Context, condition map[string]interface{}) ([]catalogmodels.ProductLink, error) {
	var entities []catalogmodels.ProductLink
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("LinkedProduct").
		Where(condition).Find(&entities).Error
	return entities, err
}

func (r *productLinkRepository) FirstWhere(ctx context.Context, condition map[string]interface{}) (*catalogmodels.ProductLink, error) {
	var entity catalogmodels.ProductLink
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("LinkedProduct").
		Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productLinkRepository) Paginate(ctx context.Context, page, pageSize int) ([]catalogmodels.ProductLink, int64, error) {
	var entities []catalogmodels.ProductLink
	var total int64

	query := r.db.WithContext(ctx).Model(&catalogmodels.ProductLink{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Product").
		Preload("LinkedProduct").
		Offset(offset).Limit(pageSize).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
