package models

import (
	"errors"
	"fmt"
	"pi-gravity/pkg/tracer"

	"gorm.io/gorm"
)

// Model represents a car's model in the database.
type Model struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	Name   string `json:"model_name" gorm:"column:model_name;not null;unique"`
	MakeID int    `json:"make_id"`
	Years  []Year `json:"years" gorm:"many2many:model_years;"`
}

// GetMake gets a make based on the make's name.
func GetModel(db *gorm.DB, modelName string) (Model, error) {
	var model Model
	result := db.Where(
		"model_name = ?",
		modelName,
	).First(&model)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model, errors.New("record not found")
	}

	if result.Error != nil {
		return model, fmt.Errorf("error: %v", result.Error)
	}

	return model, nil
}

// CreatePost inserts a new post into the database.
func CreateModel(db *gorm.DB, model *Model) error {
	result := db.Omit("id").Create(model)

	if result.Error != nil {
		return tracer.TraceError(result.Error)
	}

	return nil
}

// GetModelsByMakeAndYear gets car models by make and year.
func GetModelsByMakeAndYear(db *gorm.DB, make_id int, year int) ([]string, error) {
	var models []string
	result := db.Table("models").
		Joins("JOIN model_years ON models.id = model_years.model_id").
		Joins("JOIN years ON model_years.year_id = years.id").
		Where("years.year = ?", year).
		Where("models.make_id = ?", make_id).
		Pluck("models.model_name", &models)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models, fmt.Errorf(
			"no records found for makes made in %d",
			year,
		)
	}

	if result.Error != nil {
		return models, fmt.Errorf("error: %v", result.Error)
	}

	return models, nil
}
