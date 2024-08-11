package api

import (
	"pi-gravity/api/controllers"
	"pi-gravity/internal/cache"
	"pi-gravity/internal/config"
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

	cache := cache.NewCache()

	// Create a group for API version 1
	apiV1 := r.Group("/api/v1")
	apiV1.GET("/status", controllers.GetHealth)
	apiV1.GET("/makes/", controllers.GetMakes(c, db))
	apiV1.GET("/models", controllers.GetModels(c, db))
	apiV1.GET("/all-makes", controllers.GetAllMakes(c, db, cache))
	apiV1.GET("/all-models", controllers.GetAllModels(c, db, cache))
	apiV1.POST("/images/upload", controllers.GetModels(c, db))
}
