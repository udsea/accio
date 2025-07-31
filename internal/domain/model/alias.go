package model

// Alias represents an alternative name for a profile
type Alias struct {
	ID        uint   `gorm:"primaryKey"`
	ProfileID uint   `gorm:"index"`
	Alias     string `gorm:"index"`
}

// NewAlias creates a new alias entity
func NewAlias(profileID uint, alias string) *Alias {
	return &Alias{
		ProfileID: profileID,
		Alias:     alias,
	}
}
