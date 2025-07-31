package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/accio/internal/database"
)

// TwitchClient is a client for the Twitch API
type TwitchClient struct {
	*BaseClient
	ClientID     string
	ClientSecret string
	AccessToken  string
	TokenExpiry  time.Time
}

// TwitchUser represents a Twitch user from the API
type TwitchUser struct {
	ID              string    `json:"id"`
	Login           string    `json:"login"`
	DisplayName     string    `json:"display_name"`
	Type            string    `json:"type"`
	BroadcasterType string    `json:"broadcaster_type"`
	Description     string    `json:"description"`
	ProfileImageURL string    `json:"profile_image_url"`
	OfflineImageURL string    `json:"offline_image_url"`
	ViewCount       int       `json:"view_count"`
	CreatedAt       time.Time `json:"created_at"`
}

// TwitchSearchResponse represents a Twitch search response
type TwitchSearchResponse struct {
	Data []TwitchUser `json:"data"`
}

// TwitchTokenResponse represents a Twitch OAuth token response
type TwitchTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// NewTwitchClient creates a new Twitch API client
func NewTwitchClient() (*TwitchClient, error) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("TWITCH_CLIENT_ID and TWITCH_CLIENT_SECRET environment variables must be set")
	}

	client := &TwitchClient{
		BaseClient:   NewBaseClient(),
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	// Get initial access token
	if err := client.refreshAccessToken(); err != nil {
		return nil, err
	}

	return client, nil
}

// GetPlatformName returns the name of the platform
func (c *TwitchClient) GetPlatformName() string {
	return "Twitch"
}

// refreshAccessToken refreshes the Twitch API access token
func (c *TwitchClient) refreshAccessToken() error {
	// Build URL
	apiURL := "https://id.twitch.tv/oauth2/token"

	// Add query parameters
	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("client_secret", c.ClientSecret)
	params.Add("grant_type", "client_credentials")

	// Create request
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.UserAgent)

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status code %d", ErrAPIError, resp.StatusCode)
	}

	// Parse response
	var tokenResp TwitchTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Update client
	c.AccessToken = tokenResp.AccessToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

// ensureValidToken ensures the access token is valid
func (c *TwitchClient) ensureValidToken() error {
	if c.AccessToken == "" || time.Now().After(c.TokenExpiry) {
		return c.refreshAccessToken()
	}
	return nil
}

// GetProfileByUsername gets a Twitch profile by username
func (c *TwitchClient) GetProfileByUsername(username string) (*database.Profile, error) {
	// Ensure we have a valid token
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	// Build URL
	apiURL := "https://api.twitch.tv/helix/users"

	// Add query parameters
	params := url.Values{}
	params.Add("login", username)

	apiURL += "?" + params.Encode()

	// Create request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("User-Agent", c.UserAgent)

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusUnauthorized {
		// Token might be expired, refresh and try again
		if err := c.refreshAccessToken(); err != nil {
			return nil, err
		}
		return c.GetProfileByUsername(username)
	} else if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code %d", ErrAPIError, resp.StatusCode)
	}

	// Parse response
	var response TwitchSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if user was found
	if len(response.Data) == 0 {
		return nil, ErrNotFound
	}

	// Convert to profile
	return c.twitchUserToProfile(&response.Data[0]), nil
}

// SearchProfilesByName searches for Twitch profiles by real name
func (c *TwitchClient) SearchProfilesByName(name string) ([]*database.Profile, error) {
	// Ensure we have a valid token
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	// Twitch API doesn't support searching by real name directly
	// We'll search for channels instead, which is the closest approximation
	apiURL := "https://api.twitch.tv/helix/search/channels"

	// Add query parameters
	params := url.Values{}
	params.Add("query", name)
	params.Add("first", "10") // Limit to 10 results

	apiURL += "?" + params.Encode()

	// Create request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Client-ID", c.ClientID)
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("User-Agent", c.UserAgent)

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusUnauthorized {
		// Token might be expired, refresh and try again
		if err := c.refreshAccessToken(); err != nil {
			return nil, err
		}
		return c.SearchProfilesByName(name)
	} else if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code %d", ErrAPIError, resp.StatusCode)
	}

	// Parse response
	var response struct {
		Data []struct {
			ID              string `json:"id"`
			DisplayName     string `json:"display_name"`
			BroadcasterType string `json:"broadcaster_type"`
			GameName        string `json:"game_name"`
			Title           string `json:"title"`
			ThumbnailURL    string `json:"thumbnail_url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// For each channel, get the user details
	var profiles []*database.Profile
	for _, channel := range response.Data {
		// Get user details
		user, err := c.GetProfileByUsername(channel.DisplayName)
		if err != nil {
			// Skip this user if there's an error
			continue
		}
		profiles = append(profiles, user)
	}

	return profiles, nil
}

// GetProfileImage gets a Twitch profile image
func (c *TwitchClient) GetProfileImage(profile *database.Profile) (io.ReadCloser, error) {
	if profile.ImageURL == "" {
		return nil, fmt.Errorf("profile has no image URL")
	}

	return c.DownloadImage(profile.ImageURL)
}

// twitchUserToProfile converts a Twitch user to a profile
func (c *TwitchClient) twitchUserToProfile(user *TwitchUser) *database.Profile {
	// Create profile
	profile := &database.Profile{
		RealName:      user.DisplayName,
		Username:      user.Login,
		Platform:      "Twitch",
		ProfileURL:    fmt.Sprintf("https://twitch.tv/%s", user.Login),
		ImageURL:      user.ProfileImageURL,
		Verified:      user.BroadcasterType == "partner" || user.Type == "admin" || user.Type == "staff",
		FollowerCount: int64(user.ViewCount), // Twitch API doesn't provide follower count directly
		Bio:           user.Description,
		PlatformData:  make(map[string]string),
	}

	// Add platform data
	profile.PlatformData["user_id"] = user.ID
	profile.PlatformData["broadcaster_type"] = user.BroadcasterType
	profile.PlatformData["user_type"] = user.Type
	profile.PlatformData["view_count"] = fmt.Sprintf("%d", user.ViewCount)
	profile.PlatformData["created_at"] = user.CreatedAt.Format("2006-01-02")

	// Parse name parts
	nameParts := strings.Fields(user.DisplayName)
	if len(nameParts) > 0 {
		profile.NameParts = append(profile.NameParts, database.NamePart{
			NamePart: nameParts[0],
			PartType: "first",
		})
	}
	if len(nameParts) > 1 {
		profile.NameParts = append(profile.NameParts, database.NamePart{
			NamePart: nameParts[len(nameParts)-1],
			PartType: "last",
		})
	}
	if len(nameParts) > 2 {
		profile.NameParts = append(profile.NameParts, database.NamePart{
			NamePart: strings.Join(nameParts[1:len(nameParts)-1], " "),
			PartType: "middle",
		})
	}

	return profile
}
