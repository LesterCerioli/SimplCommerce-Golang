package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SearchService struct {
	db *sql.DB
}

func NewSearchService(db *sql.DB) *SearchService {
	return &SearchService{db: db}
}

type SearchQueryResponse struct {
	ID           string    `json:"id"`
	QueryText    string    `json:"queryText"`
	ResultsCount int       `json:"resultsCount"`
	CreatedAt    time.Time `json:"createdAt"`
}

type PopularSearchTerm struct {
	QueryText string `json:"queryText"`
	Count     int64  `json:"count"`
}

func (s *SearchService) LogSearch(ctx context.Context, queryText string, resultsCount int) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO search_queries (query_text, results_count, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`, queryText, resultsCount)
	if err != nil {
		return fmt.Errorf("failed to log search: %w", err)
	}
	return nil
}

func (s *SearchService) GetPopularSearches(ctx context.Context, limit int) ([]PopularSearchTerm, error) {
	if limit < 1 {
		limit = 10
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT query_text, COUNT(*) as cnt
		FROM search_queries
		GROUP BY query_text
		ORDER BY cnt DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular searches: %w", err)
	}
	defer rows.Close()

	var terms []PopularSearchTerm
	for rows.Next() {
		var t PopularSearchTerm
		if err := rows.Scan(&t.QueryText, &t.Count); err != nil {
			return nil, err
		}
		terms = append(terms, t)
	}
	if terms == nil {
		terms = []PopularSearchTerm{}
	}
	return terms, nil
}
