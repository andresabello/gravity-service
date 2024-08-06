// cmd/app/main.go
package main

import (
	"log"
	"os"
	"pi-gravity/internal/config"
	"pi-gravity/internal/app"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize your configuration
	config := &config.Config{} // Replace with actual initialization

	// Initialize the database connection
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create the root command
	rootCmd := app.NewAppCommand(config, db)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
