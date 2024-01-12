package database

import (
	"fmt"
	"pi-search/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	formPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Starts the Database with all the propper configuration.
func Start(config config.Config) (*gorm.DB, error) {
	sql, err := connectDB(config)
	if err != nil {
		return nil, err
	}

	err = runMigration(sql)
	if err != nil {
		return nil, err
	}

	return sql, nil
}

// Open a database connection
func connectDB(config config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=America/New_York",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		5432,
	)

	db, err := gorm.Open(formPostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// PerformMigrations performs database migrations if needed.
func runMigration(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Use the database connection to create a driver
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	// Specify the migration source (migrations directory)
	sourceURL := "file:///app/internal/database/migrations/"

	// Create a new migration instance
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return err
	}

	// Check if any migrations are pending
	err = m.Up()
	if err == migrate.ErrNoChange {
		fmt.Println("No migrations to apply.")
		return nil
	}

	return err
}
