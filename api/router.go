package api

import (
	"pi-search/api/controllers"
	"pi-search/internal/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Registers and starts the application's routes.
func StartRouter(r *gin.Engine, c *config.Config, db *gorm.DB) {
	preFlightRequestCache := 12
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        time.Duration(preFlightRequestCache) * time.Hour,
	}))

	r.GET("/status", controllers.GetHealth)
	r.GET("/ingest", controllers.IngestPosts(c, db))
	r.GET("/search", controllers.Search(c, db))
}
