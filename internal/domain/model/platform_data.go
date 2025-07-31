package model

// PlatformData represents platform-specific data for a profile
type PlatformData struct {
	ID        uint   `gorm:"primaryKey"`
	ProfileID uint   `gorm:"index"`
	DataKey   string `gorm:"index"`
	DataValue string
}

// NewPlatformData creates a new platform data entity
func NewPlatformData(profileID uint, key, value string) *PlatformData {
	return &PlatformData{
		ProfileID: profileID,
		DataKey:   key,
		DataValue: value,
	}
}
