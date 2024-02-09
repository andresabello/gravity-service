package main

import (
	"log"
	"pi-gravity/cmd"
	"pi-gravity/internal/config"
	"pi-gravity/internal/database"
)

func main() {
	// Load configuration from /app/.env file
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	db, err := database.Start(*config)
	if err != nil {
		log.Fatalf("Error loading database: %v", err)
	}
	cmd.Execute(config, db)
}
