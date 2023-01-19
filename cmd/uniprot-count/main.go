package main

import (
	"flag"
	"fmt"

	"github.com/hdevillers/go-fannot/uniprot"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := flag.String("i", "", "Input UniProt data file.")
	flag.Parse()

	if *input == "" {
		panic("You must provide a UniProt data file.")
	}

	// Create a reader
	swr := uniprot.NewReader(*input)
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
		fmt.Println("Found 1 UniProt entry in", *input)
	} else {
		fmt.Println("Found", cnt, "UniProt entries in", *input)
	}
}
