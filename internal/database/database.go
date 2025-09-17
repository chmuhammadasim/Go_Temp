package database

import (
	"fmt"
	"log"

	"go-backend/internal/config"
	"go-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database wraps the GORM database connection
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.Config) (*Database, error) {
	var dialector gorm.Dialector
	
	switch cfg.Database.Type {
	case "postgres":
		dialector = postgres.Open(cfg.GetDSN())
	case "sqlite":
		dialector = sqlite.Open(cfg.GetDSN())
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}
	
	// Set log level based on environment
	logLevel := logger.Info
	if cfg.IsProduction() {
		logLevel = logger.Error
	}
	
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	return &Database{DB: db}, nil
}

// Migrate runs database migrations
func (d *Database) Migrate() error {
	log.Println("Running database migrations...")
	
	err := d.DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	log.Println("Database migrations completed successfully")
	return nil
}

// Seed creates initial data in the database
func (d *Database) Seed() error {
	// Check if admin user already exists
	var count int64
	d.DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)
	
	if count == 0 {
		// Create default admin user
		adminUser := &models.User{
			Email:     "admin@example.com",
			Username:  "admin",
			Password:  "admin123", // Will be hashed by BeforeCreate hook
			FirstName: "Admin",
			LastName:  "User",
			Role:      models.RoleAdmin,
			IsActive:  true,
		}
		
		if err := d.DB.Create(adminUser).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
		
		log.Println("Default admin user created: admin@example.com / admin123")
	}
	
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the GORM database instance
func (d *Database) GetDB() *gorm.DB {
	return d.DB
}