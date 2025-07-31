package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is a wrapper around http.Client with additional functionality
type Client struct {
	*http.Client
	UserAgent string
	Verbose   bool
}

// NewClient creates a new HTTP client with the specified timeout
func NewClient(timeout int, verbose bool) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		UserAgent: "Accio/1.0",
		Verbose:   verbose,
	}
}

// Get performs an HTTP GET request to the specified URL
func (c *Client) Get(url string) (*http.Response, error) {
	// Replace {} in URL with the username
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set User-Agent header
	req.Header.Set("User-Agent", c.UserAgent)

	if c.Verbose {
		fmt.Printf("Making request to: %s\n", url)
	}

	return c.Client.Do(req)
}

// CheckURL checks if a URL exists and returns the response body
func (c *Client) CheckURL(url string) (bool, string, error) {
	resp, err := c.Get(url)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	// Consider 2xx status codes as "exists"
	exists := resp.StatusCode >= 200 && resp.StatusCode < 300

	return exists, string(body), nil
}

// FormatURL replaces the {} placeholder in a URL with the username
func FormatURL(urlFormat string, username string) string {
	return strings.ReplaceAll(urlFormat, "{}", username)
}
