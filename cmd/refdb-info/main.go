package main

import (
	"flag"

	"github.com/hdevillers/go-fannot/custom"
	"github.com/hdevillers/go-fannot/refdb"
)

// Initialize argument variables
var id string
var dir string

// Init function that manages input arguments
func init() {
	// Define default values and usages
	const (
		toolName     = "refdb-info"
		toolDesc     = "Show information about a given reference database"
		refdbDefault = ""
		refdbUsage   = "Id of a reference database or path of its config.json file"
		dirdbDefault = ""
		dirdbUsage   = "Directory containing the reference databases"
	)

	// Init. flags
	flag.StringVar(&id, "refdb-id", refdbDefault, refdbUsage)
	flag.StringVar(&id, "r", refdbDefault, refdbUsage)
	flag.StringVar(&dir, "refdb-dir", dirdbDefault, dirdbUsage)
	flag.StringVar(&dir, "d", dirdbDefault, dirdbUsage)

	// Shorthand associations
	shand := map[string]string{
		"refdb-id":  "r",
		"refdb-dir": "d",
	}

	order := []string{"refdb-id", "refdb-dir"}

	// Custom usage display
	flag.Usage = func() {
		custom.Usage(*flag.CommandLine, toolName, toolDesc, &order, &shand)
	}

}

func main() {
	flag.Parse()

	// Check input values and find the JSON file
	if id == "" {
		panic("You must provide the ID of the queried reference database or its config.json file.")
	}

	// Load the refdb object
	rdb := refdb.FindRefDB(id, dir)

	// Print-out the info
	rdb.PrintInfoHeader()
	rdb.PrintInfo()
}
