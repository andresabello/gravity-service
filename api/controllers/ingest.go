package controllers

import (
	"net/http"
	"pi-search/internal/config"
	"pi-search/internal/posts"
	"pi-search/pkg/tracer"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IngestPosts(config *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		websiteUrl := c.Query("website")
		appEnv := config.APPEnv

		// TODO: I order to test multiple websites sending in request, 
		// you need to identify the host and then replace that for host.docker.internal
		// Everything else stays the same
		
		if appEnv == "development" {
			websiteUrl = "http://host.docker.internal:7001"
		}

		err := posts.Fetch(db, websiteUrl)
		if err != nil {
			c.Error(tracer.TraceError(err))
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
			return
		}

		// TODO: You will have options like post type, and other filters.
		// TODO: IN WP make an endpoint to allow for status bar on ingestion.
		c.JSON(http.StatusOK, gin.H{
			"message": "The data ingestion process started!",
		})
	}
}
