package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/accio/internal/domain/model"
)

// Common errors
var (
	ErrNotFound      = errors.New("profile not found")
	ErrRateLimited   = errors.New("rate limited by API")
	ErrUnauthorized  = errors.New("unauthorized API access")
	ErrAPIError      = errors.New("API error")
	ErrInvalidParams = errors.New("invalid parameters")
)

// PlatformClient defines the interface for platform-specific API clients
type PlatformClient interface {
	// GetProfileByUsername gets a profile by username
	GetProfileByUsername(ctx context.Context, username string) (*model.Profile, error)

	// SearchProfilesByName searches for profiles by real name
	SearchProfilesByName(ctx context.Context, name string) ([]*model.Profile, error)

	// GetProfileImage gets a profile image
	GetProfileImage(ctx context.Context, profile *model.Profile) (io.ReadCloser, error)

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
func (c *BaseClient) DownloadImage(ctx context.Context, imageURL string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.New("failed to download image")
	}

	return resp.Body, nil
}
