package repository

import (
	"context"

	"github.com/accio/internal/domain/model"
)

// UserFeedbackRepository defines the interface for user feedback data access
type UserFeedbackRepository interface {
	// Create creates a new user feedback entry
	Create(ctx context.Context, userFeedback *model.UserFeedback) error

	// FindByProfileID finds user feedback entries by profile ID
	FindByProfileID(ctx context.Context, profileID uint) ([]*model.UserFeedback, error)

	// FindByFeedbackType finds user feedback entries by feedback type
	FindByFeedbackType(ctx context.Context, feedbackType string) ([]*model.UserFeedback, error)

	// Count counts all user feedback entries
	Count(ctx context.Context) (int64, error)
}
