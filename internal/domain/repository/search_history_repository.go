package repository

import (
	"context"

	"github.com/accio/internal/domain/model"
)

// SearchHistoryRepository defines the interface for search history data access
type SearchHistoryRepository interface {
	// Create creates a new search history entry
	Create(ctx context.Context, searchHistory *model.SearchHistory) error

	// FindByQuery finds search history entries by query
	FindByQuery(ctx context.Context, query string) ([]*model.SearchHistory, error)

	// FindPopular finds the most popular searches
	FindPopular(ctx context.Context, limit int) ([]*model.SearchHistory, error)

	// Count counts all search history entries
	Count(ctx context.Context) (int64, error)
}
