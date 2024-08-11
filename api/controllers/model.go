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

func GetAllModels(config *config.Config, db *gorm.DB, cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		allModels, err := models.GetAllModels(db, cache)
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
			"models":  allModels,
		})
	}
}

func GetModels(config *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		year := c.Query("year")
		yearINT, err := strconv.Atoi(year)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
			return
		}

		makeID := c.Query("make_id")
		makeINT, err := strconv.Atoi(makeID)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				gin.H{"error": tracer.TraceError(err).Error()},
			)
			return
		}

		allModels, err := models.GetModelsByMakeAndYear(db, makeINT, yearINT)
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
			"models":  allModels,
		})
	}
}
