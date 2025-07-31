package model

import (
	"testing"
)

func TestNewProfile(t *testing.T) {
	// Create a new profile
	profile := NewProfile(
		"John Doe",
		"johndoe",
		"Twitter",
		"https://twitter.com/johndoe",
		"https://pbs.twimg.com/profile_images/123456789/johndoe.jpg",
		"Software Engineer",
		true,
		1000,
	)

	// Check that the profile was created correctly
	if profile.RealName != "John Doe" {
		t.Errorf("Expected RealName to be 'John Doe', got '%s'", profile.RealName)
	}

	if profile.Username != "johndoe" {
		t.Errorf("Expected Username to be 'johndoe', got '%s'", profile.Username)
	}

	if profile.Platform != "Twitter" {
		t.Errorf("Expected Platform to be 'Twitter', got '%s'", profile.Platform)
	}

	if profile.ProfileURL != "https://twitter.com/johndoe" {
		t.Errorf("Expected ProfileURL to be 'https://twitter.com/johndoe', got '%s'", profile.ProfileURL)
	}

	if profile.ImageURL != "https://pbs.twimg.com/profile_images/123456789/johndoe.jpg" {
		t.Errorf("Expected ImageURL to be 'https://pbs.twimg.com/profile_images/123456789/johndoe.jpg', got '%s'", profile.ImageURL)
	}

	if profile.Bio != "Software Engineer" {
		t.Errorf("Expected Bio to be 'Software Engineer', got '%s'", profile.Bio)
	}

	if !profile.Verified {
		t.Errorf("Expected Verified to be true, got %v", profile.Verified)
	}

	if profile.FollowerCount != 1000 {
		t.Errorf("Expected FollowerCount to be 1000, got %d", profile.FollowerCount)
	}

	// Check that LastUpdated is set
	if profile.LastUpdated.IsZero() {
		t.Errorf("Expected LastUpdated to be set, got zero time")
	}
}

func TestProfileAddNamePart(t *testing.T) {
	// Create a new profile
	profile := NewProfile(
		"John Doe",
		"johndoe",
		"Twitter",
		"https://twitter.com/johndoe",
		"https://pbs.twimg.com/profile_images/123456789/johndoe.jpg",
		"Software Engineer",
		true,
		1000,
	)

	// Add name parts
	profile.AddNamePart("John", "first")
	profile.AddNamePart("Doe", "last")

	// Check that name parts were added correctly
	if len(profile.NameParts) != 2 {
		t.Errorf("Expected 2 name parts, got %d", len(profile.NameParts))
	}

	if profile.NameParts[0].NamePart != "John" {
		t.Errorf("Expected first name part to be 'John', got '%s'", profile.NameParts[0].NamePart)
	}

	if profile.NameParts[0].PartType != "first" {
		t.Errorf("Expected first name part type to be 'first', got '%s'", profile.NameParts[0].PartType)
	}

	if profile.NameParts[1].NamePart != "Doe" {
		t.Errorf("Expected second name part to be 'Doe', got '%s'", profile.NameParts[1].NamePart)
	}

	if profile.NameParts[1].PartType != "last" {
		t.Errorf("Expected second name part type to be 'last', got '%s'", profile.NameParts[1].PartType)
	}
}

func TestProfileAddAlias(t *testing.T) {
	// Create a new profile
	profile := NewProfile(
		"John Doe",
		"johndoe",
		"Twitter",
		"https://twitter.com/johndoe",
		"https://pbs.twimg.com/profile_images/123456789/johndoe.jpg",
		"Software Engineer",
		true,
		1000,
	)

	// Add aliases
	profile.AddAlias("jdoe")
	profile.AddAlias("john.doe")

	// Check that aliases were added correctly
	if len(profile.Aliases) != 2 {
		t.Errorf("Expected 2 aliases, got %d", len(profile.Aliases))
	}

	if profile.Aliases[0].Alias != "jdoe" {
		t.Errorf("Expected first alias to be 'jdoe', got '%s'", profile.Aliases[0].Alias)
	}

	if profile.Aliases[1].Alias != "john.doe" {
		t.Errorf("Expected second alias to be 'john.doe', got '%s'", profile.Aliases[1].Alias)
	}
}

func TestProfileAddPlatformData(t *testing.T) {
	// Create a new profile
	profile := NewProfile(
		"John Doe",
		"johndoe",
		"Twitter",
		"https://twitter.com/johndoe",
		"https://pbs.twimg.com/profile_images/123456789/johndoe.jpg",
		"Software Engineer",
		true,
		1000,
	)

	// Add platform data
	profile.AddPlatformData("user_id", "123456789")
	profile.AddPlatformData("created_at", "2020-01-01")

	// Check that platform data was added correctly
	if len(profile.PlatformData) != 2 {
		t.Errorf("Expected 2 platform data entries, got %d", len(profile.PlatformData))
	}

	if profile.PlatformData[0].DataKey != "user_id" {
		t.Errorf("Expected first platform data key to be 'user_id', got '%s'", profile.PlatformData[0].DataKey)
	}

	if profile.PlatformData[0].DataValue != "123456789" {
		t.Errorf("Expected first platform data value to be '123456789', got '%s'", profile.PlatformData[0].DataValue)
	}

	if profile.PlatformData[1].DataKey != "created_at" {
		t.Errorf("Expected second platform data key to be 'created_at', got '%s'", profile.PlatformData[1].DataKey)
	}

	if profile.PlatformData[1].DataValue != "2020-01-01" {
		t.Errorf("Expected second platform data value to be '2020-01-01', got '%s'", profile.PlatformData[1].DataValue)
	}
}

func TestProfileGetPlatformData(t *testing.T) {
	// Create a new profile
	profile := NewProfile(
		"John Doe",
		"johndoe",
		"Twitter",
		"https://twitter.com/johndoe",
		"https://pbs.twimg.com/profile_images/123456789/johndoe.jpg",
		"Software Engineer",
		true,
		1000,
	)

	// Add platform data
	profile.AddPlatformData("user_id", "123456789")
	profile.AddPlatformData("created_at", "2020-01-01")

	// Get platform data
	value, ok := profile.GetPlatformData("user_id")
	if !ok {
		t.Errorf("Expected to find platform data for key 'user_id'")
	}
	if value != "123456789" {
		t.Errorf("Expected platform data value to be '123456789', got '%s'", value)
	}

	// Get non-existent platform data
	value, ok = profile.GetPlatformData("non_existent")
	if ok {
		t.Errorf("Expected not to find platform data for key 'non_existent'")
	}
	if value != "" {
		t.Errorf("Expected platform data value to be empty, got '%s'", value)
	}
}

func TestProfileGetPlatformDataMap(t *testing.T) {
	// Create a new profile
	profile := NewProfile(
		"John Doe",
		"johndoe",
		"Twitter",
		"https://twitter.com/johndoe",
		"https://pbs.twimg.com/profile_images/123456789/johndoe.jpg",
		"Software Engineer",
		true,
		1000,
	)

	// Add platform data
	profile.AddPlatformData("user_id", "123456789")
	profile.AddPlatformData("created_at", "2020-01-01")

	// Get platform data map
	dataMap := profile.GetPlatformDataMap()

	// Check that platform data map was created correctly
	if len(dataMap) != 2 {
		t.Errorf("Expected 2 platform data entries, got %d", len(dataMap))
	}

	if dataMap["user_id"] != "123456789" {
		t.Errorf("Expected platform data value for 'user_id' to be '123456789', got '%s'", dataMap["user_id"])
	}

	if dataMap["created_at"] != "2020-01-01" {
		t.Errorf("Expected platform data value for 'created_at' to be '2020-01-01', got '%s'", dataMap["created_at"])
	}
}

func TestProfileGetFullName(t *testing.T) {
	// Create a new profile with real name
	profile1 := NewProfile(
		"John Doe",
		"johndoe",
		"Twitter",
		"https://twitter.com/johndoe",
		"https://pbs.twimg.com/profile_images/123456789/johndoe.jpg",
		"Software Engineer",
		true,
		1000,
	)

	// Add name parts
	profile1.AddNamePart("John", "first")
	profile1.AddNamePart("Doe", "last")

	// Check that GetFullName returns the real name
	if profile1.GetFullName() != "John Doe" {
		t.Errorf("Expected GetFullName to return 'John Doe', got '%s'", profile1.GetFullName())
	}

	// Create a new profile without real name
	profile2 := NewProfile(
		"",
		"johndoe",
		"Twitter",
		"https://twitter.com/johndoe",
		"https://pbs.twimg.com/profile_images/123456789/johndoe.jpg",
		"Software Engineer",
		true,
		1000,
	)

	// Add name parts
	profile2.AddNamePart("John", "first")
	profile2.AddNamePart("Doe", "last")

	// Check that GetFullName returns the constructed name
	if profile2.GetFullName() != "John Doe" {
		t.Errorf("Expected GetFullName to return 'John Doe', got '%s'", profile2.GetFullName())
	}

	// Create a new profile without real name or name parts
	profile3 := NewProfile(
		"",
		"johndoe",
		"Twitter",
		"https://twitter.com/johndoe",
		"https://pbs.twimg.com/profile_images/123456789/johndoe.jpg",
		"Software Engineer",
		true,
		1000,
	)

	// Check that GetFullName returns the username
	if profile3.GetFullName() != "johndoe" {
		t.Errorf("Expected GetFullName to return 'johndoe', got '%s'", profile3.GetFullName())
	}
}

func TestProfileUniqueKey(t *testing.T) {
	// Create a new profile
	profile := NewProfile(
		"John Doe",
		"johndoe",
		"Twitter",
		"https://twitter.com/johndoe",
		"https://pbs.twimg.com/profile_images/123456789/johndoe.jpg",
		"Software Engineer",
		true,
		1000,
	)

	// Check that UniqueKey returns the expected value
	if profile.UniqueKey() != "Twitter:johndoe" {
		t.Errorf("Expected UniqueKey to return 'Twitter:johndoe', got '%s'", profile.UniqueKey())
	}
}
