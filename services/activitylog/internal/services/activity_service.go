package services

import (
	"context"

	"github.com/simplcommerce-go/services/activitylog/internal/models"
	"github.com/simplcommerce-go/services/activitylog/internal/repositories"
)

type ActivityService interface {
	LogActivity(ctx context.Context, activityTypeName string, userID, entityID uint, entityTypeID string) error
	GetRecentActivities(ctx context.Context, limit int) ([]models.Activity, error)
	GetAllActivities(ctx context.Context) ([]models.Activity, error)
	GetMostViewed(ctx context.Context, limit int) ([]models.Activity, error)
}

type activityService struct {
	activityRepo     repositories.ActivityRepository
	activityTypeRepo repositories.ActivityTypeRepository
}

func NewActivityService(
	activityRepo repositories.ActivityRepository,
	activityTypeRepo repositories.ActivityTypeRepository,
) ActivityService {
	return &activityService{
		activityRepo:     activityRepo,
		activityTypeRepo: activityTypeRepo,
	}
}

func (s *activityService) LogActivity(ctx context.Context, activityTypeName string, userID, entityID uint, entityTypeID string) error {
	activityType, err := s.activityTypeRepo.FirstOrCreate(ctx, activityTypeName)
	if err != nil {
		return err
	}

	activity := &models.Activity{
		ActivityTypeID: activityType.ID,
		UserID:         userID,
		EntityID:       entityID,
		EntityTypeID:   entityTypeID,
	}

	return s.activityRepo.Create(ctx, activity)
}

func (s *activityService) GetRecentActivities(ctx context.Context, limit int) ([]models.Activity, error) {
	return s.activityRepo.FindRecent(ctx, limit)
}

func (s *activityService) GetAllActivities(ctx context.Context) ([]models.Activity, error) {
	return s.activityRepo.FindAll(ctx)
}

func (s *activityService) GetMostViewed(ctx context.Context, limit int) ([]models.Activity, error) {
	return s.activityRepo.FindMostViewed(ctx, limit)
}
