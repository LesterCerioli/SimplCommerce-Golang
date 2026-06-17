package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/pkg/repository"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
)

type CategoryRepository interface {
	repository.Repository[catalogmodels.Category]
	FindBySlug(ctx context.Context, slug string) (*catalogmodels.Category, error)
	GetTree(ctx context.Context) ([]catalogmodels.Category, error)
	GetByParent(ctx context.Context, parentID *uint) ([]catalogmodels.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, entity *catalogmodels.Category) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *categoryRepository) Update(ctx context.Context, entity *catalogmodels.Category) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalogmodels.Category{}, id).Error
}

func (r *categoryRepository) FindByID(ctx context.Context, id uint) (*catalogmodels.Category, error) {
	var entity catalogmodels.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]catalogmodels.Category, error) {
	var entities []catalogmodels.Category
	err := r.db.WithContext(ctx).
		Order("display_order ASC").
		Find(&entities).Error
	return entities, err
}

func (r *categoryRepository) FindWhere(ctx context.Context, condition map[string]interface{}) ([]catalogmodels.Category, error) {
	var entities []catalogmodels.Category
	err := r.db.WithContext(ctx).Where(condition).Order("display_order ASC").Find(&entities).Error
	return entities, err
}

func (r *categoryRepository) FirstWhere(ctx context.Context, condition map[string]interface{}) (*catalogmodels.Category, error) {
	var entity catalogmodels.Category
	err := r.db.WithContext(ctx).Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *categoryRepository) Paginate(ctx context.Context, page, pageSize int) ([]catalogmodels.Category, int64, error) {
	var entities []catalogmodels.Category
	var total int64

	query := r.db.WithContext(ctx).Model(&catalogmodels.Category{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("display_order ASC").
		Offset(offset).Limit(pageSize).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*catalogmodels.Category, error) {
	var entity catalogmodels.Category
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("slug = ?", slug).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *categoryRepository) GetTree(ctx context.Context) ([]catalogmodels.Category, error) {
	var entities []catalogmodels.Category
	err := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Order("display_order ASC").
		Find(&entities).Error
	return entities, err
}

func (r *categoryRepository) GetByParent(ctx context.Context, parentID *uint) ([]catalogmodels.Category, error) {
	var entities []catalogmodels.Category
	query := r.db.WithContext(ctx).Order("display_order ASC")
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	err := query.Find(&entities).Error
	return entities, err
}
