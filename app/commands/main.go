package commands

import (
	"flag"
	"fmt"
	"gametime-hub/commands/fetchers/twitter"
	"os"
)

// Initializes command-line app and flags. Uses the Flag package.
func Initialize() {
	provider := flag.String(
		"provider",
		"",
		"Specify the provider to gather news from (e.g., 'Twitter', 'Google', 'FB')",
	)

	flag.Parse()

	// Check if the "action" flag is set
	switch *provider {
	case "google":
		fmt.Println("Get Items from Google")
	case "fb-thread":
		fmt.Println("Get Items from FB Threads")
	default:
		twitter.Fetch()
	}

	os.Exit(1)
}
