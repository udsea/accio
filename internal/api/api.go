package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/accio/internal/database"
)

// Common errors
var (
	ErrNotFound      = errors.New("profile not found")
	ErrRateLimited   = errors.New("rate limited by API")
	ErrUnauthorized  = errors.New("unauthorized API access")
	ErrAPIError      = errors.New("API error")
	ErrInvalidParams = errors.New("invalid parameters")
)

// ProfileClient is the interface for all platform-specific API clients
type ProfileClient interface {
	// GetProfileByUsername gets a profile by username
	GetProfileByUsername(username string) (*database.Profile, error)

	// SearchProfilesByName searches for profiles by real name
	SearchProfilesByName(name string) ([]*database.Profile, error)

	// GetProfileImage gets a profile image
	GetProfileImage(profile *database.Profile) (io.ReadCloser, error)

	// GetPlatformName returns the name of the platform
	GetPlatformName() string
}

// BaseClient provides common functionality for all API clients
type BaseClient struct {
	HTTPClient *http.Client
	UserAgent  string
}

// NewBaseClient creates a new base client
func NewBaseClient() *BaseClient {
	return &BaseClient{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		UserAgent: "Accio/1.0 (+https://github.com/yourusername/accio)",
	}
}

// DownloadImage downloads an image from a URL
func (c *BaseClient) DownloadImage(imageURL string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// SaveImageToFile saves an image to a file
func SaveImageToFile(imageData io.ReadCloser, filePath string) error {
	defer imageData.Close()

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy image data to file
	_, err = io.Copy(file, imageData)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}

// GetClientForPlatform returns the appropriate client for a platform
func GetClientForPlatform(platform string) (ProfileClient, error) {
	switch platform {
	case "Twitter", "X":
		return NewTwitterClient()
	case "Twitch":
		return NewTwitchClient()
	case "Instagram":
		return NewInstagramClient()
	case "GitHub":
		return NewGitHubClient()
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}
