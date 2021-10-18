package main

import (
	"flag"

	"github.com/hdevillers/go-fannot/refdb"
)

func main() {
	id := flag.String("id", "", "Id of the reference database or path of the config.json file.")
	dir := flag.String("dir", ".", "Directory that contain databases.")
	flag.Parse()

	// Check input values and find the JSON file
	if *id == "" {
		panic("You must provide the ID of the queried reference database or its config.json file.")
	}

	// Load the refdb object
	rdb := refdb.FindRefDB(*id, *dir)

	// Print-out the info
	rdb.PrintInfoHeader()
	rdb.PrintInfo()
}
