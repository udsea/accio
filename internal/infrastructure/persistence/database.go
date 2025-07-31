package persistence

import (
	"fmt"
	"os"

	"github.com/accio/internal/domain/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database represents the database connection
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase() (*Database, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found, using environment variables\n")
	}

	dbURL := os.Getenv("TURSO_DATABASE_URL")
	if dbURL == "" {
		// Use SQLite as fallback
		dbURL = "file::memory:?cache=shared"
	}

	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	// Connect to database
	db, err := gorm.Open(sqlite.Open(dbURL), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate schema
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Database{DB: db}, nil
}

// autoMigrate automatically migrates the schema
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Profile{},
		&model.NamePart{},
		&model.Alias{},
		&model.PlatformData{},
		&model.SearchHistory{},
		&model.UserFeedback{},
	)
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
