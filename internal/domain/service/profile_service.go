package service

import (
	"context"
	"io"

	"github.com/accio/internal/application/dto"
	"github.com/accio/internal/domain/model"
)

// ProfileService defines the interface for profile-related operations
type ProfileService interface {
	// GetProfileByUsername gets a profile by username from a specific platform
	GetProfileByUsername(ctx context.Context, username, platform string) (*dto.ProfileDTO, error)

	// SearchProfilesByName searches for profiles by real name
	SearchProfilesByName(ctx context.Context, name string) ([]*dto.ProfileDTO, error)

	// GetProfileImage gets a profile image
	GetProfileImage(ctx context.Context, profile *model.Profile) (io.ReadCloser, error)

	// SaveProfile saves a profile to the repository
	SaveProfile(ctx context.Context, profile *model.Profile) error

	// GetSupportedPlatforms returns a list of supported platforms
	GetSupportedPlatforms() []string
}
