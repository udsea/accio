package persistence

import (
	"context"

	"github.com/accio/internal/domain/model"
	"github.com/accio/internal/domain/repository"
	"gorm.io/gorm"
)

// GormUserFeedbackRepository is a GORM implementation of UserFeedbackRepository
type GormUserFeedbackRepository struct {
	db *gorm.DB
}

// NewGormUserFeedbackRepository creates a new GormUserFeedbackRepository
func NewGormUserFeedbackRepository(db *gorm.DB) repository.UserFeedbackRepository {
	return &GormUserFeedbackRepository{
		db: db,
	}
}

// Create creates a new user feedback entry
func (r *GormUserFeedbackRepository) Create(ctx context.Context, userFeedback *model.UserFeedback) error {
	return r.db.WithContext(ctx).Create(userFeedback).Error
}

// FindByProfileID finds user feedback entries by profile ID
func (r *GormUserFeedbackRepository) FindByProfileID(ctx context.Context, profileID uint) ([]*model.UserFeedback, error) {
	var userFeedbacks []*model.UserFeedback
	err := r.db.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Order("created_at DESC").
		Find(&userFeedbacks).Error

	if err != nil {
		return nil, err
	}

	return userFeedbacks, nil
}

// FindByFeedbackType finds user feedback entries by feedback type
func (r *GormUserFeedbackRepository) FindByFeedbackType(ctx context.Context, feedbackType string) ([]*model.UserFeedback, error) {
	var userFeedbacks []*model.UserFeedback
	err := r.db.WithContext(ctx).
		Where("feedback_type = ?", feedbackType).
		Order("created_at DESC").
		Find(&userFeedbacks).Error

	if err != nil {
		return nil, err
	}

	return userFeedbacks, nil
}

// Count counts all user feedback entries
func (r *GormUserFeedbackRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserFeedback{}).Count(&count).Error
	return count, err
}
