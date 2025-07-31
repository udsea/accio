package persistence

import (
	"context"
	"errors"

	"github.com/accio/internal/domain/model"
	"github.com/accio/internal/domain/repository"
	"gorm.io/gorm"
)

// GormProfileRepository is a GORM implementation of ProfileRepository
type GormProfileRepository struct {
	db *gorm.DB
}

// NewGormProfileRepository creates a new GormProfileRepository
func NewGormProfileRepository(db *gorm.DB) repository.ProfileRepository {
	return &GormProfileRepository{
		db: db,
	}
}

// Create creates a new profile
func (r *GormProfileRepository) Create(ctx context.Context, profile *model.Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// Update updates an existing profile
func (r *GormProfileRepository) Update(ctx context.Context, profile *model.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// FindByID finds a profile by ID
func (r *GormProfileRepository) FindByID(ctx context.Context, id uint) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.WithContext(ctx).
		Preload("NameParts").
		Preload("Aliases").
		Preload("PlatformData").
		First(&profile, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

// FindByUsername finds a profile by username and platform
func (r *GormProfileRepository) FindByUsername(ctx context.Context, username, platform string) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.WithContext(ctx).
		Preload("NameParts").
		Preload("Aliases").
		Preload("PlatformData").
		Where("username = ? AND platform = ?", username, platform).
		First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

// FindByRealName finds profiles by real name
func (r *GormProfileRepository) FindByRealName(ctx context.Context, name string) ([]*model.Profile, error) {
	var profiles []*model.Profile
	err := r.db.WithContext(ctx).
		Preload("NameParts").
		Preload("Aliases").
		Preload("PlatformData").
		Where("real_name LIKE ?", "%"+name+"%").
		Find(&profiles).Error

	if err != nil {
		return nil, err
	}

	return profiles, nil
}

// FindByNamePart finds profiles by name part
func (r *GormProfileRepository) FindByNamePart(ctx context.Context, namePart string) ([]*model.Profile, error) {
	var profiles []*model.Profile
	err := r.db.WithContext(ctx).
		Preload("NameParts").
		Preload("Aliases").
		Preload("PlatformData").
		Joins("JOIN name_parts ON name_parts.profile_id = profiles.id").
		Where("name_parts.name_part LIKE ?", "%"+namePart+"%").
		Find(&profiles).Error

	if err != nil {
		return nil, err
	}

	return profiles, nil
}

// FindByAlias finds profiles by alias
func (r *GormProfileRepository) FindByAlias(ctx context.Context, alias string) ([]*model.Profile, error) {
	var profiles []*model.Profile
	err := r.db.WithContext(ctx).
		Preload("NameParts").
		Preload("Aliases").
		Preload("PlatformData").
		Joins("JOIN aliases ON aliases.profile_id = profiles.id").
		Where("aliases.alias LIKE ?", "%"+alias+"%").
		Find(&profiles).Error

	if err != nil {
		return nil, err
	}

	return profiles, nil
}

// FindAll finds all profiles with optional limit and offset
func (r *GormProfileRepository) FindAll(ctx context.Context, limit, offset int) ([]*model.Profile, error) {
	var profiles []*model.Profile
	query := r.db.WithContext(ctx).
		Preload("NameParts").
		Preload("Aliases").
		Preload("PlatformData")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&profiles).Error
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

// Count counts all profiles
func (r *GormProfileRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Profile{}).Count(&count).Error
	return count, err
}

// Delete deletes a profile
func (r *GormProfileRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Profile{}, id).Error
}
