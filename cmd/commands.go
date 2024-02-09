package cmd

import (
	"fmt"
	"pi-gravity/cmd/app"
	"pi-gravity/internal/config"

	"gorm.io/gorm"
)

// Execute starts the command execution
func Execute(config *config.Config, db *gorm.DB) {
	cmd := app.NewAppCommand(config, db)
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
