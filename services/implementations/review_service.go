package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type ReviewService struct {
	db *sql.DB
}

func NewReviewService(db *sql.DB) *ReviewService {
	return &ReviewService{db: db}
}

type ReviewResponse struct {
	ID           uint               `json:"id"`
	UserID       uint               `json:"userId"`
	Title        string             `json:"title"`
	Comment      string             `json:"comment"`
	Rating       int                `json:"rating"`
	ReviewerName string             `json:"reviewerName"`
	Status       string             `json:"status"`
	EntityTypeID string             `json:"entityTypeId"`
	EntityID     uint               `json:"entityId"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	Replies      []ReplyResponse    `json:"replies,omitempty"`
}

type ReplyResponse struct {
	ID          uint      `json:"id"`
	ReviewID    uint      `json:"reviewId"`
	UserID      uint      `json:"userId"`
	Comment     string    `json:"comment"`
	ReplierName string    `json:"replierName"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (s *ReviewService) GetProductReviews(ctx context.Context, productID uint, page, pageSize int) ([]ReviewResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM reviews_reviews WHERE entity_type_id = 'Product' AND entity_id = $1 AND status = 'Approved'
	`, productID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, title, comment, rating, reviewer_name, status, entity_type_id, entity_id, created_at, updated_at
		FROM reviews_reviews WHERE entity_type_id = 'Product' AND entity_id = $1 AND status = 'Approved'
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, productID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list reviews: %w", err)
	}
	defer rows.Close()

	var reviews []ReviewResponse
	for rows.Next() {
		var r ReviewResponse
		if err := rows.Scan(&r.ID, &r.UserID, &r.Title, &r.Comment, &r.Rating, &r.ReviewerName, &r.Status, &r.EntityTypeID, &r.EntityID, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, 0, err
		}
		replies, _ := s.getReplies(ctx, r.ID)
		r.Replies = replies
		reviews = append(reviews, r)
	}
	if reviews == nil {
		reviews = []ReviewResponse{}
	}
	return reviews, total, nil
}

func (s *ReviewService) CreateReview(ctx context.Context, userID uint, title, comment string, rating int, reviewerName, entityTypeID string, entityID uint) (*ReviewResponse, error) {
	var r ReviewResponse
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO reviews_reviews (user_id, title, comment, rating, reviewer_name, status, entity_type_id, entity_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,'Pending',$6,$7,NOW(),NOW())
		RETURNING id, user_id, title, comment, rating, reviewer_name, status, entity_type_id, entity_id, created_at, updated_at
	`, userID, title, comment, rating, reviewerName, entityTypeID, entityID).Scan(
		&r.ID, &r.UserID, &r.Title, &r.Comment, &r.Rating, &r.ReviewerName, &r.Status, &r.EntityTypeID, &r.EntityID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}
	r.Replies = []ReplyResponse{}
	return &r, nil
}

func (s *ReviewService) UpdateStatus(ctx context.Context, id uint, status string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE reviews_reviews SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	if err != nil {
		return fmt.Errorf("failed to update review status: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("review not found")
	}
	return nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, id uint) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM reviews_reviews WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("review not found")
	}
	return nil
}

func (s *ReviewService) AddReply(ctx context.Context, reviewID, userID uint, comment, replierName string) (*ReplyResponse, error) {
	_, err := s.db.ExecContext(ctx, `SELECT 1 FROM reviews_reviews WHERE id = $1`, reviewID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("review not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find review: %w", err)
	}

	var reply ReplyResponse
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO reviews_replies (review_id, user_id, comment, replier_name, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,'Approved',NOW(),NOW())
		RETURNING id, review_id, user_id, comment, replier_name, status, created_at
	`, reviewID, userID, comment, replierName).Scan(
		&reply.ID, &reply.ReviewID, &reply.UserID, &reply.Comment, &reply.ReplierName, &reply.Status, &reply.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add reply: %w", err)
	}
	return &reply, nil
}

func (s *ReviewService) getReplies(ctx context.Context, reviewID uint) ([]ReplyResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, review_id, user_id, comment, replier_name, status, created_at
		FROM reviews_replies WHERE review_id = $1 ORDER BY created_at ASC
	`, reviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []ReplyResponse
	for rows.Next() {
		var r ReplyResponse
		if err := rows.Scan(&r.ID, &r.ReviewID, &r.UserID, &r.Comment, &r.ReplierName, &r.Status, &r.CreatedAt); err != nil {
			return nil, err
		}
		replies = append(replies, r)
	}
	if replies == nil {
		replies = []ReplyResponse{}
	}
	return replies, nil
}
