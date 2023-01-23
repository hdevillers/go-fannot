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

// Initialize argument variables
var input string
var output string
var nsplit int
var compress bool

// Init function that manages input arguments
func init() {
	// Define default values and usages
	const (
		toolName        = "uniprot-split"
		toolDesc        = "Split an UniProt data file into multiple sub-files"
		inputDefault    = ""
		inputUsage      = "Input UniProt data file (dat(.gz))"
		outputDefault   = ""
		outputUsage     = "Output file basename"
		nsplitDefault   = 10
		nsplitUsage     = "Number of sub-data files wanted"
		compressDefault = false
		compressUsage   = "Compress output files"
	)

	// Init. flags
	flag.StringVar(&input, "input", inputDefault, inputUsage)
	flag.StringVar(&input, "i", inputDefault, inputUsage)
	flag.StringVar(&output, "output", outputDefault, outputUsage)
	flag.StringVar(&output, "o", outputDefault, outputUsage)
	flag.IntVar(&nsplit, "n-splits", nsplitDefault, nsplitUsage)
	flag.IntVar(&nsplit, "n", nsplitDefault, nsplitUsage)
	flag.BoolVar(&compress, "compress", compressDefault, compressUsage)
	flag.BoolVar(&compress, "c", compressDefault, compressUsage)

	// Shorthand associations
	shand := map[string]string{
		"input":    "i",
		"output":   "o",
		"n-split":  "n",
		"compress": "c",
	}

	// Usage print order
	order := []string{"input", "output", "n-splits", "compress"}

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

	if output == "" {
		panic("You must provide an output file basename.")
	}

	if nsplit < 2 {
		panic("The number of file division must be greater than 1.")
	}

	swr := uniprot.NewReader(input)
	swr.PanicOnError()
	defer swr.Close()

	writers := make([]*uniprot.Writer, nsplit)
	fileExt := ".dat"
	if compress {
		fileExt += ".gz"
	}

	for i := 0; i < nsplit; i++ {
		writers[i] = uniprot.NewWriter(output + fmt.Sprintf("%03d", i) + fileExt)
		writers[i].PanicOnError()
		defer writers[i].Close()
	}

	wi := 0
	for swr.Next() {
		writers[wi].WriteStrings(swr.GetData())
		writers[wi].WriteEntryEnd()
		writers[wi].PanicOnError()
		wi++
		if wi == nsplit {
			wi = 0
		}
	}

}
