package repository

import (
	"context"

	"github.com/accio/internal/domain/model"
)

// ProfileRepository defines the interface for profile data access
type ProfileRepository interface {
	// Create creates a new profile
	Create(ctx context.Context, profile *model.Profile) error

	// Update updates an existing profile
	Update(ctx context.Context, profile *model.Profile) error

	// FindByID finds a profile by ID
	FindByID(ctx context.Context, id uint) (*model.Profile, error)

	// FindByUsername finds a profile by username and platform
	FindByUsername(ctx context.Context, username, platform string) (*model.Profile, error)

	// FindByRealName finds profiles by real name
	FindByRealName(ctx context.Context, name string) ([]*model.Profile, error)

	// FindByNamePart finds profiles by name part
	FindByNamePart(ctx context.Context, namePart string) ([]*model.Profile, error)

	// FindByAlias finds profiles by alias
	FindByAlias(ctx context.Context, alias string) ([]*model.Profile, error)

	// FindAll finds all profiles with optional limit and offset
	FindAll(ctx context.Context, limit, offset int) ([]*model.Profile, error)

	// Count counts all profiles
	Count(ctx context.Context) (int64, error)

	// Delete deletes a profile
	Delete(ctx context.Context, id uint) error
}
