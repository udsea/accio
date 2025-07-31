package intersection

import (
	"sort"
	"strings"

	"github.com/accio/internal/output"
)

// ProfileMatch represents a username match across platforms
type ProfileMatch struct {
	Username    string          // The username that was matched
	Results     []output.Result // The results for this username across platforms
	MatchCount  int             // Number of platforms where this username was found
	Confidence  float64         // Confidence score (0.0-1.0) that these profiles belong to the same person
	NameMatches map[string]bool // Map of platforms where the real name matches
}

// AnalysisResult contains the results of an intersection analysis
type AnalysisResult struct {
	Matches        []ProfileMatch // All matches found
	TotalChecked   int            // Total number of username variations checked
	TotalFound     int            // Total number of profiles found
	UniqueProfiles int            // Number of likely unique profiles
}

// AnalyzeResults performs intersection analysis on results from multiple username checks
func AnalyzeResults(allResults map[string][]output.Result) AnalysisResult {
	// Group results by username
	matches := groupResultsByUsername(allResults)

	// Calculate confidence scores
	calculateConfidenceScores(matches)

	// Sort matches by confidence score (descending)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Confidence > matches[j].Confidence
	})

	// Count totals
	totalChecked := len(allResults)
	totalFound := 0
	for _, match := range matches {
		totalFound += match.MatchCount
	}

	// Count unique profiles (those with high confidence)
	uniqueProfiles := 0
	for _, match := range matches {
		if match.Confidence >= 0.7 {
			uniqueProfiles++
		}
	}

	return AnalysisResult{
		Matches:        matches,
		TotalChecked:   totalChecked,
		TotalFound:     totalFound,
		UniqueProfiles: uniqueProfiles,
	}
}

// groupResultsByUsername groups results by username
func groupResultsByUsername(allResults map[string][]output.Result) []ProfileMatch {
	matches := []ProfileMatch{}

	for username, results := range allResults {
		// Filter to only include found profiles
		foundResults := []output.Result{}
		for _, result := range results {
			if result.Exists {
				foundResults = append(foundResults, result)
			}
		}

		// Only include usernames with at least one match
		if len(foundResults) > 0 {
			matches = append(matches, ProfileMatch{
				Username:    username,
				Results:     foundResults,
				MatchCount:  len(foundResults),
				Confidence:  0.0, // Will be calculated later
				NameMatches: make(map[string]bool),
			})
		}
	}

	return matches
}

// calculateConfidenceScores calculates confidence scores for each match
func calculateConfidenceScores(matches []ProfileMatch) {
	// Base confidence on number of matches
	// More matches = higher confidence
	for i := range matches {
		// Start with a base confidence based on match count
		// 1 match = 0.1, 2 matches = 0.3, 3 matches = 0.5, 4+ matches = 0.7+
		switch matches[i].MatchCount {
		case 1:
			matches[i].Confidence = 0.1
		case 2:
			matches[i].Confidence = 0.3
		case 3:
			matches[i].Confidence = 0.5
		default:
			// For 4+ matches, scale up to 0.9 max
			matches[i].Confidence = 0.7 + float64(matches[i].MatchCount-4)*0.05
			if matches[i].Confidence > 0.9 {
				matches[i].Confidence = 0.9
			}
		}

		// Adjust confidence based on platform importance
		// Some platforms are more likely to be unique identifiers
		highValuePlatforms := map[string]float64{
			"GitHub":    0.1,
			"Twitter":   0.1,
			"LinkedIn":  0.15,
			"Facebook":  0.1,
			"Instagram": 0.1,
		}

		for _, result := range matches[i].Results {
			if bonus, ok := highValuePlatforms[result.Site]; ok {
				// Add bonus, but don't exceed 1.0
				matches[i].Confidence += bonus
				if matches[i].Confidence > 1.0 {
					matches[i].Confidence = 1.0
				}
			}
		}
	}
}

// FindCommonProfiles finds profiles that appear on multiple high-value platforms
func FindCommonProfiles(matches []ProfileMatch) []ProfileMatch {
	highValuePlatforms := []string{"GitHub", "Twitter", "LinkedIn", "Facebook", "Instagram"}

	commonProfiles := []ProfileMatch{}
	for _, match := range matches {
		// Count how many high-value platforms this username appears on
		platformCount := 0
		for _, result := range match.Results {
			for _, platform := range highValuePlatforms {
				if result.Site == platform {
					platformCount++
					break
				}
			}
		}

		// If the username appears on at least 2 high-value platforms, consider it common
		if platformCount >= 2 {
			commonProfiles = append(commonProfiles, match)
		}
	}

	return commonProfiles
}

// FindProfilesByPlatforms finds profiles that appear on all specified platforms
func FindProfilesByPlatforms(matches []ProfileMatch, platforms []string) []ProfileMatch {
	if len(platforms) == 0 {
		return matches
	}

	// Convert platforms to lowercase for case-insensitive comparison
	lowerPlatforms := make([]string, len(platforms))
	for i, platform := range platforms {
		lowerPlatforms[i] = strings.ToLower(platform)
	}

	filteredProfiles := []ProfileMatch{}
	for _, match := range matches {
		// Check if this profile appears on all specified platforms
		platformMatches := make(map[string]bool)
		for _, result := range match.Results {
			platformMatches[strings.ToLower(result.Site)] = true
		}

		allFound := true
		for _, platform := range lowerPlatforms {
			if !platformMatches[platform] {
				allFound = false
				break
			}
		}

		if allFound {
			filteredProfiles = append(filteredProfiles, match)
		}
	}

	return filteredProfiles
}
