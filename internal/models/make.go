package models

import (
	"errors"
	"fmt"
	"pi-gravity/pkg/tracer"

	"gorm.io/gorm"
)

// Make represents a car's make.
type Make struct {
	ID     int     `json:"id" gorm:"primaryKey"`
	Name   string  `json:"make_name" gorm:"column:make_name;not null;unique"`
	Models []Model `json:"models" gorm:"foreignKey:make_id"`
	Years  []Year  `json:"years" gorm:"many2many:make_years;"`
}

type MakeResult struct {
	ID       int    `json:"id"`
	MakeName string `json:"make_name"`
}

// GetMake gets a make based on the make's name.
func GetMake(db *gorm.DB, makeName string) (Make, error) {
	var make Make
	result := db.Where("make_name=?", makeName).First(&make)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return make, errors.New("record not found")
	}

	if result.Error != nil {
		return make, fmt.Errorf("error: %v", result.Error)
	}

	return make, nil
}

// CreateMake inserts a new make into the database.
func CreateMake(db *gorm.DB, make *Make) error {
	result := db.Omit("id").Create(make)
	if result.Error != nil {
		return tracer.TraceError(result.Error)
	}

	return nil
}

// GetMakesByYear gets car makes by year.
func GetMakesByYear(db *gorm.DB, year int) ([]MakeResult, error) {
	var makeResults []MakeResult
	result := db.Table("makes").
		Joins("JOIN make_years ON makes.id = make_years.make_id").
		Joins("JOIN years ON make_years.year_id = years.id").
		Where("years.year = ?", year).
		Select("makes.id, makes.make_name").
		Find(&makeResults)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return makeResults, fmt.Errorf(
			"no records found for makes made in %d",
			year,
		)
	}

	if result.Error != nil {
		return makeResults, fmt.Errorf("error: %v", result.Error)
	}

	return makeResults, nil
}
