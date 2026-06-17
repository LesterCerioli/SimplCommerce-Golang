package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type NotificationService struct {
	db *sql.DB
}

func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{db: db}
}

type NotificationResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	IsRead     bool      `json:"isRead"`
	EntityID   *string   `json:"entityId,omitempty"`
	EntityType string    `json:"entityType"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (s *NotificationService) Create(ctx context.Context, userID string, title, body, entityType string, entityID *string) (*NotificationResponse, error) {
	var n NotificationResponse
	var eID sql.NullString
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO notifications_notifications (user_id, title, body, is_read, entity_id, entity_type, created_at, updated_at)
		VALUES ($1,$2,$3,false,$4,$5,NOW(),NOW())
		RETURNING id, user_id, title, body, is_read, entity_id, entity_type, created_at, updated_at
	`, userID, title, body, entityID, entityType).Scan(
		&n.ID, &n.UserID, &n.Title, &n.Body, &n.IsRead, &eID, &n.EntityType, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	if eID.Valid {
		n.EntityID = &eID.String
	}
	return &n, nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID string) ([]NotificationResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, title, body, is_read, entity_id, COALESCE(entity_type,''), created_at, updated_at
		FROM notifications_notifications WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []NotificationResponse
	for rows.Next() {
		var n NotificationResponse
		var entityID sql.NullString
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.IsRead, &entityID, &n.EntityType, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		if entityID.Valid {
			n.EntityID = &entityID.String
		}
		notifications = append(notifications, n)
	}
	if notifications == nil {
		notifications = []NotificationResponse{}
	}
	return notifications, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id, userID string) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE notifications_notifications SET is_read = true, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("notification not found")
	}
	return nil
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM notifications_notifications WHERE user_id = $1 AND is_read = false
	`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", err)
	}
	return count, nil
}
