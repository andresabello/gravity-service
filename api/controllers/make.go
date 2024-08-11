package controllers

import (
	"net/http"
	"pi-gravity/internal/cache"
	"pi-gravity/internal/config"
	"pi-gravity/internal/models"
	"pi-gravity/pkg/tracer"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllMakes(config *config.Config, db *gorm.DB, cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		allMakes, err := models.GetAllMakes(db, cache)
		if err != nil {
			c.Error(tracer.TraceError(err))
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Found Posts!",
			"makes":   allMakes,
		})
	}
}

func GetMakes(config *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		year := c.Query("year")
		num, err := strconv.Atoi(year)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
			return
		}

		allMakes, err := models.GetMakesByYear(db, num)
		if err != nil {
			c.Error(tracer.TraceError(err))
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Found Posts!",
			"makes":   allMakes,
		})
	}
}
