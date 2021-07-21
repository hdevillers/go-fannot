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
	flag.Parse()

	if *input == "" {
		panic("You must provide a SwissProt data file.")
	}

	// Create a reader
	swr := swiss.NewReader(*input)
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
		fmt.Println("Found 1 SwissProt entry in", *input)
	} else {
		fmt.Println("Found", cnt, "SwissProt entries in", *input)
	}
}
