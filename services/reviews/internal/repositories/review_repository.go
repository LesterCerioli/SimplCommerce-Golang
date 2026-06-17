package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/reviews/internal/models"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Review, error)
	FindByEntity(ctx context.Context, entityTypeID string, entityID uint) ([]models.Review, error)
	FindAll(ctx context.Context) ([]models.Review, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Review, int64, error)
}

type GormReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &GormReviewRepository{db: db}
}

func (r *GormReviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *GormReviewRepository) Update(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *GormReviewRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, id).Error
}

func (r *GormReviewRepository) FindByID(ctx context.Context, id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).Preload("Replies").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *GormReviewRepository) FindByEntity(ctx context.Context, entityTypeID string, entityID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.WithContext(ctx).
		Where("entity_type_id = ? AND entity_id = ?", entityTypeID, entityID).
		Preload("Replies").
		Find(&reviews).Error
	return reviews, err
}

func (r *GormReviewRepository) FindAll(ctx context.Context) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.WithContext(ctx).Preload("Replies").Find(&reviews).Error
	return reviews, err
}

func (r *GormReviewRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Review{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Replies").Offset(offset).Limit(pageSize).Find(&reviews).Error
	if err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}
