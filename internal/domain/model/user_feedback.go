package model

import "time"

// UserFeedback represents user feedback for a profile
type UserFeedback struct {
	ID           uint   `gorm:"primaryKey"`
	ProfileID    uint   `gorm:"index"`
	FeedbackType string // 'correct', 'incorrect', 'missing'
	Comment      string
	CreatedAt    time.Time
}

// NewUserFeedback creates a new user feedback entity
func NewUserFeedback(profileID uint, feedbackType, comment string) *UserFeedback {
	return &UserFeedback{
		ProfileID:    profileID,
		FeedbackType: feedbackType,
		Comment:      comment,
		CreatedAt:    time.Now(),
	}
}
