package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hdevillers/go-fannot/swiss"
)

type Subset struct {
	Ekeep string
	Eskip string
	Tkeep string
	Tskip string
	Lmin  int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Recorder routine
func recordEntry(ec chan *[]string, re chan int, out string) {
	// Create a writer
	sww := swiss.NewWriter(out)
	sww.PanicOnError()
	defer sww.Close()

	nrec := 0

	for e := range ec {
		sww.WriteStrings(e)
		sww.WriteEntryEnd()
		nrec++
	}

	// Throw the number of recorded entries
	re <- nrec
}

func (s *Subset) parseFile(ec chan *[]string, th chan int, in string) {
	// Create a reader
	swr := swiss.NewReader(in)
	swr.PanicOnError()
	defer swr.Close()

	ntot := 0

	for swr.Next() {
		// Parse the entry
		e := swr.LightParse()
		ntot++

		if e.Length < s.Lmin {
			continue
		}

		if s.Eskip != "" {
			if e.TestEvidence(s.Eskip) {
				continue
			}
		}

		if s.Tskip != "" {
			if e.TestTaxonomy(s.Tskip) {
				continue
			}
		}

		if s.Ekeep != "" {
			if !e.TestEvidence(s.Ekeep) {
				continue
			}
		}

		if s.Tkeep != "" {
			if !e.TestTaxonomy(s.Tkeep) {
				continue
			}
		}
		var dt []string
		dt = *swr.GetData()

		ec <- &dt
	}

	// Throw the number of scanned entries
	th <- ntot
}

func main() {
	input := flag.String("i", "", "Input SwissProt data file.")
	output := flag.String("o", "", "Output data file.")
	ekeep := flag.String("e", "", "Evident keep instruction (regex).")
	eskip := flag.String("E", "", "Evidence skip instruction (regex).")
	tkeep := flag.String("t", "", "Taxonomy keep instruction (regex).")
	tskip := flag.String("T", "", "Taxonomy skip instruction (regex).")
	lmin := flag.Int("l", 30, "Minimal protein length (aa).")
	flag.Parse()

	if *input == "" {
		panic("You must provide a SwissProt data file.")
	}

	if *output == "" {
		panic("You must provide an output file name.")
	}

	if *ekeep == "" && *eskip == "" && *tkeep == "" && *tskip == "" {
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

	// Launch the recording routine
	go recordEntry(entryChan, recordChan, *output)

	// Init. a new subset object
	s := Subset{*ekeep, *eskip, *tkeep, *tskip, *lmin}

	// Launch reading routine(s)
	for _, file := range files {
		go s.parseFile(entryChan, threadChan, file)
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
