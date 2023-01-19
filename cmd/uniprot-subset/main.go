package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hdevillers/go-fannot/uniprot"
)

func main() {
	input := flag.String("i", "", "Input UniProt data file.")
	output := flag.String("o", "", "Output data file.")
	ekeep := flag.String("e", "", "Evident keep instruction (regex).")
	eskip := flag.String("E", "", "Evidence skip instruction (regex).")
	tkeep := flag.String("t", "", "Taxonomy keep instruction (regex).")
	tskip := flag.String("T", "", "Taxonomy skip instruction (regex).")
	dkeep := flag.String("d", "", "Description/function keep instruction (regex).")
	dskip := flag.String("D", "", "Description/function skip instruction (regex).")
	lmin := flag.Int("l", 30, "Minimal protein length (aa).")
	flag.Parse()

	if *input == "" {
		panic("You must provide a UniProt data file.")
	}

	if *output == "" {
		panic("You must provide an output file name.")
	}

	if *ekeep == "" && *eskip == "" && *tkeep == "" && *tskip == "" && *dkeep == "" && *dskip == "" {
		panic("You must provide at least one keep/skip instruction.")
	}

	// Check if input is a single file or a base name for multiple files
	files := make([]string, 0)
	if _, err := os.Stat(*input); errors.Is(err, os.ErrNotExist) {
		// This is probably not a single file, then look for multiple files
		files, err = filepath.Glob(*input + "*")
		if err != nil {
			panic(err)
		}
		if files == nil {
			panic("Failed to found files from the provided pattern.")
		}
	} else {
		files = append(files, *input)
	}

	// Initalize the channel
	entryChan := make(chan *[]string)
	threadChan := make(chan int)
	recordChan := make(chan int)

	// Initialze output writer
	sww := uniprot.NewSubsetWriter(*output)
	sww.Writer.PanicOnError()
	defer sww.Writer.Close()

	// Launch the recording routine
	go sww.RecordEntry(entryChan, recordChan)

	// Init. a new subset object
	s := uniprot.Subset{
		Ekeep: *ekeep, Eskip: *eskip,
		Tkeep: *tkeep, Tskip: *tskip,
		Dkeep: *dkeep, Dskip: *dskip,
		Lmin: *lmin}

	// Launch reading routine(s)
	for _, file := range files {
		// Parsing optimization if reading function is not necessary
		if *dkeep == "" && *dskip == "" {
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
