package matcher

import (
	"testing"
)

func TestParseFullName(t *testing.T) {
	testCases := []struct {
		name     string
		fullName string
		expected NameInfo
	}{
		{
			name:     "Empty name",
			fullName: "",
			expected: NameInfo{},
		},
		{
			name:     "First name only",
			fullName: "John",
			expected: NameInfo{FirstName: "John"},
		},
		{
			name:     "First and last name",
			fullName: "John Doe",
			expected: NameInfo{FirstName: "John", LastName: "Doe"},
		},
		{
			name:     "Full name with middle name",
			fullName: "John Michael Doe",
			expected: NameInfo{FirstName: "John", MiddleName: "Michael", LastName: "Doe"},
		},
		{
			name:     "Full name with multiple middle names",
			fullName: "John Michael James Doe",
			expected: NameInfo{FirstName: "John", MiddleName: "Michael James", LastName: "Doe"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseFullName(tc.fullName)

			if result.FirstName != tc.expected.FirstName {
				t.Errorf("Expected FirstName to be %s, got %s", tc.expected.FirstName, result.FirstName)
			}

			if result.MiddleName != tc.expected.MiddleName {
				t.Errorf("Expected MiddleName to be %s, got %s", tc.expected.MiddleName, result.MiddleName)
			}

			if result.LastName != tc.expected.LastName {
				t.Errorf("Expected LastName to be %s, got %s", tc.expected.LastName, result.LastName)
			}
		})
	}
}

func TestGenerateUsernames(t *testing.T) {
	nameInfo := NameInfo{
		FirstName: "John",
		LastName:  "Doe",
		BirthYear: 1990,
	}

	usernames := nameInfo.GenerateUsernames()

	// Check that we have some usernames
	if len(usernames) == 0 {
		t.Error("Expected usernames to be generated, got empty slice")
	}

	// Check for some expected variations
	expectedVariations := []string{
		"john",
		"doe",
		"johndoe",
		"john.doe",
		"john_doe",
		"jdoe",
		"john1990",
		"john90",
	}

	for _, expected := range expectedVariations {
		found := false
		for _, username := range usernames {
			if username == expected {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected to find username variation %s, but it was not generated", expected)
		}
	}
}

func TestGenerateCommonUsernames(t *testing.T) {
	nameInfo := NameInfo{
		FirstName: "John",
		LastName:  "Doe",
		BirthYear: 1990,
	}

	usernames := nameInfo.GenerateCommonUsernames()

	// Check that we have some usernames
	if len(usernames) == 0 {
		t.Error("Expected usernames to be generated, got empty slice")
	}

	// Check that common usernames is a subset of all usernames
	allUsernames := nameInfo.GenerateUsernames()
	if len(usernames) > len(allUsernames) {
		t.Error("Expected common usernames to be a subset of all usernames")
	}
}

func TestGenerateAdvancedUsernames(t *testing.T) {
	nameInfo := NameInfo{
		FirstName: "John",
		LastName:  "Doe",
	}

	usernames := nameInfo.GenerateAdvancedUsernames()

	// Check that we have some usernames
	if len(usernames) == 0 {
		t.Error("Expected usernames to be generated, got empty slice")
	}

	// Check for some expected advanced variations
	expectedVariations := []string{
		"doejohn",
		"doe.john",
		"j0hn", // leetspeak
		"d03",  // leetspeak
		"realjohn",
		"official_doe",
	}

	for _, expected := range expectedVariations {
		found := false
		for _, username := range usernames {
			if username == expected {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected to find username variation %s, but it was not generated", expected)
		}
	}
}
