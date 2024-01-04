package controllers

import (
	"net/http"
	"pi-search/internal/config"
	"pi-search/internal/posts"
	"pi-search/pkg/tracer"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Search(config *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		searchString := c.Query("query")
		page := c.Query("page")
		if page == "" {
			page = "0"
		}

		allPosts, err := posts.GetPostsByQueryString(db, searchString, page)
		if err != nil {
			c.Error(tracer.TraceError(err))
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
			return
		}

		// Return with all the found posts from the search
		// We are searching by query string
		c.JSON(http.StatusOK, gin.H{
			"message": "Found Posts!",
			"posts":   allPosts,
		})
	}
}
