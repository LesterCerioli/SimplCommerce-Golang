package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ActivityService struct {
	db *sql.DB
}

func NewActivityService(db *sql.DB) *ActivityService {
	return &ActivityService{db: db}
}

type ActivityResponse struct {
	ID             uint      `json:"id"`
	ActivityTypeID uint      `json:"activityTypeId"`
	ActivityType   string    `json:"activityType"`
	UserID         uint      `json:"userId"`
	EntityID       uint      `json:"entityId"`
	EntityTypeID   string    `json:"entityTypeId"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (s *ActivityService) LogActivity(ctx context.Context, activityTypeName string, userID, entityID uint, entityTypeID string) error {
	var activityTypeID uint
	err := s.db.QueryRowContext(ctx, `SELECT id FROM activity_log_activity_types WHERE name = $1`, activityTypeName).Scan(&activityTypeID)
	if err == sql.ErrNoRows {
		err = s.db.QueryRowContext(ctx, `
			INSERT INTO activity_log_activity_types (name, created_at, updated_at) VALUES ($1, NOW(), NOW())
			RETURNING id
		`, activityTypeName).Scan(&activityTypeID)
		if err != nil {
			return fmt.Errorf("failed to create activity type: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to find activity type: %w", err)
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO activity_log_activities (activity_type_id, user_id, entity_id, entity_type_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,NOW(),NOW())
	`, activityTypeID, userID, entityID, entityTypeID)
	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}
	return nil
}

func (s *ActivityService) GetActivities(ctx context.Context, page, pageSize int) ([]ActivityResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM activity_log_activities`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count activities: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT a.id, a.activity_type_id, COALESCE(at.name,''), a.user_id, a.entity_id, COALESCE(a.entity_type_id,''), a.created_at
		FROM activity_log_activities a
		LEFT JOIN activity_log_activity_types at ON at.id = a.activity_type_id
		ORDER BY a.created_at DESC LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list activities: %w", err)
	}
	defer rows.Close()

	var activities []ActivityResponse
	for rows.Next() {
		var a ActivityResponse
		if err := rows.Scan(&a.ID, &a.ActivityTypeID, &a.ActivityType, &a.UserID, &a.EntityID, &a.EntityTypeID, &a.CreatedAt); err != nil {
			return nil, 0, err
		}
		activities = append(activities, a)
	}
	if activities == nil {
		activities = []ActivityResponse{}
	}
	return activities, total, nil
}
