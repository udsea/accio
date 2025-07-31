package model

// NamePart represents a part of a person's name
type NamePart struct {
	ID        uint   `gorm:"primaryKey"`
	ProfileID uint   `gorm:"index"`
	NamePart  string `gorm:"index"`
	PartType  string // 'first', 'middle', 'last', 'nickname'
}

// NewNamePart creates a new name part entity
func NewNamePart(profileID uint, namePart, partType string) *NamePart {
	return &NamePart{
		ProfileID: profileID,
		NamePart:  namePart,
		PartType:  partType,
	}
}
