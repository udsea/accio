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

// GitHubClient is a client for the GitHub API
type GitHubClient struct {
	*BaseClient
	Token string
}

// GitHubUser represents a GitHub user from the API
type GitHubUser struct {
	Login             string    `json:"login"`
	ID                int       `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	TwitterUsername   string    `json:"twitter_username"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// GitHubSearchResponse represents a GitHub search response
type GitHubSearchResponse struct {
	TotalCount        int          `json:"total_count"`
	IncompleteResults bool         `json:"incomplete_results"`
	Items             []GitHubUser `json:"items"`
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient() (*GitHubClient, error) {
	token := os.Getenv("GITHUB_TOKEN")

	client := &GitHubClient{
		BaseClient: NewBaseClient(),
		Token:      token,
	}

	return client, nil
}

// GetPlatformName returns the name of the platform
func (c *GitHubClient) GetPlatformName() string {
	return "GitHub"
}

// GetProfileByUsername gets a GitHub profile by username
func (c *GitHubClient) GetProfileByUsername(username string) (*database.Profile, error) {
	// Build URL
	apiURL := fmt.Sprintf("https://api.github.com/users/%s", url.PathEscape(username))

	// Create request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", c.UserAgent)
	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	} else if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	} else if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code %d", ErrAPIError, resp.StatusCode)
	}

	// Parse response
	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to profile
	return c.githubUserToProfile(&user), nil
}

// SearchProfilesByName searches for GitHub profiles by real name
func (c *GitHubClient) SearchProfilesByName(name string) ([]*database.Profile, error) {
	// Build URL
	apiURL := "https://api.github.com/search/users"

	// Add query parameters
	params := url.Values{}
	params.Add("q", name+" in:name")
	params.Add("per_page", "10")

	apiURL += "?" + params.Encode()

	// Create request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", c.UserAgent)
	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	} else if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code %d", ErrAPIError, resp.StatusCode)
	}

	// Parse response
	var searchResp GitHubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// For each user in the search results, get the full profile
	var profiles []*database.Profile
	for _, user := range searchResp.Items {
		profile, err := c.GetProfileByUsername(user.Login)
		if err != nil {
			// Skip this user if there's an error
			continue
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// GetProfileImage gets a GitHub profile image
func (c *GitHubClient) GetProfileImage(profile *database.Profile) (io.ReadCloser, error) {
	if profile.ImageURL == "" {
		return nil, fmt.Errorf("profile has no image URL")
	}

	return c.DownloadImage(profile.ImageURL)
}

// githubUserToProfile converts a GitHub user to a profile
func (c *GitHubClient) githubUserToProfile(user *GitHubUser) *database.Profile {
	// Create profile
	profile := &database.Profile{
		RealName:      user.Name,
		Username:      user.Login,
		Platform:      "GitHub",
		ProfileURL:    user.HTMLURL,
		ImageURL:      user.AvatarURL,
		Verified:      user.SiteAdmin,
		FollowerCount: int64(user.Followers),
		Bio:           user.Bio,
		PlatformData:  make(map[string]string),
	}

	// Add platform data
	profile.PlatformData["user_id"] = fmt.Sprintf("%d", user.ID)
	profile.PlatformData["public_repos"] = fmt.Sprintf("%d", user.PublicRepos)
	profile.PlatformData["public_gists"] = fmt.Sprintf("%d", user.PublicGists)
	profile.PlatformData["following"] = fmt.Sprintf("%d", user.Following)
	profile.PlatformData["created_at"] = user.CreatedAt.Format("2006-01-02")
	profile.PlatformData["updated_at"] = user.UpdatedAt.Format("2006-01-02")

	if user.Email != "" {
		profile.PlatformData["email"] = user.Email
	}
	if user.Company != "" {
		profile.PlatformData["company"] = user.Company
	}
	if user.Location != "" {
		profile.PlatformData["location"] = user.Location
	}
	if user.Blog != "" {
		profile.PlatformData["blog"] = user.Blog
	}
	if user.TwitterUsername != "" {
		profile.PlatformData["twitter_username"] = user.TwitterUsername
	}

	// Parse name parts
	if user.Name != "" {
		nameParts := strings.Fields(user.Name)
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
	}

	return profile
}
