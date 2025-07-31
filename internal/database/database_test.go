package database

import (
	"os"
	"testing"
	"time"
)

// TestNewClient tests creating a new database client
func TestNewClient(t *testing.T) {
	// Skip if no database URL is provided
	dbURL := os.Getenv("TURSO_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TURSO_DATABASE_URL not set, skipping database tests")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create database client: %v", err)
	}
	defer client.Close()
}

// TestAddAndGetProfile tests adding and retrieving a profile
func TestAddAndGetProfile(t *testing.T) {
	// Skip if no database URL is provided
	dbURL := os.Getenv("TURSO_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TURSO_DATABASE_URL not set, skipping database tests")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create database client: %v", err)
	}
	defer client.Close()

	// Create a test profile
	profile := &Profile{
		RealName:      "Test User",
		Username:      "testuser_" + time.Now().Format("20060102150405"),
		Platform:      "TestPlatform",
		ProfileURL:    "https://example.com/testuser",
		ImageURL:      "https://example.com/testuser.jpg",
		Verified:      true,
		FollowerCount: 1000,
		Bio:           "This is a test user",
		NameParts: []NamePart{
			{NamePart: "Test", PartType: "first"},
			{NamePart: "User", PartType: "last"},
		},
		Aliases: []string{"tester", "testaccount"},
		PlatformData: map[string]string{
			"user_id": "12345",
			"joined":  "2020-01-01",
		},
	}

	// Add the profile
	err = client.AddProfile(profile)
	if err != nil {
		t.Fatalf("Failed to add profile: %v", err)
	}

	// Get the profile
	retrievedProfile, err := client.GetProfileByUsername(profile.Username, profile.Platform)
	if err != nil {
		t.Fatalf("Failed to get profile: %v", err)
	}

	// Check that the profile was retrieved correctly
	if retrievedProfile == nil {
		t.Fatalf("Profile not found")
	}

	if retrievedProfile.RealName != profile.RealName {
		t.Errorf("Expected RealName to be %s, got %s", profile.RealName, retrievedProfile.RealName)
	}

	if retrievedProfile.Username != profile.Username {
		t.Errorf("Expected Username to be %s, got %s", profile.Username, retrievedProfile.Username)
	}

	if retrievedProfile.Platform != profile.Platform {
		t.Errorf("Expected Platform to be %s, got %s", profile.Platform, retrievedProfile.Platform)
	}

	if retrievedProfile.ProfileURL != profile.ProfileURL {
		t.Errorf("Expected ProfileURL to be %s, got %s", profile.ProfileURL, retrievedProfile.ProfileURL)
	}

	if retrievedProfile.ImageURL != profile.ImageURL {
		t.Errorf("Expected ImageURL to be %s, got %s", profile.ImageURL, retrievedProfile.ImageURL)
	}

	if retrievedProfile.Verified != profile.Verified {
		t.Errorf("Expected Verified to be %v, got %v", profile.Verified, retrievedProfile.Verified)
	}

	if retrievedProfile.FollowerCount != profile.FollowerCount {
		t.Errorf("Expected FollowerCount to be %d, got %d", profile.FollowerCount, retrievedProfile.FollowerCount)
	}

	if retrievedProfile.Bio != profile.Bio {
		t.Errorf("Expected Bio to be %s, got %s", profile.Bio, retrievedProfile.Bio)
	}

	// Check name parts
	if len(retrievedProfile.NameParts) != len(profile.NameParts) {
		t.Errorf("Expected %d name parts, got %d", len(profile.NameParts), len(retrievedProfile.NameParts))
	}

	// Check aliases
	if len(retrievedProfile.Aliases) != len(profile.Aliases) {
		t.Errorf("Expected %d aliases, got %d", len(profile.Aliases), len(retrievedProfile.Aliases))
	}

	// Check platform data
	if len(retrievedProfile.PlatformData) != len(profile.PlatformData) {
		t.Errorf("Expected %d platform data entries, got %d", len(profile.PlatformData), len(retrievedProfile.PlatformData))
	}
}

// TestSearchProfilesByName tests searching for profiles by name
func TestSearchProfilesByName(t *testing.T) {
	// Skip if no database URL is provided
	dbURL := os.Getenv("TURSO_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TURSO_DATABASE_URL not set, skipping database tests")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create database client: %v", err)
	}
	defer client.Close()

	// Create a unique test name
	testName := "SearchTest_" + time.Now().Format("20060102150405")

	// Create a test profile
	profile := &Profile{
		RealName:      testName + " User",
		Username:      "searchuser_" + time.Now().Format("20060102150405"),
		Platform:      "TestPlatform",
		ProfileURL:    "https://example.com/searchuser",
		ImageURL:      "https://example.com/searchuser.jpg",
		Verified:      true,
		FollowerCount: 1000,
		Bio:           "This is a test user for search",
		NameParts: []NamePart{
			{NamePart: testName, PartType: "first"},
			{NamePart: "User", PartType: "last"},
		},
		Aliases: []string{"searcher", "findme"},
		PlatformData: map[string]string{
			"user_id": "67890",
			"joined":  "2021-01-01",
		},
	}

	// Add the profile
	err = client.AddProfile(profile)
	if err != nil {
		t.Fatalf("Failed to add profile: %v", err)
	}

	// Search for the profile
	profiles, err := client.SearchProfilesByName(testName)
	if err != nil {
		t.Fatalf("Failed to search profiles: %v", err)
	}

	// Check that the profile was found
	if len(profiles) == 0 {
		t.Fatalf("No profiles found")
	}

	found := false
	for _, p := range profiles {
		if p.Username == profile.Username {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected to find profile with username %s", profile.Username)
	}
}
