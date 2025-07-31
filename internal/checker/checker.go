package checker

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Checker handles the process of checking usernames across different sites
type Checker struct {
	Client    *http.Client
	Verbose   bool
	UserAgent string
	Mutex     sync.Mutex
	Stats     CheckStats
}

// CheckStats tracks statistics about the checking process
type CheckStats struct {
	Total     int
	Found     int
	NotFound  int
	Errors    int
	StartTime time.Time
	EndTime   time.Time
}

// NewChecker creates a new Checker instance with the specified timeout
func NewChecker(timeout int, verbose bool) *Checker {
	return &Checker{
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				MaxConnsPerHost:     100,
				IdleConnTimeout:     30 * time.Second,
			},
		},
		Verbose:   verbose,
		UserAgent: "Accio/1.0",
		Stats: CheckStats{
			StartTime: time.Now(),
		},
	}
}

// CheckUsername checks if a username exists on a given site
func (c *Checker) CheckUsername(username, site, url string) (bool, error) {
	if c.Verbose {
		fmt.Printf("Checking %s on %s\n", username, site)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.Client.Timeout)
	defer cancel()

	// Create a new request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.incrementErrors()
		return false, err
	}

	// Set a user agent to avoid being blocked
	req.Header.Set("User-Agent", c.UserAgent)

	// Make the request
	resp, err := c.Client.Do(req)
	if err != nil {
		c.incrementErrors()
		return false, err
	}
	defer resp.Body.Close()

	// Check if the response indicates the username exists
	// This is a simple implementation that just checks status codes
	// A more sophisticated implementation would check for specific patterns in the response
	exists := resp.StatusCode >= 200 && resp.StatusCode < 300

	// Some sites return 200 even if the user doesn't exist, so we'd need more checks
	// This is just a placeholder for now
	if site == "GitHub" && exists {
		// For GitHub, check if the response URL contains "404" which indicates the user doesn't exist
		if strings.Contains(resp.Request.URL.String(), "404") {
			exists = false
		}
	}

	// Update stats
	if exists {
		c.incrementFound()
	} else {
		c.incrementNotFound()
	}

	return exists, nil
}

// CheckUsernameWithRetry checks a username with retry logic
func (c *Checker) CheckUsernameWithRetry(username, site, url string, maxRetries int) (bool, error) {
	var lastErr error

	for retry := 0; retry < maxRetries; retry++ {
		exists, err := c.CheckUsername(username, site, url)
		if err == nil {
			return exists, nil
		}

		lastErr = err

		// Wait before retrying (with exponential backoff)
		time.Sleep(time.Duration(retry+1) * 500 * time.Millisecond)
	}

	return false, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// GetStats returns the current statistics
func (c *Checker) GetStats() CheckStats {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// Update end time
	c.Stats.EndTime = time.Now()
	return c.Stats
}

// incrementFound increments the found counter
func (c *Checker) incrementFound() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Stats.Found++
	c.Stats.Total++
}

// incrementNotFound increments the not found counter
func (c *Checker) incrementNotFound() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Stats.NotFound++
	c.Stats.Total++
}

// incrementErrors increments the errors counter
func (c *Checker) incrementErrors() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Stats.Errors++
	c.Stats.Total++
}
