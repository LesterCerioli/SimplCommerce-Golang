package services

import (
	"context"

	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/repositories"
)

type MediaService interface {
	Create(ctx context.Context, media *models.Media) (*models.Media, error)
	Update(ctx context.Context, id uint, media *models.Media) (*models.Media, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Media, error)
	FindAll(ctx context.Context) ([]models.Media, error)
	Paginate(ctx context.Context, page, pageSize int) ([]models.Media, int64, error)
}

type mediaService struct {
	mediaRepo repositories.MediaRepository
}

func NewMediaService(mediaRepo repositories.MediaRepository) MediaService {
	return &mediaService{mediaRepo: mediaRepo}
}

func (s *mediaService) Create(ctx context.Context, media *models.Media) (*models.Media, error) {
	if err := s.mediaRepo.Create(ctx, media); err != nil {
		return nil, err
	}
	return media, nil
}

func (s *mediaService) Update(ctx context.Context, id uint, media *models.Media) (*models.Media, error) {
	existing, err := s.mediaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Caption = media.Caption
	existing.FileSize = media.FileSize
	existing.FileName = media.FileName
	existing.MediaType = media.MediaType

	if err := s.mediaRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *mediaService) Delete(ctx context.Context, id uint) error {
	return s.mediaRepo.Delete(ctx, id)
}

func (s *mediaService) FindByID(ctx context.Context, id uint) (*models.Media, error) {
	return s.mediaRepo.FindByID(ctx, id)
}

func (s *mediaService) FindAll(ctx context.Context) ([]models.Media, error) {
	return s.mediaRepo.FindAll(ctx)
}

func (s *mediaService) Paginate(ctx context.Context, page, pageSize int) ([]models.Media, int64, error) {
	return s.mediaRepo.Paginate(ctx, page, pageSize)
}
