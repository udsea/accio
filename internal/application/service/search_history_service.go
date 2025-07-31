package service

import (
	"context"

	"github.com/accio/internal/domain/model"
	"github.com/accio/internal/domain/repository"
)

// SearchHistoryService defines the interface for search history operations
type SearchHistoryService interface {
	// RecordSearch records a search query
	RecordSearch(ctx context.Context, query string, resultCount int) error

	// GetPopularSearches gets the most popular searches
	GetPopularSearches(ctx context.Context, limit int) ([]*model.SearchHistory, error)
}

// SearchHistoryServiceImpl implements the SearchHistoryService interface
type SearchHistoryServiceImpl struct {
	searchHistoryRepo repository.SearchHistoryRepository
}

// NewSearchHistoryService creates a new SearchHistoryServiceImpl
func NewSearchHistoryService(searchHistoryRepo repository.SearchHistoryRepository) SearchHistoryService {
	return &SearchHistoryServiceImpl{
		searchHistoryRepo: searchHistoryRepo,
	}
}

// RecordSearch records a search query
func (s *SearchHistoryServiceImpl) RecordSearch(ctx context.Context, query string, resultCount int) error {
	searchHistory := model.NewSearchHistory(query, resultCount)
	return s.searchHistoryRepo.Create(ctx, searchHistory)
}

// GetPopularSearches gets the most popular searches
func (s *SearchHistoryServiceImpl) GetPopularSearches(ctx context.Context, limit int) ([]*model.SearchHistory, error) {
	return s.searchHistoryRepo.FindPopular(ctx, limit)
}
