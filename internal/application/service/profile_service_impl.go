package service

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/accio/internal/domain/model"
	"github.com/accio/internal/domain/repository"
	"github.com/accio/internal/domain/service"
	"github.com/accio/internal/infrastructure/api"
)

// ProfileServiceImpl implements the ProfileService interface
type ProfileServiceImpl struct {
	profileRepo     repository.ProfileRepository
	platformClients map[string]api.PlatformClient
	clientsMutex    sync.RWMutex
}

// NewProfileService creates a new ProfileServiceImpl
func NewProfileService(profileRepo repository.ProfileRepository) service.ProfileService {
	return &ProfileServiceImpl{
		profileRepo:     profileRepo,
		platformClients: make(map[string]api.PlatformClient),
	}
}

// RegisterPlatformClient registers a platform client
func (s *ProfileServiceImpl) RegisterPlatformClient(client api.PlatformClient) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	s.platformClients[client.GetPlatformName()] = client
}

// GetProfileByUsername gets a profile by username from a specific platform
func (s *ProfileServiceImpl) GetProfileByUsername(ctx context.Context, username, platform string) (*model.Profile, error) {
	// First, try to get from repository
	profile, err := s.profileRepo.FindByUsername(ctx, username, platform)
	if err != nil {
		return nil, err
	}

	// If found in repository, return it
	if profile != nil {
		return profile, nil
	}

	// If not found in repository, try to get from platform API
	s.clientsMutex.RLock()
	client, ok := s.platformClients[platform]
	s.clientsMutex.RUnlock()

	if !ok {
		return nil, errors.New("unsupported platform")
	}

	// Get from platform API
	profile, err = client.GetProfileByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if profile != nil {
		if err := s.profileRepo.Create(ctx, profile); err != nil {
			// Log error but continue
			// log.Printf("Error saving profile to repository: %v", err)
		}
	}

	return profile, nil
}

// SearchProfilesByName searches for profiles by real name
func (s *ProfileServiceImpl) SearchProfilesByName(ctx context.Context, name string) ([]*model.Profile, error) {
	// First, try to get from repository
	profiles, err := s.profileRepo.FindByRealName(ctx, name)
	if err != nil {
		return nil, err
	}

	// If found in repository, return them
	if len(profiles) > 0 {
		return profiles, nil
	}

	// If not found in repository, try to get from platform APIs
	var allProfiles []*model.Profile
	var wg sync.WaitGroup
	var mu sync.Mutex

	s.clientsMutex.RLock()
	for _, client := range s.platformClients {
		wg.Add(1)
		go func(client api.PlatformClient) {
			defer wg.Done()

			// Get from platform API
			clientProfiles, err := client.SearchProfilesByName(ctx, name)
			if err != nil {
				// Log error but continue
				// log.Printf("Error searching profiles from %s: %v", client.GetPlatformName(), err)
				return
			}

			// Save to repository
			for _, profile := range clientProfiles {
				if err := s.profileRepo.Create(ctx, profile); err != nil {
					// Log error but continue
					// log.Printf("Error saving profile to repository: %v", err)
				}
			}

			// Add to results
			mu.Lock()
			allProfiles = append(allProfiles, clientProfiles...)
			mu.Unlock()
		}(client)
	}
	s.clientsMutex.RUnlock()

	// Wait for all goroutines to finish
	wg.Wait()

	return allProfiles, nil
}

// GetProfileImage gets a profile image
func (s *ProfileServiceImpl) GetProfileImage(ctx context.Context, profile *model.Profile) (io.ReadCloser, error) {
	// Get platform client
	s.clientsMutex.RLock()
	client, ok := s.platformClients[profile.Platform]
	s.clientsMutex.RUnlock()

	if !ok {
		return nil, errors.New("unsupported platform")
	}

	// Get image from platform API
	return client.GetProfileImage(ctx, profile)
}

// SaveProfile saves a profile to the repository
func (s *ProfileServiceImpl) SaveProfile(ctx context.Context, profile *model.Profile) error {
	// Check if profile already exists
	existingProfile, err := s.profileRepo.FindByUsername(ctx, profile.Username, profile.Platform)
	if err != nil {
		return err
	}

	if existingProfile != nil {
		// Update existing profile
		profile.ID = existingProfile.ID
		return s.profileRepo.Update(ctx, profile)
	}

	// Create new profile
	return s.profileRepo.Create(ctx, profile)
}

// GetSupportedPlatforms returns a list of supported platforms
func (s *ProfileServiceImpl) GetSupportedPlatforms() []string {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()

	platforms := make([]string, 0, len(s.platformClients))
	for platform := range s.platformClients {
		platforms = append(platforms, platform)
	}

	return platforms
}
