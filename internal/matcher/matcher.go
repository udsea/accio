package matcher

import (
	"fmt"
	"strings"
	"unicode"
)

// NameInfo contains information about a person's name
type NameInfo struct {
	FirstName  string
	MiddleName string
	LastName   string
	BirthYear  int // Optional birth year for username variations
}

// NewNameInfo creates a new NameInfo instance
func NewNameInfo(firstName, middleName, lastName string, birthYear int) *NameInfo {
	return &NameInfo{
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		BirthYear:  birthYear,
	}
}

// ParseFullName parses a full name into first, middle, and last names
func ParseFullName(fullName string) *NameInfo {
	parts := strings.Fields(fullName)

	switch len(parts) {
	case 0:
		return &NameInfo{}
	case 1:
		return &NameInfo{FirstName: parts[0]}
	case 2:
		return &NameInfo{
			FirstName: parts[0],
			LastName:  parts[1],
		}
	default:
		return &NameInfo{
			FirstName:  parts[0],
			MiddleName: strings.Join(parts[1:len(parts)-1], " "),
			LastName:   parts[len(parts)-1],
		}
	}
}

// GenerateUsernames generates possible username variations based on the name
func (n *NameInfo) GenerateUsernames() []string {
	if n.FirstName == "" {
		return []string{}
	}

	variations := []string{}

	// Basic variations
	variations = append(variations, strings.ToLower(n.FirstName))

	if n.LastName != "" {
		variations = append(variations,
			strings.ToLower(n.LastName),
			strings.ToLower(n.FirstName+n.LastName),
			strings.ToLower(n.FirstName+"."+n.LastName),
			strings.ToLower(n.FirstName+"_"+n.LastName),
			strings.ToLower(string(n.FirstName[0])+n.LastName),
		)
	}

	// Add middle initial if available
	if n.MiddleName != "" {
		middleInitial := string(n.MiddleName[0])
		variations = append(variations,
			strings.ToLower(n.FirstName+middleInitial+n.LastName),
			strings.ToLower(n.FirstName+"."+middleInitial+"."+n.LastName),
			strings.ToLower(n.FirstName+"_"+middleInitial+"_"+n.LastName),
		)
	}

	// Add birth year variations if available
	if n.BirthYear > 0 {
		yearStr := fmt.Sprintf("%d", n.BirthYear)
		shortYearStr := yearStr[2:] // Last two digits

		// Create year variations
		yearVariations := []string{}
		for _, username := range variations {
			yearVariations = append(yearVariations,
				username+yearStr,
				username+shortYearStr,
				username+"_"+yearStr,
				username+"_"+shortYearStr,
			)
		}

		variations = append(variations, yearVariations...)
	}

	// Remove duplicates
	return removeDuplicates(variations)
}

// GenerateCommonUsernames generates the most common username patterns
// This is a subset of GenerateUsernames for quicker searches
func (n *NameInfo) GenerateCommonUsernames() []string {
	if n.FirstName == "" {
		return []string{}
	}

	variations := []string{}

	// Most common variations
	firstName := strings.ToLower(n.FirstName)
	variations = append(variations, firstName)

	if n.LastName != "" {
		lastName := strings.ToLower(n.LastName)
		variations = append(variations,
			lastName,
			firstName+lastName,
			firstName+"."+lastName,
			string(firstName[0])+lastName,
		)
	}

	// Add birth year to first name if available
	if n.BirthYear > 0 {
		yearStr := fmt.Sprintf("%d", n.BirthYear)
		shortYearStr := yearStr[2:] // Last two digits

		variations = append(variations,
			firstName+yearStr,
			firstName+shortYearStr,
		)
	}

	return removeDuplicates(variations)
}

// GenerateAdvancedUsernames generates more sophisticated username variations
func (n *NameInfo) GenerateAdvancedUsernames() []string {
	basic := n.GenerateUsernames()
	advanced := []string{}

	// Add leetspeak variations
	for _, username := range basic {
		advanced = append(advanced, toLeetspeak(username))
	}

	// Add reversed names
	if n.LastName != "" {
		advanced = append(advanced,
			strings.ToLower(n.LastName+n.FirstName),
			strings.ToLower(n.LastName+"."+n.FirstName),
			strings.ToLower(n.LastName+"_"+n.FirstName),
		)
	}

	// Add common prefixes/suffixes
	commonAffixes := []string{"the", "real", "official", "its", "im", "mr", "ms", "dr"}
	for _, affix := range commonAffixes {
		if n.FirstName != "" {
			advanced = append(advanced,
				strings.ToLower(affix+n.FirstName),
				strings.ToLower(affix+"_"+n.FirstName),
			)
		}

		if n.LastName != "" {
			advanced = append(advanced,
				strings.ToLower(affix+n.LastName),
				strings.ToLower(affix+"_"+n.LastName),
			)
		}
	}

	return append(basic, removeDuplicates(advanced)...)
}

// toLeetspeak converts a string to basic leetspeak
func toLeetspeak(s string) string {
	leetMap := map[rune]string{
		'a': "4",
		'e': "3",
		'i': "1",
		'o': "0",
		's': "5",
		't': "7",
	}

	var result strings.Builder
	for _, char := range s {
		if leet, ok := leetMap[unicode.ToLower(char)]; ok {
			result.WriteString(leet)
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
