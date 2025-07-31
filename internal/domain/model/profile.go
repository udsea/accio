package model

import (
	"time"
)

// Profile represents a social media profile as a domain entity
type Profile struct {
	ID            uint   `gorm:"primaryKey"`
	RealName      string `gorm:"index"`
	Username      string `gorm:"index"`
	Platform      string `gorm:"index"`
	ProfileURL    string
	ImageURL      string
	Verified      bool
	FollowerCount int64
	Bio           string
	LastUpdated   time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relationships
	NameParts    []NamePart     `gorm:"foreignKey:ProfileID"`
	Aliases      []Alias        `gorm:"foreignKey:ProfileID"`
	PlatformData []PlatformData `gorm:"foreignKey:ProfileID"`
}

// NewProfile creates a new profile entity
func NewProfile(realName, username, platform, profileURL, imageURL, bio string, verified bool, followerCount int64) *Profile {
	return &Profile{
		RealName:      realName,
		Username:      username,
		Platform:      platform,
		ProfileURL:    profileURL,
		ImageURL:      imageURL,
		Bio:           bio,
		Verified:      verified,
		FollowerCount: followerCount,
		LastUpdated:   time.Now(),
	}
}

// AddNamePart adds a name part to the profile
func (p *Profile) AddNamePart(namePart, partType string) {
	p.NameParts = append(p.NameParts, NamePart{
		NamePart:  namePart,
		PartType:  partType,
		ProfileID: p.ID,
	})
}

// AddAlias adds an alias to the profile
func (p *Profile) AddAlias(alias string) {
	p.Aliases = append(p.Aliases, Alias{
		Alias:     alias,
		ProfileID: p.ID,
	})
}

// AddPlatformData adds platform-specific data to the profile
func (p *Profile) AddPlatformData(key, value string) {
	p.PlatformData = append(p.PlatformData, PlatformData{
		DataKey:   key,
		DataValue: value,
		ProfileID: p.ID,
	})
}

// GetPlatformData gets platform-specific data from the profile
func (p *Profile) GetPlatformData(key string) (string, bool) {
	for _, data := range p.PlatformData {
		if data.DataKey == key {
			return data.DataValue, true
		}
	}
	return "", false
}

// GetPlatformDataMap returns all platform data as a map
func (p *Profile) GetPlatformDataMap() map[string]string {
	result := make(map[string]string)
	for _, data := range p.PlatformData {
		result[data.DataKey] = data.DataValue
	}
	return result
}

// GetFirstName returns the first name of the profile
func (p *Profile) GetFirstName() string {
	for _, part := range p.NameParts {
		if part.PartType == "first" {
			return part.NamePart
		}
	}
	return ""
}

// GetLastName returns the last name of the profile
func (p *Profile) GetLastName() string {
	for _, part := range p.NameParts {
		if part.PartType == "last" {
			return part.NamePart
		}
	}
	return ""
}

// GetFullName returns the full name of the profile
func (p *Profile) GetFullName() string {
	if p.RealName != "" {
		return p.RealName
	}

	firstName := p.GetFirstName()
	lastName := p.GetLastName()

	if firstName != "" && lastName != "" {
		return firstName + " " + lastName
	}

	return p.Username
}

// UniqueKey returns a unique key for the profile
func (p *Profile) UniqueKey() string {
	return p.Platform + ":" + p.Username
}
