package main

import (
	"flag"
	"fmt"

	"github.com/hdevillers/go-fannot/fannot"
)

func main() {
	query := flag.String("query", "", "Input query fasta file.")
	refdb := flag.String("refdb", "", "List of reference DB (coma separator).")
	dirdb := flag.String("dirdb", "", "Sub-directory that contains the reference DBs.")
	rules := flag.String("rules", "", "JSON file containing similarity levels.")
	ipsin := flag.String("ips", "", "InterProScan output predictions (TSV format).")
	threads := flag.Int("threads", 4, "Number of threads.")
	flag.Parse()

	if *query == "" {
		panic("You must provide an input query file.")
	}
	if *refdb == "" {
		panic("You must provide at least one reference DB.")
	}

	// Initialize the functional annotation strucutre
	fa := fannot.NewFannot(*query)

	// Reset rules if a JSON is provided
	if *rules != "" {
		fa.Param = *fannot.NewParamFromJson(*rules)
	}

	// Parse the list of reference DB
	fa.GetDBs(*refdb, *dirdb)

	// Load ips if provided
	if *ipsin != "" {
		err := fa.Ips.LoadIpsData(*ipsin)
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
		for i := 0; i < *threads; i++ {
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
		for i := 0; i < *threads; i++ {
			<-threadChan
		}

		// If every sequence has a function, then stop
		if nq == 0 {
			break REFDB
		}
	}

	// Complete with IPS annotation if provided
	if *ipsin != "" {
		fa.AddIpsAnnot()
	}

	// Printout the results
	fmt.Print(fannot.Header())
	//fannot.PrintFAResultsHeader()
	for i := 0; i < fa.NQueries; i++ {
		fmt.Print(fa.Results[i].ToString(fa.Queries[i].Id))
	}
}
