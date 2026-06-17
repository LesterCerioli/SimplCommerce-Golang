package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/activitylog/internal/models"
)

type ActivityTypeRepository interface {
	Create(ctx context.Context, activityType *models.ActivityType) error
	FindByName(ctx context.Context, name string) (*models.ActivityType, error)
	FirstOrCreate(ctx context.Context, name string) (*models.ActivityType, error)
}

type GormActivityTypeRepository struct {
	db *gorm.DB
}

func NewActivityTypeRepository(db *gorm.DB) ActivityTypeRepository {
	return &GormActivityTypeRepository{db: db}
}

func (r *GormActivityTypeRepository) Create(ctx context.Context, activityType *models.ActivityType) error {
	return r.db.WithContext(ctx).Create(activityType).Error
}

func (r *GormActivityTypeRepository) FindByName(ctx context.Context, name string) (*models.ActivityType, error) {
	var activityType models.ActivityType
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&activityType).Error
	if err != nil {
		return nil, err
	}
	return &activityType, nil
}

func (r *GormActivityTypeRepository) FirstOrCreate(ctx context.Context, name string) (*models.ActivityType, error) {
	var activityType models.ActivityType
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&activityType).Error
	if err != nil {
		activityType = models.ActivityType{Name: name}
		if createErr := r.db.WithContext(ctx).Create(&activityType).Error; createErr != nil {
			return nil, createErr
		}
	}
	return &activityType, nil
}
