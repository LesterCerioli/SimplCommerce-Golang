package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/identity/internal/models"
)

type MediaRepository interface {
	Create(ctx context.Context, media *models.Media) error
	Update(ctx context.Context, media *models.Media) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Media, error)
	FindAll(ctx context.Context) ([]models.Media, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Media, int64, error)
}

type GormMediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &GormMediaRepository{db: db}
}

func (r *GormMediaRepository) Create(ctx context.Context, media *models.Media) error {
	return r.db.WithContext(ctx).Create(media).Error
}

func (r *GormMediaRepository) Update(ctx context.Context, media *models.Media) error {
	return r.db.WithContext(ctx).Save(media).Error
}

func (r *GormMediaRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Media{}, id).Error
}

func (r *GormMediaRepository) FindByID(ctx context.Context, id uint) (*models.Media, error) {
	var media models.Media
	err := r.db.WithContext(ctx).First(&media, id).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *GormMediaRepository) FindAll(ctx context.Context) ([]models.Media, error) {
	var medias []models.Media
	err := r.db.WithContext(ctx).Find(&medias).Error
	return medias, err
}

func (r *GormMediaRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.Media, int64, error) {
	var medias []models.Media
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Media{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&medias).Error
	if err != nil {
		return nil, 0, err
	}

	return medias, total, nil
}
