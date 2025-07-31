package container

import (
	"fmt"
	"log"

	appservice "github.com/accio/internal/application/service"
	"github.com/accio/internal/domain/repository"
	domainservice "github.com/accio/internal/domain/service"
	"github.com/accio/internal/image"
	"github.com/accio/internal/infrastructure/api"
	"github.com/accio/internal/infrastructure/persistence"
)

// Container is a dependency injection container
type Container struct {
	// Database
	Database *persistence.Database

	// Repositories
	ProfileRepository       repository.ProfileRepository
	SearchHistoryRepository repository.SearchHistoryRepository
	UserFeedbackRepository  repository.UserFeedbackRepository

	// Services
	ProfileService       domainservice.ProfileService
	SearchHistoryService appservice.SearchHistoryService

	// Platform clients
	PlatformClients map[string]api.PlatformClient

	// Image processor
	ImageProcessor *image.ImageProcessor

	// Seeder
	Seeder *persistence.Seeder
}

// NewContainer creates a new dependency injection container
func NewContainer() (*Container, error) {
	container := &Container{
		PlatformClients: make(map[string]api.PlatformClient),
	}

	// Initialize database
	db, err := persistence.NewDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	container.Database = db

	// Initialize repositories
	container.ProfileRepository = persistence.NewGormProfileRepository(db.DB)
	container.SearchHistoryRepository = persistence.NewGormSearchHistoryRepository(db.DB)
	container.UserFeedbackRepository = persistence.NewGormUserFeedbackRepository(db.DB)

	// Initialize services
	container.ProfileService = appservice.NewProfileService(container.ProfileRepository)
	container.SearchHistoryService = appservice.NewSearchHistoryService(container.SearchHistoryRepository)

	// Initialize platform clients
	if err := container.initializePlatformClients(); err != nil {
		log.Printf("Warning: Failed to initialize some platform clients: %v", err)
	}

	// Initialize image processor
	container.ImageProcessor = image.NewImageProcessor()

	// Initialize seeder
	container.Seeder = persistence.NewSeeder(db.DB)

	return container, nil
}

// initializePlatformClients initializes platform clients
func (c *Container) initializePlatformClients() error {
	// Initialize Twitter client
	twitterClient, err := api.NewTwitterClient()
	if err == nil {
		c.PlatformClients["Twitter"] = twitterClient
		profileService, ok := c.ProfileService.(*appservice.ProfileServiceImpl)
		if ok {
			profileService.RegisterPlatformClient(twitterClient)
		}
	}

	// Initialize other platform clients here...

	return nil
}

// Close closes the container and releases resources
func (c *Container) Close() error {
	if c.Database != nil {
		return c.Database.Close()
	}
	return nil
}

// SeedDatabase seeds the database with initial data
func (c *Container) SeedDatabase() error {
	if c.Seeder != nil {
		return c.Seeder.SeedDatabase()
	}
	return nil
}
