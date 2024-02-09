package app

import (
	"log"
	"pi-gravity/cmd/crawl"
	"pi-gravity/internal/app"
	"pi-gravity/internal/config"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func NewAppCommand(config *config.Config, db *gorm.DB) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "app",
		Short: "Gravity Support Application",
		Long:  "Supports services provided by custom gravity forms wordpress plugins",
		Run: func(cmd *cobra.Command, args []string) {
			// Run your server logic here
			app := app.NewApp(config, db)
			err := app.Router.Run(":8090")
			if err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		},
	}

	// Add your commands here
	rootCmd.AddCommand(crawl.NewCrawlCommand(db))

	return rootCmd
}
