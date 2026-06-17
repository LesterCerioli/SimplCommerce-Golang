package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/notifications/internal/models"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	Update(ctx context.Context, notification *models.Notification) error
	FindByID(ctx context.Context, id uint) (*models.Notification, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.Notification, error)
	UnreadCount(ctx context.Context, userID uint) (int64, error)
	MarkAsRead(ctx context.Context, id uint) error
}

type GormNotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &GormNotificationRepository{db: db}
}

func (r *GormNotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *GormNotificationRepository) Update(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Save(notification).Error
}

func (r *GormNotificationRepository) FindByID(ctx context.Context, id uint) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.WithContext(ctx).First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *GormNotificationRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *GormNotificationRepository) UnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *GormNotificationRepository) MarkAsRead(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}
