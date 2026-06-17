package services

import (
	"context"
	"errors"

	"github.com/simplcommerce-go/services/reviews/internal/models"
	"github.com/simplcommerce-go/services/reviews/internal/repositories"
)

var (
	ErrReviewNotFound = errors.New("review not found")
)

type ReviewService interface {
	Create(ctx context.Context, review *models.Review) error
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Review, error)
	GetProductReviews(ctx context.Context, productID uint) ([]models.Review, error)
	GetEntityReviews(ctx context.Context, entityTypeID string, entityID uint) ([]models.Review, error)
	Approve(ctx context.Context, id uint) error
	Reject(ctx context.Context, id uint) error
}

type reviewService struct {
	reviewRepo repositories.ReviewRepository
}

func NewReviewService(reviewRepo repositories.ReviewRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo}
}

func (s *reviewService) Create(ctx context.Context, review *models.Review) error {
	return s.reviewRepo.Create(ctx, review)
}

func (s *reviewService) Update(ctx context.Context, review *models.Review) error {
	return s.reviewRepo.Update(ctx, review)
}

func (s *reviewService) Delete(ctx context.Context, id uint) error {
	return s.reviewRepo.Delete(ctx, id)
}

func (s *reviewService) FindByID(ctx context.Context, id uint) (*models.Review, error) {
	review, err := s.reviewRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrReviewNotFound
	}
	return review, nil
}

func (s *reviewService) GetProductReviews(ctx context.Context, productID uint) ([]models.Review, error) {
	return s.reviewRepo.FindByEntity(ctx, "Product", productID)
}

func (s *reviewService) GetEntityReviews(ctx context.Context, entityTypeID string, entityID uint) ([]models.Review, error) {
	return s.reviewRepo.FindByEntity(ctx, entityTypeID, entityID)
}

func (s *reviewService) Approve(ctx context.Context, id uint) error {
	review, err := s.reviewRepo.FindByID(ctx, id)
	if err != nil {
		return ErrReviewNotFound
	}
	review.Status = "Approved"
	return s.reviewRepo.Update(ctx, review)
}

func (s *reviewService) Reject(ctx context.Context, id uint) error {
	review, err := s.reviewRepo.FindByID(ctx, id)
	if err != nil {
		return ErrReviewNotFound
	}
	review.Status = "Rejected"
	return s.reviewRepo.Update(ctx, review)
}
