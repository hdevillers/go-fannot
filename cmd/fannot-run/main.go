package main

import (
	"flag"

	"github.com/hdevillers/go-fannot/custom"
	"github.com/hdevillers/go-fannot/fannot"
)

// Initialize argument variables
var input string
var output string
var refdb string
var dirdb string
var rules string
var ipsin string
var threads int

// Init function that manages input arguments
func init() {
	// Define default values and usages
	const (
		toolName       = "fannot-run"
		toolDesc       = "Run the functional annotation pipeline"
		inputDefault   = ""
		inputUsage     = "Input query file (fasta)"
		outputDefault  = ""
		outputUsage    = "Output file path (tsv)"
		refdbDefault   = ""
		refdbUsage     = "List of reference database ids (coma separator)"
		dirdbDefault   = ""
		dirdbUsage     = "Directory containing the reference databases"
		rulesDefault   = ""
		rulesUsage     = "JSON file containing similarity levels"
		ipsDefault     = ""
		ipsUsage       = "InterProScan output predictions (TSV format)"
		threadsDefault = 4
		threadsUsage   = "Number of threads"
	)

	// Init. flags
	flag.StringVar(&input, "input", inputDefault, inputUsage)
	flag.StringVar(&input, "i", inputDefault, inputUsage)
	flag.StringVar(&output, "output", outputDefault, outputUsage)
	flag.StringVar(&output, "o", outputDefault, outputUsage)
	flag.StringVar(&refdb, "refdb-id", refdbDefault, refdbUsage)
	flag.StringVar(&refdb, "r", refdbDefault, refdbUsage)
	flag.StringVar(&dirdb, "refdb-dir", dirdbDefault, dirdbUsage)
	flag.StringVar(&dirdb, "d", dirdbDefault, dirdbUsage)
	flag.StringVar(&rules, "rules", rulesDefault, rulesUsage)
	flag.StringVar(&ipsin, "ips", ipsDefault, ipsUsage)
	flag.IntVar(&threads, "threads", threadsDefault, threadsUsage)
	flag.IntVar(&threads, "t", threadsDefault, threadsUsage)

	// Shorthand associations
	shand := map[string]string{
		"input":     "i",
		"output":    "o",
		"refdb-id":  "r",
		"refdb-dir": "d",
		"threads":   "t",
	}

	// Usage print order
	order := []string{"input", "output", "refdb-id", "refdb-dir", "rules", "ips", "threads"}

	// Custom usage display
	flag.Usage = func() {
		custom.Usage(*flag.CommandLine, toolName, toolDesc, &order, &shand)
	}
}

func main() {
	flag.Parse()

	if input == "" {
		panic("You must provide an input query file.")
	}
	if refdb == "" {
		panic("You must provide at least one reference DB.")
	}
	// The number of threads must be greatter than 0
	if threads <= 0 {
		panic("At least one thread is required.")
	}

	// Initialize the functional annotation strucutre
	fa := fannot.NewFannot(input)

	// Reset rules if a JSON is provided
	if rules != "" {
		fa.Param = *fannot.NewParamFromJson(rules)
	}

	// Parse the list of reference DB
	fa.GetDBs(refdb, dirdb)

	// Load ips if provided
	if ipsin != "" {
		err := fa.Ips.LoadIpsData(ipsin)
		if err != nil {
			panic(err)
		}
	}

	// Init. annotation templates/format
	fa.NoteFormat = *fannot.NewFormat(fa.Param.TemplateNote)
	fa.ProductFormat = *fannot.NewFormat(fa.Param.TemplateProduct)
	fa.GeneNameFormat = *fannot.NewFormat(fa.Param.TemplateGeneName)
	fa.FunctionFormat = *fannot.NewFormat(fa.Param.TemplateFunction)

REFDB:
	for fa.NextDB() {
		// Create the channels for multithreading
		queryChan := make(chan int)
		threadChan := make(chan int)

		// Launch parallel go routines
		for i := 0; i < threads; i++ {
			go fa.FindFunction(queryChan, threadChan)
		}

		// throw gene index that require a function
		nq := 0 // Number of thrown queries
		for i := 0; i < fa.NQueries; i++ {
			if !fa.Finished[i] {
				nq++
				queryChan <- i
			} else if fa.DBs[fa.DBi].OverWrite && fa.Results[i].Status <= fa.Param.MaxStatusOW {
				// Try to overwrite the annotation
				nq++
				queryChan <- i
			}
		}
		close(queryChan)

		// Wait for all threads
		for i := 0; i < threads; i++ {
			<-threadChan
		}

		// If every sequence has a function, then stop
		if nq == 0 {
			break REFDB
		}
	}

	// Complete with IPS annotation if provided
	if ipsin != "" {
		fa.AddIpsAnnot()
	}

	// Printout the results
	fa.WriteOut(output)
}
