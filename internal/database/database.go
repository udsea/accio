package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

//go:embed schema.sql
var schemaFS embed.FS

// Profile represents a social media profile in the database
type Profile struct {
	ID            int64
	RealName      string
	Username      string
	Platform      string
	ProfileURL    string
	ImageURL      string
	Verified      bool
	FollowerCount int64
	Bio           string
	LastUpdated   time.Time
	NameParts     []NamePart
	Aliases       []string
	PlatformData  map[string]string
}

// NamePart represents a part of a person's name
type NamePart struct {
	ID        int64
	ProfileID int64
	NamePart  string
	PartType  string // 'first', 'middle', 'last', 'nickname'
}

// Client is a database client for Turso
type Client struct {
	db *sql.DB
}

// NewClient creates a new database client
func NewClient() (*Client, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found, using environment variables\n")
	}

	dbURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if dbURL == "" {
		return nil, fmt.Errorf("TURSO_DATABASE_URL environment variable not set")
	}

	// Construct connection string with auth token if provided
	connStr := dbURL
	if authToken != "" {
		connStr = fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
	}

	// Open database connection
	db, err := sql.Open("libsql", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	client := &Client{db: db}

	// Initialize database schema
	if err := client.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return client, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	return c.db.Close()
}

// initSchema initializes the database schema
func (c *Client) initSchema() error {
	// Read schema file
	schemaBytes, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Split schema into individual statements
	schema := string(schemaBytes)
	statements := strings.Split(schema, ";")

	// Execute each statement
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		_, err := c.db.Exec(statement)
		if err != nil {
			return fmt.Errorf("failed to execute schema statement: %w", err)
		}
	}

	return nil
}

// AddProfile adds a new profile to the database
func (c *Client) AddProfile(profile *Profile) error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert profile
	result, err := tx.Exec(
		`INSERT INTO profiles 
		(real_name, username, platform, profile_url, image_url, verified, follower_count, bio) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(username, platform) 
		DO UPDATE SET 
			real_name = excluded.real_name,
			profile_url = excluded.profile_url,
			image_url = excluded.image_url,
			verified = excluded.verified,
			follower_count = excluded.follower_count,
			bio = excluded.bio,
			last_updated = CURRENT_TIMESTAMP
		RETURNING id`,
		profile.RealName, profile.Username, profile.Platform, profile.ProfileURL,
		profile.ImageURL, profile.Verified, profile.FollowerCount, profile.Bio,
	)
	if err != nil {
		return fmt.Errorf("failed to insert profile: %w", err)
	}

	// Get profile ID
	var profileID int64
	profileID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get profile ID: %w", err)
	}
	profile.ID = profileID

	// Insert name parts
	for _, namePart := range profile.NameParts {
		_, err = tx.Exec(
			`INSERT INTO name_parts (profile_id, name_part, part_type) VALUES (?, ?, ?)
			ON CONFLICT(profile_id, name_part, part_type) DO NOTHING`,
			profileID, namePart.NamePart, namePart.PartType,
		)
		if err != nil {
			return fmt.Errorf("failed to insert name part: %w", err)
		}
	}

	// Insert aliases
	for _, alias := range profile.Aliases {
		_, err = tx.Exec(
			`INSERT INTO aliases (profile_id, alias) VALUES (?, ?)
			ON CONFLICT(profile_id, alias) DO NOTHING`,
			profileID, alias,
		)
		if err != nil {
			return fmt.Errorf("failed to insert alias: %w", err)
		}
	}

	// Insert platform data
	for key, value := range profile.PlatformData {
		_, err = tx.Exec(
			`INSERT INTO platform_data (profile_id, data_key, data_value) VALUES (?, ?, ?)
			ON CONFLICT(profile_id, data_key) DO UPDATE SET data_value = excluded.data_value`,
			profileID, key, value,
		)
		if err != nil {
			return fmt.Errorf("failed to insert platform data: %w", err)
		}
	}

	return tx.Commit()
}

// GetProfileByUsername gets a profile by username and platform
func (c *Client) GetProfileByUsername(username, platform string) (*Profile, error) {
	row := c.db.QueryRow(
		`SELECT id, real_name, username, platform, profile_url, image_url, verified, follower_count, bio, last_updated
		FROM profiles
		WHERE username = ? AND platform = ?`,
		username, platform,
	)

	profile := &Profile{
		PlatformData: make(map[string]string),
	}
	err := row.Scan(
		&profile.ID, &profile.RealName, &profile.Username, &profile.Platform,
		&profile.ProfileURL, &profile.ImageURL, &profile.Verified, &profile.FollowerCount,
		&profile.Bio, &profile.LastUpdated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan profile: %w", err)
	}

	// Get name parts
	rows, err := c.db.Query(
		`SELECT id, name_part, part_type FROM name_parts WHERE profile_id = ?`,
		profile.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query name parts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var namePart NamePart
		namePart.ProfileID = profile.ID
		err := rows.Scan(&namePart.ID, &namePart.NamePart, &namePart.PartType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan name part: %w", err)
		}
		profile.NameParts = append(profile.NameParts, namePart)
	}

	// Get aliases
	rows, err = c.db.Query(
		`SELECT alias FROM aliases WHERE profile_id = ?`,
		profile.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query aliases: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alias string
		err := rows.Scan(&alias)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alias: %w", err)
		}
		profile.Aliases = append(profile.Aliases, alias)
	}

	// Get platform data
	rows, err = c.db.Query(
		`SELECT data_key, data_value FROM platform_data WHERE profile_id = ?`,
		profile.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query platform data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan platform data: %w", err)
		}
		profile.PlatformData[key] = value
	}

	return profile, nil
}

// SearchProfilesByName searches for profiles by real name
func (c *Client) SearchProfilesByName(name string) ([]*Profile, error) {
	// Record search in history
	_, err := c.db.Exec(
		`INSERT INTO search_history (query, result_count) VALUES (?, 0)`,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to record search: %w", err)
	}

	// Search for profiles
	rows, err := c.db.Query(
		`SELECT DISTINCT p.id, p.real_name, p.username, p.platform, p.profile_url, 
		p.image_url, p.verified, p.follower_count, p.bio, p.last_updated
		FROM profiles p
		LEFT JOIN name_parts np ON p.id = np.profile_id
		LEFT JOIN aliases a ON p.id = a.profile_id
		WHERE p.real_name LIKE ? OR np.name_part LIKE ? OR a.alias LIKE ?
		ORDER BY p.verified DESC, p.follower_count DESC`,
		"%"+name+"%", "%"+name+"%", "%"+name+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search profiles: %w", err)
	}
	defer rows.Close()

	var profiles []*Profile
	profileMap := make(map[int64]*Profile)

	for rows.Next() {
		profile := &Profile{
			PlatformData: make(map[string]string),
		}
		err := rows.Scan(
			&profile.ID, &profile.RealName, &profile.Username, &profile.Platform,
			&profile.ProfileURL, &profile.ImageURL, &profile.Verified, &profile.FollowerCount,
			&profile.Bio, &profile.LastUpdated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan profile: %w", err)
		}

		// Check if we've already seen this profile
		if _, ok := profileMap[profile.ID]; !ok {
			profileMap[profile.ID] = profile
			profiles = append(profiles, profile)
		}
	}

	// Update search history with result count
	_, err = c.db.Exec(
		`UPDATE search_history SET result_count = ? WHERE query = ? AND id = (
			SELECT id FROM search_history WHERE query = ? ORDER BY timestamp DESC LIMIT 1
		)`,
		len(profiles), name, name,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update search history: %w", err)
	}

	// Load additional data for each profile
	for _, profile := range profiles {
		// Get name parts
		nameParts, err := c.getNameParts(profile.ID)
		if err != nil {
			return nil, err
		}
		profile.NameParts = nameParts

		// Get aliases
		aliases, err := c.getAliases(profile.ID)
		if err != nil {
			return nil, err
		}
		profile.Aliases = aliases

		// Get platform data
		platformData, err := c.getPlatformData(profile.ID)
		if err != nil {
			return nil, err
		}
		profile.PlatformData = platformData
	}

	return profiles, nil
}

// getNameParts gets name parts for a profile
func (c *Client) getNameParts(profileID int64) ([]NamePart, error) {
	rows, err := c.db.Query(
		`SELECT id, name_part, part_type FROM name_parts WHERE profile_id = ?`,
		profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query name parts: %w", err)
	}
	defer rows.Close()

	var nameParts []NamePart
	for rows.Next() {
		var namePart NamePart
		namePart.ProfileID = profileID
		err := rows.Scan(&namePart.ID, &namePart.NamePart, &namePart.PartType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan name part: %w", err)
		}
		nameParts = append(nameParts, namePart)
	}

	return nameParts, nil
}

// getAliases gets aliases for a profile
func (c *Client) getAliases(profileID int64) ([]string, error) {
	rows, err := c.db.Query(
		`SELECT alias FROM aliases WHERE profile_id = ?`,
		profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query aliases: %w", err)
	}
	defer rows.Close()

	var aliases []string
	for rows.Next() {
		var alias string
		err := rows.Scan(&alias)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alias: %w", err)
		}
		aliases = append(aliases, alias)
	}

	return aliases, nil
}

// getPlatformData gets platform data for a profile
func (c *Client) getPlatformData(profileID int64) (map[string]string, error) {
	rows, err := c.db.Query(
		`SELECT data_key, data_value FROM platform_data WHERE profile_id = ?`,
		profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query platform data: %w", err)
	}
	defer rows.Close()

	platformData := make(map[string]string)
	for rows.Next() {
		var key, value string
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan platform data: %w", err)
		}
		platformData[key] = value
	}

	return platformData, nil
}

// AddUserFeedback adds user feedback for a profile
func (c *Client) AddUserFeedback(profileID int64, feedbackType, comment string) error {
	_, err := c.db.Exec(
		`INSERT INTO user_feedback (profile_id, feedback_type, comment) VALUES (?, ?, ?)`,
		profileID, feedbackType, comment,
	)
	if err != nil {
		return fmt.Errorf("failed to add user feedback: %w", err)
	}

	return nil
}

// GetPopularSearches gets the most popular searches
func (c *Client) GetPopularSearches(limit int) ([]string, error) {
	rows, err := c.db.Query(
		`SELECT query, COUNT(*) as count
		FROM search_history
		GROUP BY query
		ORDER BY count DESC
		LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular searches: %w", err)
	}
	defer rows.Close()

	var searches []string
	for rows.Next() {
		var query string
		var count int
		err := rows.Scan(&query, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search: %w", err)
		}
		searches = append(searches, query)
	}

	return searches, nil
}
