package main

import (
	"flag"

	"github.com/hdevillers/go-fannot/fannot"
)

func main() {
	query := flag.String("query", "", "Input query fasta file.")
	refdb := flag.String("refdb", "", "List of reference DB (coma separator).")
	dirdb := flag.String("dirdb", "", "Sub-directory that contains the reference DBs.")
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

	// Parse the list of reference DB
	fa.GetDBs(*refdb, *dirdb)

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
			} else if fa.DBs[fa.DBi].OverWrite && fa.Results[i].Status == 1 {
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

	// Printout the results
	fannot.PrintFAResultsHeader()
	for i := 0; i < fa.NQueries; i++ {
		fa.Results[i].PrintFAResult(fa.Queries[i].Id)
	}
}
