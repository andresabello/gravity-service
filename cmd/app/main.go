package main

import (
	"log"
	"pi-search/internal/app"
)

func main() {
	app := app.NewApp()
	err := app.Router.Run(":8090")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
