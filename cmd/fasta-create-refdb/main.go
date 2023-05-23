package main

import (
	"flag"

	"github.com/hdevillers/go-fannot/custom"
	"github.com/hdevillers/go-fannot/refdb"
)

// Initialize argument variables
var input string
var name string
var outdir string
var equal bool
var ow bool
var unre bool
var gn bool
var desc string

// Init function that manages input arguments
func init() {
	// Define default values and usages
	const (
		toolName     = "fasta-create-refdb"
		toolDesc     = "Create a reference database from a Fasta file"
		inputDefault = ""
		inputUsage   = "Input Fasta file"
		refdbDefault = ""
		refdbUsage   = "Id of the reference database to create"
		dirdbDefault = ""
		dirdbUsage   = "Directory containing the reference databases"
		equalDefault = false
		equalUsage   = "Indicate that the reference database contains genes from the query"
		owDefault    = false
		owUsage      = "Indicate that annotations from this database can overwrite annotation from other databases"
		unreDefault  = false
		unreUsage    = "Indicate that annotations are unreviewed (e.g., from TrEmbl)"
		descDefault  = "No description"
		descUsage    = "Short description of the reference database"
	)

	// Init. flags
	flag.StringVar(&input, "input", inputDefault, inputUsage)
	flag.StringVar(&input, "i", inputDefault, inputUsage)
	flag.StringVar(&name, "refdb-id", refdbDefault, refdbUsage)
	flag.StringVar(&name, "r", refdbDefault, refdbUsage)
	flag.StringVar(&outdir, "refdb-dir", dirdbDefault, dirdbUsage)
	flag.StringVar(&outdir, "d", dirdbDefault, dirdbUsage)
	flag.BoolVar(&equal, "equal", equalDefault, equalUsage)
	flag.BoolVar(&equal, "e", equalDefault, equalUsage)
	flag.BoolVar(&ow, "overwrite", owDefault, owUsage)
	flag.BoolVar(&ow, "w", owDefault, owUsage)
	flag.BoolVar(&unre, "unreviewed", unreDefault, unreUsage)
	flag.BoolVar(&unre, "u", unreDefault, unreUsage)
	flag.StringVar(&desc, "description", descDefault, descUsage)
	flag.StringVar(&desc, "D", descDefault, descUsage)

	// Shorthand associations
	shand := map[string]string{
		"input":       "i",
		"refdb-id":    "r",
		"refdb-dir":   "d",
		"equal":       "e",
		"overwrite":   "w",
		"unreviewed":  "u",
		"description": "D",
	}

	// Usage print order
	order := []string{"input", "refdb-id", "refdb-dir", "equal", "overwrite", "unreviewed", "description"}

	// Custom usage display
	flag.Usage = func() {
		custom.Usage(*flag.CommandLine, toolName, toolDesc, &order, &shand)
	}
}

func main() {
	flag.Parse()

	if input == "" {
		panic("You must provide a UniProt data file.")
	}
	if name == "" {
		panic("You must provide a name for the new reference database.")
	}

	// Create the refdb object
	//rdb := refdb.NewRefdb(outdir, name, input, desc, equal, ow, !unre, gn)
	rdb := refdb.NewRefdb(outdir, name, input, desc, equal, ow, !unre)
	// Load the data
	rdb.LoadFasta()

	// Save the json config
	rdb.WriteJson()
}
