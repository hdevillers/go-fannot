package main

import (
	"flag"

	"github.com/hdevillers/go-fannot/refdb"
)

func main() {
	input := flag.String("input", "", "Input swissProt data file.")
	name := flag.String("id", "", "Name of the reference database.")
	outdir := flag.String("outdir", ".", "Output directory.")
	desc := flag.String("desc", "No description", "Database description.")
	flag.Parse()

	if *input == "" {
		panic("You must provide a SwissProt data file.")
	}
	if *name == "" {
		panic("You must provide a name for the new reference database.")
	}

	// Create the refdb object
	rdb := refdb.NewRefdb(*outdir, *name, *input, *desc)

	// Load the data
	rdb.LoadSource()

	// Save the json config
	rdb.WriteJson()
}
