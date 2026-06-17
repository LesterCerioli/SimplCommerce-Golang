package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/reviews/internal/models"
)

type ReplyRepository interface {
	Create(ctx context.Context, reply *models.Reply) error
	Update(ctx context.Context, reply *models.Reply) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Reply, error)
	FindByReviewID(ctx context.Context, reviewID uint) ([]models.Reply, error)
}

type GormReplyRepository struct {
	db *gorm.DB
}

func NewReplyRepository(db *gorm.DB) ReplyRepository {
	return &GormReplyRepository{db: db}
}

func (r *GormReplyRepository) Create(ctx context.Context, reply *models.Reply) error {
	return r.db.WithContext(ctx).Create(reply).Error
}

func (r *GormReplyRepository) Update(ctx context.Context, reply *models.Reply) error {
	return r.db.WithContext(ctx).Save(reply).Error
}

func (r *GormReplyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Reply{}, id).Error
}

func (r *GormReplyRepository) FindByID(ctx context.Context, id uint) (*models.Reply, error) {
	var reply models.Reply
	err := r.db.WithContext(ctx).Preload("Review").First(&reply, id).Error
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func (r *GormReplyRepository) FindByReviewID(ctx context.Context, reviewID uint) ([]models.Reply, error) {
	var replies []models.Reply
	err := r.db.WithContext(ctx).Where("review_id = ?", reviewID).Find(&replies).Error
	return replies, err
}
