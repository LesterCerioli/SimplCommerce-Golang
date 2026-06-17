package services

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/search/internal/models"
)

type SearchRepository interface {
	Create(ctx context.Context, query *models.Query) error
	FindPopular(ctx context.Context, limit int) ([]models.Query, error)
}

type GormSearchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) SearchRepository {
	return &GormSearchRepository{db: db}
}

func (r *GormSearchRepository) Create(ctx context.Context, query *models.Query) error {
	return r.db.WithContext(ctx).Create(query).Error
}

func (r *GormSearchRepository) FindPopular(ctx context.Context, limit int) ([]models.Query, error) {
	var queries []models.Query
	err := r.db.WithContext(ctx).
		Order("results_count desc").
		Limit(limit).
		Find(&queries).Error
	return queries, err
}

type SearchService interface {
	LogSearch(ctx context.Context, queryText string, resultsCount int) error
	GetPopularSearches(ctx context.Context, limit int) ([]models.Query, error)
}

type searchService struct {
	searchRepo SearchRepository
}

func NewSearchService(searchRepo SearchRepository) SearchService {
	return &searchService{searchRepo: searchRepo}
}

func (s *searchService) LogSearch(ctx context.Context, queryText string, resultsCount int) error {
	query := &models.Query{
		QueryText:    queryText,
		ResultsCount: resultsCount,
	}
	return s.searchRepo.Create(ctx, query)
}

func (s *searchService) GetPopularSearches(ctx context.Context, limit int) ([]models.Query, error) {
	return s.searchRepo.FindPopular(ctx, limit)
}
