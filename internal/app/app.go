package app

import (
	"log"
	"pi-search/api"
	"pi-search/internal/config"
	"pi-search/internal/database"

	"github.com/gin-gonic/gin"
)

// App represents the application.
type App struct {
	Router *gin.Engine
}

// NewApp creates a new instance of the application.
func NewApp() *App {
	// Load configuration from /app/.env file
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	db, err := database.Start(*config)
	if err != nil {
		log.Fatalf("Error loading database: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	api.StartRouter(router, config, db)

	return &App{
		Router: router,
	}
}
