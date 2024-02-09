package models

import (
	"errors"
	"fmt"
	"pi-gravity/pkg/tracer"

	"gorm.io/gorm"
)

// Post represents a blog post from Wordpress
type Year struct {
	ID     int     `json:"id" gorm:"primaryKey"`
	Year   int     `json:"year" gorm:"year"`
	Makes  []Make  `json:"makes" gorm:"many2many:make_years;"`
	Models []Model `json:"models" gorm:"many2many:model_years;"`
}

// GetYear gets a year based on the lookup year.
func GetYear(db *gorm.DB, lookupYear int) (Year, error) {
	var year Year
	result := db.Where("year=?", lookupYear).First(&year)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return year, errors.New("record not found")
	}

	if result.Error != nil {
		return year, fmt.Errorf("error: %v", result.Error)
	}

	return year, nil
}

// CreateYear inserts a new year into the database.
func CreateYear(db *gorm.DB, year *Year) error {
	result := db.Omit("id").Create(year)

	if result.Error != nil {
		return tracer.TraceError(result.Error)
	}

	return nil
}
