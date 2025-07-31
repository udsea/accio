package persistence

import (
	"context"

	"github.com/accio/internal/domain/model"
	"github.com/accio/internal/domain/repository"
	"gorm.io/gorm"
)

// GormSearchHistoryRepository is a GORM implementation of SearchHistoryRepository
type GormSearchHistoryRepository struct {
	db *gorm.DB
}

// NewGormSearchHistoryRepository creates a new GormSearchHistoryRepository
func NewGormSearchHistoryRepository(db *gorm.DB) repository.SearchHistoryRepository {
	return &GormSearchHistoryRepository{
		db: db,
	}
}

// Create creates a new search history entry
func (r *GormSearchHistoryRepository) Create(ctx context.Context, searchHistory *model.SearchHistory) error {
	return r.db.WithContext(ctx).Create(searchHistory).Error
}

// FindByQuery finds search history entries by query
func (r *GormSearchHistoryRepository) FindByQuery(ctx context.Context, query string) ([]*model.SearchHistory, error) {
	var searchHistories []*model.SearchHistory
	err := r.db.WithContext(ctx).
		Where("query = ?", query).
		Order("created_at DESC").
		Find(&searchHistories).Error

	if err != nil {
		return nil, err
	}

	return searchHistories, nil
}

// FindPopular finds the most popular searches
func (r *GormSearchHistoryRepository) FindPopular(ctx context.Context, limit int) ([]*model.SearchHistory, error) {
	var searchHistories []*model.SearchHistory
	err := r.db.WithContext(ctx).
		Model(&model.SearchHistory{}).
		Select("query, COUNT(*) as count, MAX(created_at) as created_at, MAX(id) as id, SUM(result_count) as result_count").
		Group("query").
		Order("count DESC").
		Limit(limit).
		Find(&searchHistories).Error

	if err != nil {
		return nil, err
	}

	return searchHistories, nil
}

// Count counts all search history entries
func (r *GormSearchHistoryRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.SearchHistory{}).Count(&count).Error
	return count, err
}
