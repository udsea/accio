package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/accio/internal/domain/model"
)

// TwitterClient is a client for the Twitter API
type TwitterClient struct {
	*BaseClient
	BearerToken string
}

// TwitterUser represents a Twitter user from the API
type TwitterUser struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profile_image_url"`
	Description     string `json:"description"`
	Verified        bool   `json:"verified"`
	PublicMetrics   struct {
		FollowersCount int `json:"followers_count"`
		FollowingCount int `json:"following_count"`
		TweetCount     int `json:"tweet_count"`
		ListedCount    int `json:"listed_count"`
	} `json:"public_metrics"`
	CreatedAt string `json:"created_at"`
}

// TwitterSearchResponse represents a Twitter search response
type TwitterSearchResponse struct {
	Data []TwitterUser `json:"data"`
	Meta struct {
		ResultCount int    `json:"result_count"`
		NextToken   string `json:"next_token"`
	} `json:"meta"`
}

// NewTwitterClient creates a new Twitter API client
func NewTwitterClient() (PlatformClient, error) {
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	if bearerToken == "" {
		return nil, fmt.Errorf("TWITTER_BEARER_TOKEN environment variable not set")
	}

	return &TwitterClient{
		BaseClient:  NewBaseClient(),
		BearerToken: bearerToken,
	}, nil
}

// GetPlatformName returns the name of the platform
func (c *TwitterClient) GetPlatformName() string {
	return "Twitter"
}

// GetProfileByUsername gets a Twitter profile by username
func (c *TwitterClient) GetProfileByUsername(ctx context.Context, username string) (*model.Profile, error) {
	// Clean username (remove @ if present)
	username = strings.TrimPrefix(username, "@")

	// Build URL
	apiURL := fmt.Sprintf("https://api.twitter.com/2/users/by/username/%s", url.PathEscape(username))

	// Add query parameters for additional fields
	apiURL += "?user.fields=id,name,username,profile_image_url,description,verified,public_metrics,created_at"

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("User-Agent", c.UserAgent)

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	} else if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	} else if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code %d", ErrAPIError, resp.StatusCode)
	}

	// Parse response
	var response struct {
		Data TwitterUser `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to profile
	return c.twitterUserToProfile(&response.Data), nil
}

// SearchProfilesByName searches for Twitter profiles by real name
func (c *TwitterClient) SearchProfilesByName(ctx context.Context, name string) ([]*model.Profile, error) {
	// Build URL
	apiURL := "https://api.twitter.com/2/users/search"

	// Add query parameters
	params := url.Values{}
	params.Add("query", name)
	params.Add("max_results", "10")
	params.Add("user.fields", "id,name,username,profile_image_url,description,verified,public_metrics,created_at")

	apiURL += "?" + params.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("User-Agent", c.UserAgent)

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	} else if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code %d", ErrAPIError, resp.StatusCode)
	}

	// Parse response
	var response TwitterSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to profiles
	var profiles []*model.Profile
	for _, user := range response.Data {
		profiles = append(profiles, c.twitterUserToProfile(&user))
	}

	return profiles, nil
}

// GetProfileImage gets a Twitter profile image
func (c *TwitterClient) GetProfileImage(ctx context.Context, profile *model.Profile) (io.ReadCloser, error) {
	if profile.ImageURL == "" {
		return nil, fmt.Errorf("profile has no image URL")
	}

	// Twitter API returns a small image by default, replace with original size
	imageURL := strings.Replace(profile.ImageURL, "_normal", "", 1)

	return c.DownloadImage(ctx, imageURL)
}

// twitterUserToProfile converts a Twitter user to a profile
func (c *TwitterClient) twitterUserToProfile(user *TwitterUser) *model.Profile {
	// Create profile
	profile := model.NewProfile(
		user.Name,
		user.Username,
		"Twitter",
		fmt.Sprintf("https://twitter.com/%s", user.Username),
		user.ProfileImageURL,
		user.Description,
		user.Verified,
		int64(user.PublicMetrics.FollowersCount),
	)

	// Parse created at date
	createdAt, _ := time.Parse(time.RFC3339, user.CreatedAt)

	// Add platform data
	profile.AddPlatformData("user_id", user.ID)
	profile.AddPlatformData("following_count", fmt.Sprintf("%d", user.PublicMetrics.FollowingCount))
	profile.AddPlatformData("tweet_count", fmt.Sprintf("%d", user.PublicMetrics.TweetCount))
	profile.AddPlatformData("listed_count", fmt.Sprintf("%d", user.PublicMetrics.ListedCount))
	profile.AddPlatformData("created_at", createdAt.Format("2006-01-02"))

	// Parse name parts
	nameParts := strings.Fields(user.Name)
	if len(nameParts) > 0 {
		profile.AddNamePart(nameParts[0], "first")
	}
	if len(nameParts) > 1 {
		profile.AddNamePart(nameParts[len(nameParts)-1], "last")
	}
	if len(nameParts) > 2 {
		profile.AddNamePart(strings.Join(nameParts[1:len(nameParts)-1], " "), "middle")
	}

	return profile
}
