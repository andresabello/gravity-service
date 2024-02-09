package app

import (
	"pi-gravity/api"
	"pi-gravity/internal/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// App represents the application.
type App struct {
	Router *gin.Engine
}

// NewApp creates a new instance of the application.
func NewApp(config *config.Config, db *gorm.DB) *App {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	api.StartRouter(router, config, db)

	return &App{
		Router: router,
	}
}
