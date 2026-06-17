package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/cms/internal/models"
)

type PageRepository interface {
	Create(ctx context.Context, page *models.Page) error
	Update(ctx context.Context, page *models.Page) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Page, error)
	FindBySlug(ctx context.Context, slug string) (*models.Page, error)
	FindAll(ctx context.Context) ([]models.Page, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Page, int64, error)
}

type GormPageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
	return &GormPageRepository{db: db}
}

func (r *GormPageRepository) Create(ctx context.Context, page *models.Page) error {
	return r.db.WithContext(ctx).Create(page).Error
}

func (r *GormPageRepository) Update(ctx context.Context, page *models.Page) error {
	return r.db.WithContext(ctx).Save(page).Error
}

func (r *GormPageRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Page{}, id).Error
}

func (r *GormPageRepository) FindByID(ctx context.Context, id uint) (*models.Page, error) {
	var page models.Page
	err := r.db.WithContext(ctx).First(&page, id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *GormPageRepository) FindBySlug(ctx context.Context, slug string) (*models.Page, error) {
	var page models.Page
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *GormPageRepository) FindAll(ctx context.Context) ([]models.Page, error) {
	var pages []models.Page
	err := r.db.WithContext(ctx).Find(&pages).Error
	return pages, err
}

func (r *GormPageRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.Page, int64, error) {
	var pages []models.Page
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Page{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&pages).Error
	if err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}
