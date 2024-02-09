package crawl

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"pi-gravity/internal/models"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var crawledData string
var location string

func NewCrawlCommand(db *gorm.DB) *cobra.Command {
	myCommand := &cobra.Command{
		Use:   "crawlcommand",
		Short: "Crawls car years, makes and models from the web.",
		Run: func(cmd *cobra.Command, args []string) {
			crawl(cmd, args, db)
		},
	}

	return myCommand
}

func crawl(cmd *cobra.Command, args []string, db *gorm.DB) {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"),
		chromedp.WindowSize(500, 300),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	startYear := 2024
	endYear := 1942
	// Loop through the years in reverse order
	for year := startYear; year >= endYear; year-- {
		fmt.Println("CURRENT YEAR IS: ", year)
		// does year exist?
		yearDB, err := models.GetYear(db, year)
		if err != nil {
			yearDB = models.Year{
				Year: year,
			}
			err := models.CreateYear(db, &yearDB)
			if err != nil {
				panic(fmt.Errorf(
					"unable to create year %d; Error %s",
					yearDB.Year,
					err,
				))
			}
		}

		url := fmt.Sprintf("https://shop.advanceautoparts.com/capi/v33/vehicles/makes?type=3&year=%d", year)
		visitWebError := chromedp.Run(ctx,
			getData(url),
		)

		if visitWebError != nil {
			log.Println(visitWebError.Error())
		}

		var makesData []string
		err = json.Unmarshal([]byte(crawledData), &makesData)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, make := range makesData {
			fmt.Println("CURRENT MAKE IS: ", make)
			makeDB, err := models.GetMake(db, make)
			if err != nil {
				makeDB = models.Make{
					Name: make,
					Years: []models.Year{
						yearDB,
					},
				}
				err := models.CreateMake(db, &makeDB)
				if err != nil {
					panic(fmt.Errorf(
						"unable to create make.On Year %d, make is %s; Error: %s",
						yearDB.Year,
						make,
						err,
					))
				}
			} else {
				makeYear := models.MakeYear{MakeID: makeDB.ID, YearID: yearDB.ID}
				err := models.AssociateMakeYear(db, makeYear)
				if err != nil {
					panic(fmt.Errorf(
						"unable to create make_year. Record %d and %d error %s",
						makeYear.YearID,
						makeYear.MakeID,
						err,
					))
				}
			}

			url := fmt.Sprintf("https://shop.advanceautoparts.com/capi/v33/vehicles/models?type=3&year=%d&make=%s", year, make)
			visitWebError = chromedp.Run(ctx,
				getData(url),
			)

			if visitWebError != nil {
				log.Println(visitWebError.Error())
			}

			var modelsData []string
			err = json.Unmarshal([]byte(crawledData), &modelsData)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			for _, model := range modelsData {
				fmt.Println("CURRENT MODEL IS: ", model)
				// does model exist? - same name and year
				// no then insert model.
				modelDB, err := models.GetModel(db, model)
				if err != nil {
					modelDB = models.Model{
						Name:   model,
						MakeID: makeDB.ID,
						Years: []models.Year{
							yearDB,
						},
					}
					err := models.CreateModel(db, &modelDB)
					if err != nil {
						panic(fmt.Errorf(
							"unable to create model on year %d model name %s; Error %s",
							yearDB.Year,
							model,
							err,
						))
					}
				} else {
					modelYear := models.ModelYear{ModelID: modelDB.ID, YearID: yearDB.ID}
					err := models.AssociateModelYear(db, modelYear)
					if err != nil {
						panic(fmt.Errorf(
							"unable to create model_year. Year %d and Model %d error %s",
							modelYear.YearID,
							modelYear.ModelID,
							err,
						))
					}
				}
			}
		}
	}
}

func getData(url string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Location(&location),
		chromedp.Text(`pre`, &crawledData, chromedp.ByQueryAll),
	}
}
