package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hdevillers/go-fannot/custom"
	"github.com/hdevillers/go-fannot/uniprot"
)

// Initialize argument variables
var input string
var output string
var ekeep string
var eskip string
var tkeep string
var tskip string
var dkeep string
var dskip string
var lmin int

// Init function that manages input arguments
func init() {
	// Define default values and usages
	const (
		toolName      = "uniprot-subset"
		toolDesc      = "Extract a subset of data from an UniProt data file"
		inputDefault  = ""
		inputUsage    = "Input UniProt data file (dat(.gz))"
		outputDefault = ""
		outputUsage   = "Output UniProt data file (dat(.gz))"
		ekeepDefault  = ""
		ekeepUsage    = "Evidence keep instruction (regex)"
		eskipDefault  = ""
		eskipUsage    = "Evidence skip instruction (regex)"
		tkeepDefault  = ""
		tkeepUsage    = "Taxonomy keep instruction (regex)"
		tskipDefault  = ""
		tskipUsage    = "Taxonomy skip instruction (regex)"
		dkeepDefault  = ""
		dkeepUsage    = "Description/function keep instruction (regex)"
		dskipDefault  = ""
		dskipUsage    = "Description/function skip instruction (regex)"
		lminDefault   = 30
		lminUsage     = "Minimal protein length (aa)"
	)

	// Init. flags
	flag.StringVar(&input, "input", inputDefault, inputUsage)
	flag.StringVar(&input, "i", inputDefault, inputUsage)
	flag.StringVar(&output, "output", outputDefault, outputUsage)
	flag.StringVar(&output, "o", outputDefault, outputUsage)
	flag.StringVar(&ekeep, "evidence-keep", ekeepDefault, ekeepUsage)
	flag.StringVar(&ekeep, "e", ekeepDefault, ekeepUsage)
	flag.StringVar(&eskip, "evidence-skip", eskipDefault, eskipUsage)
	flag.StringVar(&eskip, "E", eskipDefault, eskipUsage)
	flag.StringVar(&tkeep, "taxonomy-keep", tkeepDefault, tkeepUsage)
	flag.StringVar(&tkeep, "t", tkeepDefault, tkeepUsage)
	flag.StringVar(&tskip, "taxonomy-skip", tskipDefault, tskipUsage)
	flag.StringVar(&tskip, "T", tskipDefault, tskipUsage)
	flag.StringVar(&dkeep, "description-keep", dkeepDefault, dkeepUsage)
	flag.StringVar(&dkeep, "d", dkeepDefault, dkeepUsage)
	flag.StringVar(&dskip, "description-skip", dskipDefault, dskipUsage)
	flag.StringVar(&dskip, "D", dskipDefault, dskipUsage)
	flag.IntVar(&lmin, "min-length", lminDefault, lminUsage)
	flag.IntVar(&lmin, "l", lminDefault, lminUsage)

	// Shorthand associations
	shand := map[string]string{
		"input":            "i",
		"output":           "o",
		"evidence-keep":    "e",
		"evidence-skip":    "E",
		"taxonomy-keep":    "t",
		"taxonomy-skip":    "T",
		"description-keep": "d",
		"description-skip": "D",
		"min-length":       "l",
	}

	// Usage print order
	order := []string{"input", "output", "evidence-keep", "evidence-skip", "taxonomy-keep", "taxonomy-skip", "description-keep", "description-skip", "min-length"}

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
		panic("You must provide an output file name.")
	}
	// User must provide at least one keep/skip instruction
	if ekeep == "" && eskip == "" && tkeep == "" && tskip == "" && dkeep == "" && dskip == "" {
		panic("You must provide at least one keep/skip instruction.")
	}
	// Minimal length must be a positive integer
	if lmin < 0 {
		panic("The minimal protein length threshold must be positive.")
	}

	// Check if input is a single file or a base name for multiple files
	files := make([]string, 0)
	if _, err := os.Stat(input); errors.Is(err, os.ErrNotExist) {
		// This is probably not a single file, then look for multiple files
		files, err = filepath.Glob(input + "*")
		if err != nil {
			panic(err)
		}
		if files == nil {
			panic("Failed to found files from the provided pattern.")
		}
	} else {
		files = append(files, input)
	}

	// Initalize the channel
	entryChan := make(chan *[]string)
	threadChan := make(chan int)
	recordChan := make(chan int)

	// Initialze output writer
	sww := uniprot.NewSubsetWriter(output)
	sww.Writer.PanicOnError()
	defer sww.Writer.Close()

	// Launch the recording routine
	go sww.RecordEntry(entryChan, recordChan)

	// Init. a new subset object
	s := uniprot.Subset{
		Ekeep: ekeep, Eskip: eskip,
		Tkeep: tkeep, Tskip: tskip,
		Dkeep: dkeep, Dskip: dskip,
		Lmin: lmin}

	// Launch reading routine(s)
	for _, file := range files {
		// Parsing optimization if reading function is not necessary
		if dkeep == "" && dskip == "" {
			go s.LightParseFile(entryChan, threadChan, file)
		} else {
			go s.ParseFile(entryChan, threadChan, file)
		}
	}

	// Wait for reading threads
	tot := 0
	for i := 0; i < len(files); i++ {
		tot += <-threadChan
	}
	close(entryChan)

	// Wait for the recorder
	kpt := <-recordChan

	fmt.Println("Scan", tot, "entries and kept", kpt, "ones.")
}
