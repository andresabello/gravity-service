package main

import (
	"gametime-hub/commands"
	"gametime-hub/controllers"
	"gametime-hub/database"
)

func main() {
	database.Start()
	controllers.Start()
}

func init() {
	commands.Initialize()
}
