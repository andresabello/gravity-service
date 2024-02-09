package models

import (
	"pi-gravity/pkg/tracer"

	"gorm.io/gorm"
)

// MakeYear represents the many-to-many relationship between makes and years.
type MakeYear struct {
	MakeID int
	YearID int
}

// AssociateMakeYear makes the association when year and make exist.
func AssociateMakeYear(db *gorm.DB, makeYear MakeYear) error {
	var existingAssociation MakeYear
	result := db.Where(
		"make_id = ? AND year_id = ?",
		makeYear.MakeID,
		makeYear.YearID,
	).First(&existingAssociation)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return tracer.TraceError(result.Error)
	}

	if result.RowsAffected > 0 {
		return nil
	}

	result = db.Create(&makeYear)
	if result.Error != nil {
		return tracer.TraceError(result.Error)
	}

	return nil
}
