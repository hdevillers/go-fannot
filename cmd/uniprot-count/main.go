package main

import (
	"flag"
	"fmt"

	"github.com/hdevillers/go-fannot/custom"
	"github.com/hdevillers/go-fannot/uniprot"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Initialize argument variable
var input string

// Init function that manages input arguments
func init() {
	// Define default values and usages
	const (
		toolName     = "uniprot-count"
		toolDesc     = "Count the number of entries in an UniProt data file"
		inputDefault = ""
		inputUsage   = "Input UniProt data file (dat(.gz))"
	)

	// Init. flags
	flag.StringVar(&input, "input", inputDefault, inputUsage)
	flag.StringVar(&input, "i", inputDefault, inputUsage)

	// Shorthand associations
	shand := map[string]string{
		"input": "i",
	}

	// Usage print order
	order := []string{"input"}

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

	// Create a reader
	swr := uniprot.NewReader(input)
	defer swr.Close()

	// Count entry
	cnt := 0
	for swr.Next() {
		cnt++
	}

	// Display the number of entry
	if cnt == 0 {
		panic("No entry found, please check the input file format.")
	} else if cnt == 1 {
		fmt.Println("Found 1 UniProt entry in", input)
	} else {
		fmt.Println("Found", cnt, "UniProt entries in", input)
	}
}
