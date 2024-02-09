package models

import (
	"pi-gravity/pkg/tracer"

	"gorm.io/gorm"
)

// ModelYear represents the many-to-many relationship between models and years.
type ModelYear struct {
	ModelID int
	YearID  int
}

// AssociateModelYear makes the association when year and model exist.
func AssociateModelYear(db *gorm.DB, modelYear ModelYear) error {
	var existingAssociation ModelYear
	result := db.Where(
		"model_id = ? AND year_id = ?",
		modelYear.ModelID,
		modelYear.YearID,
	).First(&existingAssociation)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return tracer.TraceError(result.Error)
	}

	if result.RowsAffected > 0 {
		return nil
	}

	result = db.Create(&modelYear)
	if result.Error != nil {
		return tracer.TraceError(result.Error)
	}

	return nil
}
