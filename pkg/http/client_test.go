package http

import (
	"testing"
)

func TestFormatURL(t *testing.T) {
	testCases := []struct {
		name      string
		urlFormat string
		username  string
		expected  string
	}{
		{
			name:      "Simple replacement",
			urlFormat: "https://example.com/{}",
			username:  "testuser",
			expected:  "https://example.com/testuser",
		},
		{
			name:      "Multiple replacements",
			urlFormat: "https://{}.example.com/users/{}",
			username:  "testuser",
			expected:  "https://testuser.example.com/users/testuser",
		},
		{
			name:      "No placeholder",
			urlFormat: "https://example.com/users",
			username:  "testuser",
			expected:  "https://example.com/users",
		},
		{
			name:      "Empty username",
			urlFormat: "https://example.com/{}",
			username:  "",
			expected:  "https://example.com/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatURL(tc.urlFormat, tc.username)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	// Test with default values
	client := NewClient(10, false)
	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.UserAgent != "Accio/1.0" {
		t.Errorf("Expected UserAgent to be Accio/1.0, got %s", client.UserAgent)
	}

	if client.Verbose != false {
		t.Errorf("Expected Verbose to be false, got %v", client.Verbose)
	}

	// Test with custom values
	client = NewClient(20, true)
	if client.Verbose != true {
		t.Errorf("Expected Verbose to be true, got %v", client.Verbose)
	}
}
