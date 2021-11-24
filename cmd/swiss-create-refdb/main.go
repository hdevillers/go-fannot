package main

import (
	"flag"

	"github.com/hdevillers/go-fannot/refdb"
)

func main() {
	input := flag.String("input", "", "Input swissProt data file.")
	name := flag.String("id", "", "Name of the reference database.")
	outdir := flag.String("outdir", ".", "Output directory.")
	equal := flag.Bool("equal", false, "Indicate that the reference contains genes from the query.")
	ow := flag.Bool("overwrite", false, "Indicate that annotations from this DB can overwrite annotation from other DB.")
	unre := flag.Bool("unreviewed", false, "Indicate if annotation are unreviewed (from TrEmbl).")
	gn := flag.Bool("gene-name", false, "Indicate if gene name can be transfered in query features.")
	desc := flag.String("desc", "No description", "Database description.")
	flag.Parse()

	if *input == "" {
		panic("You must provide a SwissProt data file.")
	}
	if *name == "" {
		panic("You must provide a name for the new reference database.")
	}

	// Create the refdb object
	rdb := refdb.NewRefdb(*outdir, *name, *input, *desc, *equal, *ow, !*unre, *gn)

	// Load the data
	rdb.LoadSource()

	// Save the json config
	rdb.WriteJson()
}
