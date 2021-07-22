package main

import (
	"flag"
	"fmt"

	"github.com/hdevillers/go-fannot/swiss"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := flag.String("i", "", "Input SwissProt data file.")
	ekeep := flag.String("e", "", "Evident keep instruction (regex).")
	eskip := flag.String("E", "", "Evidence skip instruction (regex).")
	tkeep := flag.String("t", "", "Taxonomy keep instruction (regex).")
	tskip := flag.String("T", "", "Taxonomy skip instruction (regex).")
	flag.Parse()

	if *input == "" {
		panic("You must provide a SwissProt data file.")
	}

	if *ekeep == "" && *eskip == "" && *tkeep == "" && *tskip == "" {
		panic("You must provide at least one keep/skip instruction.")
	}

	// Create a reader
	swr := swiss.NewReader(*input)
	defer swr.Close()

	// Init. entry counter
	tot := 0
	kpt := 0

	for swr.Next() {
		// Parse the entry
		e := swr.Parse()
		tot++

		if *eskip != "" {
			if e.TestEvidence(*eskip) {
				continue
			}
		}

		if *tskip != "" {
			if e.TestTaxonomy(*tskip) {
				continue
			}
		}

		if *ekeep != "" {
			if !e.TestEvidence(*ekeep) {
				continue
			}
		}

		if *tkeep != "" {
			if !e.TestTaxonomy(*tkeep) {
				continue
			}
		}

		// At that step, consider the entry as keepable!
		kpt++

	}

	fmt.Println("Scan", tot, "entries and kept", kpt, "ones.")
}
