package crawl

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"pi-gravity/internal/config"
	"time"

	"pi-gravity/internal/models"
	"pi-gravity/pkg/tracer"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// Crawl fetches and processes car data from a web API
func Crawl(config *config.Config, db *gorm.DB) {
	startYear := config.CarStartYear
	endYear := config.CarEndYear

	for year := startYear; year >= endYear; year-- {
		fmt.Println("CURRENT YEAR IS:", year)
		yearDB, err := models.GetYear(db, year)
		if err != nil {
			yearDB = models.Year{Year: year}
			if err := models.CreateYear(db, &yearDB); err != nil {
				tracer.TraceError(fmt.Errorf("unable to create year %d; Error %s", yearDB.Year, err))
				continue
			}
		}

		endpointURL := fmt.Sprintf("%s/makes?type=3&year=%d", config.CarAPIURL, year)
		fmt.Println("Fetching URL:", endpointURL)

		responseBody, err := fetchWithRetries(endpointURL, 3)
		if err != nil {
			log.Printf("Failed to fetch data from URL %s after retries: %v", endpointURL, err)
			continue
		}

		var makesData []string
		if err := json.Unmarshal([]byte(responseBody), &makesData); err != nil {
			log.Printf("Error unmarshalling response body: %v", err)
			log.Printf("Response body was: %s", responseBody)
			continue
		}

		for _, make := range makesData {
			fmt.Println("CURRENT MAKE IS:", make)
			makeDB, err := models.GetMake(db, make)
			if err != nil {
				makeDB = models.Make{
					Name: make,
					Years: []models.Year{
						yearDB,
					},
				}
				if err := models.CreateMake(db, &makeDB); err != nil {
					tracer.TraceError(fmt.Errorf("unable to create make. On Year %d, make is %s; Error: %s", yearDB.Year, make, err))
					continue
				}
			}

			makeYear := models.MakeYear{MakeID: makeDB.ID, YearID: yearDB.ID}
			if err := models.AssociateMakeYear(db, makeYear); err != nil {
				tracer.TraceError(fmt.Errorf("unable to create make_year. Record %d and %d error %s", makeYear.YearID, makeYear.MakeID, err))
				continue
			}

			modelURL := fmt.Sprintf(
				"%s/models?type=3&year=%d&make=%s",
				config.CarAPIURL, year,
				url.QueryEscape(make),
			)
			fmt.Println("Fetching URL:", modelURL)

			responseBody, err = fetchWithRetries(modelURL, 3)
			if err != nil {
				log.Printf("Failed to fetch data from URL %s after retries: %v", modelURL, err)
				continue
			}

			var modelsData []string
			if err := json.Unmarshal([]byte(responseBody), &modelsData); err != nil {
				log.Printf("Error unmarshalling response body: %v", err)
				log.Printf("Response body was: %s", responseBody)
				continue
			}

			for _, model := range modelsData {
				fmt.Println("CURRENT MODEL IS:", model)
				modelDB, err := models.GetModel(db, model)
				if err != nil {
					modelDB = models.Model{
						Name:   model,
						MakeID: makeDB.ID,
						Years: []models.Year{
							yearDB,
						},
					}
					if err := models.CreateModel(db, &modelDB); err != nil {
						log.Printf("unable to create model on year %d model name %s; Error %s", yearDB.Year, model, err)
						continue
					}
				}

				modelYear := models.ModelYear{ModelID: modelDB.ID, YearID: yearDB.ID}
				if err := models.AssociateModelYear(db, modelYear); err != nil {
					log.Printf("unable to create model_year. Year %d and Model %d error %s", modelYear.YearID, modelYear.ModelID, err)
				}
			}
		}
	}
}

// fetchWithRetries fetches data from a URL with retry logic
func fetchWithRetries(url string, maxRetries int) (string, error) {
	var responseBody string
	var err error
	for retryCount := 0; retryCount < maxRetries; retryCount++ {
		responseBody, err = getData(url)
		if err == nil {
			return responseBody, nil
		}
		log.Printf("Retry %d/%d for URL %s: %v", retryCount+1, maxRetries, url, err)
		time.Sleep(2 * time.Second) // optional: sleep between retries
	}
	return "", fmt.Errorf("failed to fetch data from URL %s after %d retries: %v", url, maxRetries, err)
}

// getData fetches data from a URL and returns the response body
func getData(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create GET request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "PostmanRuntime/7.37.3")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Check if the response body is HTML
	if isHTMLResponse(body) {
		return "", fmt.Errorf("received HTML response instead of JSON: %s", string(body))
	}

	return string(body), nil
}

// Helper function to check if response is HTML
func isHTMLResponse(body []byte) bool {
	// HTML responses usually start with "<!DOCTYPE html>" or "<html>"
	return len(body) > 0 && (body[0] == '<')
}

// NewCrawlCommand creates a new crawl command
func NewCrawlCommand(config *config.Config, db *gorm.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "crawl",
		Short: "Crawl data from the web",
		Run: func(cmd *cobra.Command, args []string) {
			Crawl(config, db)
		},
	}
}
