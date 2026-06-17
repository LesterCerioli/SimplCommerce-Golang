package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/activitylog/internal/models"
)

type ActivityRepository interface {
	Create(ctx context.Context, activity *models.Activity) error
	FindAll(ctx context.Context) ([]models.Activity, error)
	FindMostViewed(ctx context.Context, limit int) ([]models.Activity, error)
	FindRecent(ctx context.Context, limit int) ([]models.Activity, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Activity, int64, error)
}

type GormActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &GormActivityRepository{db: db}
}

func (r *GormActivityRepository) Create(ctx context.Context, activity *models.Activity) error {
	return r.db.WithContext(ctx).Create(activity).Error
}

func (r *GormActivityRepository) FindAll(ctx context.Context) ([]models.Activity, error) {
	var activities []models.Activity
	err := r.db.WithContext(ctx).Preload("ActivityType").Order("created_at desc").Find(&activities).Error
	return activities, err
}

func (r *GormActivityRepository) FindMostViewed(ctx context.Context, limit int) ([]models.Activity, error) {
	var activities []models.Activity
	err := r.db.WithContext(ctx).
		Preload("ActivityType").
		Joins("JOIN activity_log_activity_types ON activity_log_activity_types.id = activities.activity_type_id").
		Where("activity_log_activity_types.name = ?", "EntityView").
		Group("activities.entity_id, activities.entity_type_id, activities.id").
		Order("count(activities.id) desc").
		Limit(limit).
		Find(&activities).Error
	return activities, err
}

func (r *GormActivityRepository) FindRecent(ctx context.Context, limit int) ([]models.Activity, error) {
	var activities []models.Activity
	err := r.db.WithContext(ctx).
		Preload("ActivityType").
		Order("created_at desc").
		Limit(limit).
		Find(&activities).Error
	return activities, err
}

func (r *GormActivityRepository) Paginate(ctx context.Context, page, pageSize int) ([]models.Activity, int64, error) {
	var activities []models.Activity
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Activity{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("ActivityType").Offset(offset).Limit(pageSize).Order("created_at desc").Find(&activities).Error
	if err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}
