package services

import (
	"context"

	"github.com/simplcommerce-go/services/reviews/internal/models"
	"github.com/simplcommerce-go/services/reviews/internal/repositories"
)

type ReplyService interface {
	Create(ctx context.Context, reply *models.Reply) error
	Update(ctx context.Context, reply *models.Reply) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Reply, error)
	FindByReviewID(ctx context.Context, reviewID uint) ([]models.Reply, error)
}

type replyService struct {
	replyRepo repositories.ReplyRepository
}

func NewReplyService(replyRepo repositories.ReplyRepository) ReplyService {
	return &replyService{replyRepo: replyRepo}
}

func (s *replyService) Create(ctx context.Context, reply *models.Reply) error {
	return s.replyRepo.Create(ctx, reply)
}

func (s *replyService) Update(ctx context.Context, reply *models.Reply) error {
	return s.replyRepo.Update(ctx, reply)
}

func (s *replyService) Delete(ctx context.Context, id uint) error {
	return s.replyRepo.Delete(ctx, id)
}

func (s *replyService) FindByID(ctx context.Context, id uint) (*models.Reply, error) {
	return s.replyRepo.FindByID(ctx, id)
}

func (s *replyService) FindByReviewID(ctx context.Context, reviewID uint) ([]models.Reply, error) {
	return s.replyRepo.FindByReviewID(ctx, reviewID)
}
