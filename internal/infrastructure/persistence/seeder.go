package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/accio/internal/domain/model"
	"gorm.io/gorm"
)

// Seeder handles database seeding
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		db: db,
	}
}

// SeedDatabase seeds the database with initial data
func (s *Seeder) SeedDatabase() error {
	// Check if database is already seeded
	var count int64
	if err := s.db.Model(&model.Profile{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		fmt.Println("Database already seeded")
		return nil
	}

	fmt.Println("Seeding database...")

	// Create celebrities
	if err := s.seedCelebrities(); err != nil {
		return err
	}

	fmt.Println("Database seeded successfully")
	return nil
}

// seedCelebrities seeds the database with celebrity profiles
func (s *Seeder) seedCelebrities() error {
	celebrities := []struct {
		RealName      string
		Username      string
		Platform      string
		ProfileURL    string
		ImageURL      string
		Bio           string
		Verified      bool
		FollowerCount int64
		NameParts     []struct {
			NamePart string
			PartType string
		}
		Aliases      []string
		PlatformData map[string]string
	}{
		{
			RealName:      "Imane Anys",
			Username:      "pokimane",
			Platform:      "Twitch",
			ProfileURL:    "https://www.twitch.tv/pokimane",
			ImageURL:      "https://static-cdn.jtvnw.net/jtv_user_pictures/pokimane-profile_image-5ab2436567f7d9cd-300x300.jpeg",
			Bio:           "Variety streamer, content creator, and co-founder of OfflineTV",
			Verified:      true,
			FollowerCount: 9200000,
			NameParts: []struct {
				NamePart string
				PartType string
			}{
				{NamePart: "Imane", PartType: "first"},
				{NamePart: "Anys", PartType: "last"},
			},
			Aliases: []string{"poki", "pokimane"},
			PlatformData: map[string]string{
				"user_id":    "44445592",
				"created_at": "2012-11-17",
			},
		},
		{
			RealName:      "Shah Rukh Khan",
			Username:      "iamsrk",
			Platform:      "Twitter",
			ProfileURL:    "https://twitter.com/iamsrk",
			ImageURL:      "https://pbs.twimg.com/profile_images/1194461669341822976/8Qy9d7HY_400x400.jpg",
			Bio:           "Actor, Producer, and co-owner of Knight Riders Group",
			Verified:      true,
			FollowerCount: 42000000,
			NameParts: []struct {
				NamePart string
				PartType string
			}{
				{NamePart: "Shah", PartType: "first"},
				{NamePart: "Rukh", PartType: "middle"},
				{NamePart: "Khan", PartType: "last"},
			},
			Aliases: []string{"SRK", "King Khan", "King of Bollywood"},
			PlatformData: map[string]string{
				"user_id":    "101311381",
				"created_at": "2010-01-02",
			},
		},
		{
			RealName:      "Elon Musk",
			Username:      "elonmusk",
			Platform:      "Twitter",
			ProfileURL:    "https://twitter.com/elonmusk",
			ImageURL:      "https://pbs.twimg.com/profile_images/1683325380441128960/yRsRRjGO_400x400.jpg",
			Bio:           "CEO of Tesla, SpaceX, and X",
			Verified:      true,
			FollowerCount: 158000000,
			NameParts: []struct {
				NamePart string
				PartType string
			}{
				{NamePart: "Elon", PartType: "first"},
				{NamePart: "Musk", PartType: "last"},
			},
			Aliases: []string{"Technoking", "Dogefather"},
			PlatformData: map[string]string{
				"user_id":    "44196397",
				"created_at": "2009-06-02",
			},
		},
		{
			RealName:      "Linus Torvalds",
			Username:      "torvalds",
			Platform:      "GitHub",
			ProfileURL:    "https://github.com/torvalds",
			ImageURL:      "https://avatars.githubusercontent.com/u/1024025?v=4",
			Bio:           "Creator of Linux and Git",
			Verified:      false,
			FollowerCount: 190000,
			NameParts: []struct {
				NamePart string
				PartType string
			}{
				{NamePart: "Linus", PartType: "first"},
				{NamePart: "Torvalds", PartType: "last"},
			},
			Aliases: []string{"Linux Creator"},
			PlatformData: map[string]string{
				"user_id":    "1024025",
				"created_at": "2011-09-03",
			},
		},
	}

	ctx := context.Background()

	for _, celebrity := range celebrities {
		// Create profile
		profile := model.NewProfile(
			celebrity.RealName,
			celebrity.Username,
			celebrity.Platform,
			celebrity.ProfileURL,
			celebrity.ImageURL,
			celebrity.Bio,
			celebrity.Verified,
			celebrity.FollowerCount,
		)

		// Set timestamps
		profile.CreatedAt = time.Now()
		profile.UpdatedAt = time.Now()
		profile.LastUpdated = time.Now()

		// Save profile to get ID
		if err := s.db.WithContext(ctx).Create(profile).Error; err != nil {
			return err
		}

		// Add name parts
		for _, namePart := range celebrity.NameParts {
			np := model.NewNamePart(profile.ID, namePart.NamePart, namePart.PartType)
			if err := s.db.WithContext(ctx).Create(np).Error; err != nil {
				return err
			}
		}

		// Add aliases
		for _, alias := range celebrity.Aliases {
			a := model.NewAlias(profile.ID, alias)
			if err := s.db.WithContext(ctx).Create(a).Error; err != nil {
				return err
			}
		}

		// Add platform data
		for key, value := range celebrity.PlatformData {
			pd := model.NewPlatformData(profile.ID, key, value)
			if err := s.db.WithContext(ctx).Create(pd).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
