package dto

// ProfileDTO represents a profile data transfer object
type ProfileDTO struct {
	ID            uint              `json:"id,omitempty"`
	RealName      string            `json:"real_name"`
	Username      string            `json:"username"`
	Platform      string            `json:"platform"`
	ProfileURL    string            `json:"profile_url"`
	ImageURL      string            `json:"image_url"`
	Bio           string            `json:"bio"`
	Verified      bool              `json:"verified"`
	FollowerCount int64             `json:"follower_count"`
	NameParts     []NamePartDTO     `json:"name_parts,omitempty"`
	Aliases       []string          `json:"aliases,omitempty"`
	PlatformData  map[string]string `json:"platform_data,omitempty"`
}

// NamePartDTO represents a name part data transfer object
type NamePartDTO struct {
	NamePart string `json:"name_part"`
	PartType string `json:"part_type"`
}
