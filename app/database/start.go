package database

import (
	"fmt"
	"gametime-hub/config"
	"log"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Starts the Database with all the propper configuration.
func Start() {
	// Load configuration from .env file
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Construct the database URL
	err = runMigration(config.ConstructDatabaseURL())
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
}

// PerformMigrations performs database migrations if needed.
func runMigration(databaseURL string) error {
	// Open a database connection
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Use the database connection to create a driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Specify the migration source (migrations directory)
	sourceURL := "file:///app/database/migrations/"

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
