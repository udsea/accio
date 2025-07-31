package model

import "time"

// SearchHistory represents a search query and its results
type SearchHistory struct {
	ID          uint   `gorm:"primaryKey"`
	Query       string `gorm:"index"`
	ResultCount int
	CreatedAt   time.Time
}

// NewSearchHistory creates a new search history entity
func NewSearchHistory(query string, resultCount int) *SearchHistory {
	return &SearchHistory{
		Query:       query,
		ResultCount: resultCount,
		CreatedAt:   time.Now(),
	}
}
