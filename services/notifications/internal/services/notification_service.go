package services

import (
	"context"
	"errors"

	"github.com/simplcommerce-go/services/notifications/internal/models"
	"github.com/simplcommerce-go/services/notifications/internal/repositories"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

type NotificationService interface {
	Create(ctx context.Context, notification *models.Notification) error
	FindByID(ctx context.Context, id uint) (*models.Notification, error)
	MarkAsRead(ctx context.Context, id uint) error
	GetUserNotifications(ctx context.Context, userID uint) ([]models.Notification, error)
	UnreadCount(ctx context.Context, userID uint) (int64, error)
}

type notificationService struct {
	notificationRepo repositories.NotificationRepository
}

func NewNotificationService(notificationRepo repositories.NotificationRepository) NotificationService {
	return &notificationService{notificationRepo: notificationRepo}
}

func (s *notificationService) Create(ctx context.Context, notification *models.Notification) error {
	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationService) FindByID(ctx context.Context, id uint) (*models.Notification, error) {
	notification, err := s.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrNotificationNotFound
	}
	return notification, nil
}

func (s *notificationService) MarkAsRead(ctx context.Context, id uint) error {
	_, err := s.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return ErrNotificationNotFound
	}
	return s.notificationRepo.MarkAsRead(ctx, id)
}

func (s *notificationService) GetUserNotifications(ctx context.Context, userID uint) ([]models.Notification, error) {
	return s.notificationRepo.FindByUserID(ctx, userID)
}

func (s *notificationService) UnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.notificationRepo.UnreadCount(ctx, userID)
}
