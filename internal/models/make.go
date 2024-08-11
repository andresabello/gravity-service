package models

import (
	"errors"
	"fmt"
	"pi-gravity/internal/cache"
	"pi-gravity/pkg/tracer"
	"sort"
	"sync"

	"gorm.io/gorm"
)

// Make represents a car's make.
type Make struct {
	ID     int     `json:"id" gorm:"primaryKey"`
	Name   string  `json:"make_name" gorm:"column:make_name;not null;unique"`
	Models []Model `json:"models" gorm:"foreignKey:make_id"`
	Years  []Year  `json:"years" gorm:"many2many:make_years;"`
}

type Cars struct {
	Makes []Make `json:"makes" gorm:"foreignKey:id"`
}

type MakeResult struct {
	ID       int    `json:"id"`
	MakeName string `json:"make_name"`
}

type CarData struct {
	MakeID    int
	MakeName  string
	ModelID   int
	ModelName string
	YearID    int
	YearName  string
}

func GetAllMakes(db *gorm.DB, cacheInstance *cache.Cache) (interface{}, error) {
	var carData []CarData

	cacheKey := "make_model_year_data"
	// Check if the result is already cached
	if result, exists := cacheInstance.Get(cacheKey); exists {
		return result, nil
	}

	// Optimized query to fetch all necessary data
	db.Table("makes").
		Select("makes.id as make_id, makes.make_name as make_name, models.id as model_id, models.model_name as model_name, years.id as year_id, years.year as year_name").
		Joins("left join models on models.make_id = makes.id").
		Joins("left join make_years on make_years.make_id = makes.id").
		Joins("left join years on years.id = make_years.year_id").
		Order("makes.make_name, models.model_name, years.year").
		Scan(&carData)

	makeModelYearMap := make(map[string]*cache.MakeModelYear)

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Channel to process each record concurrently
	dataChan := make(chan CarData)

	// Start multiple workers to process the data concurrently
	numWorkers := 4 // You can adjust this based on your system's capabilities
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for data := range dataChan {
				mu.Lock()
				if _, exists := makeModelYearMap[data.MakeName]; !exists {
					makeModelYearMap[data.MakeName] = &cache.MakeModelYear{
						Make:   data.MakeName,
						Models: []string{},
						Years:  []string{},
					}
				}

				// Add model if it's not already in the list
				if !contains(makeModelYearMap[data.MakeName].Models, data.ModelName) {
					makeModelYearMap[data.MakeName].Models = append(makeModelYearMap[data.MakeName].Models, data.ModelName)
				}

				// Add year if it's not already in the list
				if !contains(makeModelYearMap[data.MakeName].Years, data.YearName) {
					makeModelYearMap[data.MakeName].Years = append(makeModelYearMap[data.MakeName].Years, data.YearName)
				}
				mu.Unlock()
			}
		}()
	}

	// Feed the data into the channel
	for _, data := range carData {
		dataChan <- data
	}

	// Close the channel and wait for all workers to finish
	close(dataChan)
	wg.Wait()

	// Convert the map to a slice
	makeModelYearList := []cache.MakeModelYear{}
	for _, value := range makeModelYearMap {
		// Sort models and years alphabetically
		sort.Strings(value.Models)
		sort.Strings(value.Years)
		makeModelYearList = append(makeModelYearList, *value)
	}

	// Sort the makes alphabetically
	sort.Slice(makeModelYearList, func(i, j int) bool {
		return makeModelYearList[i].Make < makeModelYearList[j].Make
	})

	// Cache the result
	cacheInstance.Set(cacheKey, makeModelYearList)

	return makeModelYearList, nil
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

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
