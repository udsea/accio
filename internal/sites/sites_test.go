package sites

import (
	"testing"
)

func TestGetSites(t *testing.T) {
	sites := GetSites()

	// Check that we have sites
	if len(sites) == 0 {
		t.Error("Expected sites to be returned, got empty slice")
	}

	// Check that each site has required fields
	for i, site := range sites {
		if site.Name == "" {
			t.Errorf("Site at index %d has empty name", i)
		}
		if site.URL == "" {
			t.Errorf("Site at index %d has empty URL", i)
		}
		if site.URLFormat == "" {
			t.Errorf("Site at index %d has empty URLFormat", i)
		}
		if site.CheckMethod == "" {
			t.Errorf("Site at index %d has empty CheckMethod", i)
		}
	}
}

func TestGetSiteByName(t *testing.T) {
	// Test getting an existing site
	site, found := GetSiteByName("GitHub")
	if !found {
		t.Error("Expected to find GitHub site, but it was not found")
	}
	if site.Name != "GitHub" {
		t.Errorf("Expected site name to be GitHub, got %s", site.Name)
	}

	// Test getting a non-existent site
	_, found = GetSiteByName("NonExistentSite")
	if found {
		t.Error("Expected not to find NonExistentSite, but it was found")
	}
}
