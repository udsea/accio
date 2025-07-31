package intersection

import (
	"testing"

	"github.com/accio/internal/output"
)

func TestAnalyzeResults(t *testing.T) {
	// Create test data
	allResults := map[string][]output.Result{
		"johndoe": {
			{Site: "GitHub", URL: "https://github.com/johndoe", Exists: true},
			{Site: "Twitter", URL: "https://twitter.com/johndoe", Exists: true},
			{Site: "Facebook", URL: "https://facebook.com/johndoe", Exists: false},
		},
		"john.doe": {
			{Site: "GitHub", URL: "https://github.com/john.doe", Exists: false},
			{Site: "Twitter", URL: "https://twitter.com/john.doe", Exists: true},
			{Site: "Facebook", URL: "https://facebook.com/john.doe", Exists: true},
		},
		"jdoe": {
			{Site: "GitHub", URL: "https://github.com/jdoe", Exists: true},
			{Site: "Twitter", URL: "https://twitter.com/jdoe", Exists: false},
			{Site: "Facebook", URL: "https://facebook.com/jdoe", Exists: false},
		},
	}

	// Analyze results
	result := AnalyzeResults(allResults)

	// Check that we have the correct number of matches
	if len(result.Matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(result.Matches))
	}

	// Check that the total checked is correct
	if result.TotalChecked != 3 {
		t.Errorf("Expected TotalChecked to be 3, got %d", result.TotalChecked)
	}

	// Check that the total found is correct (2 + 2 + 1 = 5)
	if result.TotalFound != 5 {
		t.Errorf("Expected TotalFound to be 5, got %d", result.TotalFound)
	}

	// Check that matches are sorted by confidence (descending)
	for i := 1; i < len(result.Matches); i++ {
		if result.Matches[i-1].Confidence < result.Matches[i].Confidence {
			t.Errorf("Expected matches to be sorted by confidence (descending)")
		}
	}
}

func TestFindCommonProfiles(t *testing.T) {
	// Create test matches
	matches := []ProfileMatch{
		{
			Username: "johndoe",
			Results: []output.Result{
				{Site: "GitHub", URL: "https://github.com/johndoe", Exists: true},
				{Site: "Twitter", URL: "https://twitter.com/johndoe", Exists: true},
				{Site: "Reddit", URL: "https://reddit.com/user/johndoe", Exists: true},
			},
			MatchCount: 3,
			Confidence: 0.7,
		},
		{
			Username: "john.doe",
			Results: []output.Result{
				{Site: "GitHub", URL: "https://github.com/john.doe", Exists: true},
				{Site: "Flickr", URL: "https://flickr.com/john.doe", Exists: true},
			},
			MatchCount: 2,
			Confidence: 0.4,
		},
		{
			Username: "jdoe",
			Results: []output.Result{
				{Site: "Reddit", URL: "https://reddit.com/user/jdoe", Exists: true},
			},
			MatchCount: 1,
			Confidence: 0.1,
		},
	}

	// Find common profiles
	commonProfiles := FindCommonProfiles(matches)

	// Only the first profile should be common (appears on GitHub and Twitter)
	if len(commonProfiles) != 1 {
		t.Errorf("Expected 1 common profile, got %d", len(commonProfiles))
	}

	if len(commonProfiles) > 0 && commonProfiles[0].Username != "johndoe" {
		t.Errorf("Expected common profile to be johndoe, got %s", commonProfiles[0].Username)
	}
}

func TestFindProfilesByPlatforms(t *testing.T) {
	// Create test matches
	matches := []ProfileMatch{
		{
			Username: "johndoe",
			Results: []output.Result{
				{Site: "GitHub", URL: "https://github.com/johndoe", Exists: true},
				{Site: "Twitter", URL: "https://twitter.com/johndoe", Exists: true},
				{Site: "Reddit", URL: "https://reddit.com/user/johndoe", Exists: true},
			},
			MatchCount: 3,
			Confidence: 0.7,
		},
		{
			Username: "john.doe",
			Results: []output.Result{
				{Site: "GitHub", URL: "https://github.com/john.doe", Exists: true},
				{Site: "Flickr", URL: "https://flickr.com/john.doe", Exists: true},
			},
			MatchCount: 2,
			Confidence: 0.4,
		},
	}

	// Test with no platforms (should return all matches)
	filteredProfiles := FindProfilesByPlatforms(matches, []string{})
	if len(filteredProfiles) != 2 {
		t.Errorf("Expected 2 profiles with no platform filter, got %d", len(filteredProfiles))
	}

	// Test with GitHub platform (both profiles should match)
	filteredProfiles = FindProfilesByPlatforms(matches, []string{"GitHub"})
	if len(filteredProfiles) != 2 {
		t.Errorf("Expected 2 profiles with GitHub filter, got %d", len(filteredProfiles))
	}

	// Test with GitHub and Twitter platforms (only johndoe should match)
	filteredProfiles = FindProfilesByPlatforms(matches, []string{"GitHub", "Twitter"})
	if len(filteredProfiles) != 1 {
		t.Errorf("Expected 1 profile with GitHub and Twitter filter, got %d", len(filteredProfiles))
	}

	if len(filteredProfiles) > 0 && filteredProfiles[0].Username != "johndoe" {
		t.Errorf("Expected filtered profile to be johndoe, got %s", filteredProfiles[0].Username)
	}

	// Test with non-existent platform (no profiles should match)
	filteredProfiles = FindProfilesByPlatforms(matches, []string{"LinkedIn"})
	if len(filteredProfiles) != 0 {
		t.Errorf("Expected 0 profiles with LinkedIn filter, got %d", len(filteredProfiles))
	}
}
