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

// Model represents a car's model in the database.
type Model struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	Name   string `json:"model_name" gorm:"column:model_name;not null;unique"`
	MakeID int    `json:"make_id"`
	Years  []Year `json:"years" gorm:"many2many:model_years;"`
}

func GetAllModels(db *gorm.DB, cacheInstance *cache.Cache) (interface{}, error) {
	var carData []CarData

	cacheKey := "model_make_year_data"
	// Check if the result is already cached
	if result, exists := cacheInstance.Get(cacheKey); exists {
		return result, nil
	}

	// Optimized query to fetch all necessary data
	db.Table("models").
		Select("models.id as model_id, models.model_name as model_name, makes.id as make_id, makes.make_name as make_name, years.id as year_id, years.year as year_name").
		Joins("left join makes on makes.id = models.make_id").
		Joins("left join model_years on model_years.model_id = models.id").
		Joins("left join years on years.id = model_years.year_id").
		Order("model_name, make_name, year_name").
		Scan(&carData)

	modelMakeYearMap := make(map[string]*cache.ModelMakeYear)

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
				if _, exists := modelMakeYearMap[data.ModelName]; !exists {
					modelMakeYearMap[data.ModelName] = &cache.ModelMakeYear{
						Model: data.ModelName,
						Makes: []string{},
						Years: []string{},
					}
				}

				// Add make if it's not already in the list
				if !contains(modelMakeYearMap[data.ModelName].Makes, data.MakeName) {
					modelMakeYearMap[data.ModelName].Makes = append(modelMakeYearMap[data.ModelName].Makes, data.MakeName)
				}

				// Add year if it's not already in the list
				if !contains(modelMakeYearMap[data.ModelName].Years, data.YearName) {
					modelMakeYearMap[data.ModelName].Years = append(modelMakeYearMap[data.ModelName].Years, data.YearName)
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
	modelMakeYearList := []cache.ModelMakeYear{}
	for _, value := range modelMakeYearMap {
		// Sort models and years alphabetically
		sort.Strings(value.Makes)
		sort.Strings(value.Years)
		modelMakeYearList = append(modelMakeYearList, *value)
	}

	// Sort the models alphabetically
	sort.Slice(modelMakeYearList, func(i, j int) bool {
		return modelMakeYearList[i].Model < modelMakeYearList[j].Model
	})

	// Cache the result
	cacheInstance.Set(cacheKey, modelMakeYearList)

	return modelMakeYearList, nil
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
