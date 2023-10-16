package main

import (
	"top-reiki-healing/commands"
	"top-reiki-healing/controllers"
	"top-reiki-healing/database"
)

func main() {
	database.Start()
	controllers.Start()
}

func init() {
	commands.Initialize()
}
