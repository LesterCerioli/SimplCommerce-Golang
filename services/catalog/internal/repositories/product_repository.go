package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/pkg/repository"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
)

type ProductRepository interface {
	repository.Repository[catalogmodels.Product]
	FindBySlug(ctx context.Context, slug string) (*catalogmodels.Product, error)
	Search(ctx context.Context, query string, minPrice, maxPrice *float64, categoryID *uint, page, pageSize int) ([]catalogmodels.Product, int64, error)
	GetByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]catalogmodels.Product, int64, error)
	GetFeatured(ctx context.Context, count int) ([]catalogmodels.Product, error)
	GetRelated(ctx context.Context, productID uint, count int) ([]catalogmodels.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, entity *catalogmodels.Product) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *productRepository) Update(ctx context.Context, entity *catalogmodels.Product) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalogmodels.Product{}, id).Error
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*catalogmodels.Product, error) {
	var entity catalogmodels.Product
	err := r.db.WithContext(ctx).
		Preload("Categories").
		Preload("Brand").
		Preload("AttributeValues.Attribute").
		Preload("OptionValues.Option").
		Preload("OptionCombinations.Option").
		Preload("Medias.Media").
		Preload("ThumbnailImage").
		First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]catalogmodels.Product, error) {
	var entities []catalogmodels.Product
	err := r.db.WithContext(ctx).
		Preload("Categories").
		Preload("Brand").
		Find(&entities).Error
	return entities, err
}

func (r *productRepository) FindWhere(ctx context.Context, condition map[string]interface{}) ([]catalogmodels.Product, error) {
	var entities []catalogmodels.Product
	err := r.db.WithContext(ctx).
		Preload("Categories").
		Preload("Brand").
		Where(condition).Find(&entities).Error
	return entities, err
}

func (r *productRepository) FirstWhere(ctx context.Context, condition map[string]interface{}) (*catalogmodels.Product, error) {
	var entity catalogmodels.Product
	err := r.db.WithContext(ctx).
		Preload("Categories").
		Preload("Brand").
		Preload("AttributeValues.Attribute").
		Preload("OptionValues.Option").
		Preload("OptionCombinations.Option").
		Preload("Medias.Media").
		Preload("ThumbnailImage").
		Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productRepository) Paginate(ctx context.Context, page, pageSize int) ([]catalogmodels.Product, int64, error) {
	var entities []catalogmodels.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&catalogmodels.Product{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Categories").
		Preload("Brand").
		Offset(offset).Limit(pageSize).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *productRepository) FindBySlug(ctx context.Context, slug string) (*catalogmodels.Product, error) {
	var entity catalogmodels.Product
	err := r.db.WithContext(ctx).
		Preload("Categories").
		Preload("Brand").
		Preload("AttributeValues.Attribute").
		Preload("AttributeValues.Attribute.Group").
		Preload("OptionValues.Option").
		Preload("OptionCombinations.Option").
		Preload("Medias.Media").
		Preload("ThumbnailImage").
		Preload("LinkedProducts.LinkedProduct").
		Where("slug = ?", slug).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productRepository) Search(ctx context.Context, query string, minPrice, maxPrice *float64, categoryID *uint, page, pageSize int) ([]catalogmodels.Product, int64, error) {
	var entities []catalogmodels.Product
	var total int64

	db := r.db.WithContext(ctx).Model(&catalogmodels.Product{})

	if query != "" {
		db = db.Where("name ILIKE ? OR short_description ILIKE ? OR sku ILIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}
	if minPrice != nil {
		db = db.Where("price >= ?", *minPrice)
	}
	if maxPrice != nil {
		db = db.Where("price <= ?", *maxPrice)
	}
	if categoryID != nil {
		db = db.Joins("JOIN catalog_product_categories ON catalog_product_categories.product_id = catalog_products.id").
			Where("catalog_product_categories.category_id = ?", *categoryID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.
		Preload("Categories").
		Preload("Brand").
		Offset(offset).Limit(pageSize).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *productRepository) GetByCategory(ctx context.Context, categoryID uint, page, pageSize int) ([]catalogmodels.Product, int64, error) {
	var entities []catalogmodels.Product
	var total int64

	query := r.db.WithContext(ctx).
		Model(&catalogmodels.Product{}).
		Joins("JOIN catalog_product_categories ON catalog_product_categories.product_id = catalog_products.id").
		Where("catalog_product_categories.category_id = ?", categoryID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("Categories").
		Preload("Brand").
		Offset(offset).Limit(pageSize).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *productRepository) GetFeatured(ctx context.Context, count int) ([]catalogmodels.Product, error) {
	var entities []catalogmodels.Product
	err := r.db.WithContext(ctx).
		Where("is_featured = ? AND is_published = ?", true, true).
		Preload("Categories").
		Preload("Brand").
		Limit(count).
		Find(&entities).Error
	return entities, err
}

func (r *productRepository) GetRelated(ctx context.Context, productID uint, count int) ([]catalogmodels.Product, error) {
	var entities []catalogmodels.Product

	var linkedIDs []uint
	r.db.WithContext(ctx).
		Model(&catalogmodels.ProductLink{}).
		Where("product_id = ?", productID).
		Pluck("linked_product_id", &linkedIDs)

	if len(linkedIDs) == 0 {
		return entities, nil
	}

	err := r.db.WithContext(ctx).
		Where("id IN ?", linkedIDs).
		Preload("Categories").
		Preload("Brand").
		Limit(count).
		Find(&entities).Error
	return entities, err
}
