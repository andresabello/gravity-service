package app

import (
	"log"
	"pi-gravity/api"
	"pi-gravity/internal/config"
	"pi-gravity/internal/crawl"

	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// App represents the application.
type App struct {
	Router *gin.Engine
}

// NewApp creates a new instance of the application.
func newApp(config *config.Config, db *gorm.DB) *App {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	api.StartRouter(router, config, db)

	return &App{
		Router: router,
	}
}


func NewAppCommand(config *config.Config, db *gorm.DB) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "app",
		Short: "Gravity Support Application",
		Long:  "Supports services provided by custom gravity forms wordpress plugins",
		Run: func(cmd *cobra.Command, args []string) {
			// Run your server logic here
			app := newApp(config, db)
			err := app.Router.Run(":8090")
			if err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		},
	}

	// Add your commands here
	rootCmd.AddCommand(crawl.NewCrawlCommand(config, db))

	return rootCmd
}