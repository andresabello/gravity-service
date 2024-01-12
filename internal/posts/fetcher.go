package posts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pi-search/pkg/tracer"
	"pi-search/shared/utilities"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type ItemLink struct {
	Href string `json:"href"`
}

type Links struct {
	WpItems []ItemLink `json:"wp:items"`
}

type PostType struct {
	Slug  string `json:"slug"`
	Links Links  `json:"_links"`
}

// FetchPosts stores posts into postgres database.
func Fetch(db *gorm.DB, baseUrl string) error {
	// Get all post types so that we can fetch them all
	endpoints, err := fetchEndpoints(baseUrl)
	if err != nil {
		return tracer.TraceError(err)
	}

	// Create a wait group to wait for all Goroutines to finish
	var wg sync.WaitGroup

	// Create a channel to collect errors from Goroutines
	errCh := make(chan error, len(endpoints))

	for _, url := range endpoints {
		wg.Add(1)
		// Start the process asynchronously.
		go fetchRequest(url, db, errCh, &wg)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		fmt.Println(tracer.TraceError(err))
	}

	return nil
}

func fetchEndpoints(baseURL string) ([]string, error) {
	//TODO: this should be an option that the user has. User should be able to dictate
	unnecessaryTypes := []string{
		"attachment",
		"nav_menu_item",
		"wp_block",
		"wp_template",
		"wp_template_part",
		"wp_navigation",
		"tribe_organizer",
		"tribe_venue",
		"rank_math_schema",
		"tribe_events",
	}
	response, err := http.Get(fmt.Sprintf("%s/wp-json/wp/v2/types", baseURL))
	if err != nil {
		return nil, tracer.TraceError(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, tracer.TraceError(
			fmt.Errorf("unexpected status code %d", response.StatusCode),
		)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, tracer.TraceError(err)
	}

	cleanJSON, err := cleanup(body)
	if err != nil {
		return nil, tracer.TraceError(err)
	}

	// Unmarshal the JSON data into the struct
	var data map[string]PostType
	err = json.Unmarshal(cleanJSON, &data)
	if err != nil {
		return nil, tracer.TraceError(err)
	}

	// Extract wp:items links into a slice of strings
	var wpItemsLinks []string
	for _, pt := range data {
		if utilities.ContainsString(unnecessaryTypes, pt.Slug) {
			continue
		}

		for _, link := range pt.Links.WpItems {

			// Parse the URL
			parsedURL, err := url.Parse(link.Href)
			if err != nil {
				return nil, tracer.TraceError(err)
			}

			// Extract the path
			path := parsedURL.Path
			wpItemsLinks = append(wpItemsLinks, baseURL+path)
		}
	}

	return wpItemsLinks, nil
}

func fetchRequest(
	baseURL string,
	db *gorm.DB,
	errCh chan<- error,
	wg *sync.WaitGroup,
) {
	// Ensure that the WaitGroup is decremented when the function exits
	defer wg.Done()

	perPage := 10
	for page := 1; ; page++ {
		response, err := http.Get(fmt.Sprintf(
			"%s?page=%d&per_page=%d",
			baseURL,
			page,
			perPage,
		))
		if err != nil {
			errCh <- tracer.TraceError(err)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			errCh <- tracer.TraceError(
				fmt.Errorf("unexpected status code %d", response.StatusCode),
			)
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			errCh <- tracer.TraceError(err)
		}

		cleanJSON, err := cleanup(body)
		if err != nil {
			errCh <- tracer.TraceError(err)
		}
		posts, err := postUnmarshaler(cleanJSON, baseURL)
		if err != nil {
			errCh <- tracer.TraceError(err)
		}

		for _, post := range posts {
			err = ingestUsingQueue(db, post)
			if err != nil {
				errCh <- tracer.TraceError(err)
			}
		}

		// Do I Need to check if this the last page???
		if len(posts) < perPage {
			break
		}
	}
}

func cleanup(body []byte) ([]byte, error) {
	// TODO: Only read JSON and not the strings from wordpress errors
	// Skip lines starting with "Deprecated" and Warning
	lines := strings.Split(string(body), "\n")
	var cleanJSON []byte
	for _, line := range lines {
		if !strings.HasPrefix(line, "Deprecated") && !strings.HasPrefix(line, "Warning") {
			cleanJSON = append(cleanJSON, []byte(line)...)
		}
	}

	return cleanJSON, nil
}

func postUnmarshaler(cleanJSON []byte, baseURL string) ([]*Post, error) {
	var posts []*Post
	err := json.Unmarshal(cleanJSON, &posts)
	if err != nil {
		return nil, tracer.TraceError(err)
	}

	for _, post := range posts {
		post.Source = baseURL
	}

	return posts, nil
}
