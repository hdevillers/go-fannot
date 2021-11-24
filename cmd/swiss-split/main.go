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
	output := flag.String("o", "", "Output file basename.")
	nsplit := flag.Int("n", 10, "Number of sub-data files wanted.")
	compress := flag.Bool("c", false, "Compress output files.")
	flag.Parse()

	if *input == "" {
		panic("You must provide a SwissProt data file.")
	}

	if *output == "" {
		panic("You must provide an output file basename.")
	}

	if *nsplit < 2 {
		panic("The number of file division must be greater than 1.")
	}

	swr := swiss.NewReader(*input)
	swr.PanicOnError()
	defer swr.Close()

	writers := make([]*swiss.Writer, *nsplit)
	fileExt := ".dat"
	if *compress {
		fileExt += ".gz"
	}

	for i := 0; i < *nsplit; i++ {
		writers[i] = swiss.NewWriter(*output + fmt.Sprintf("%03d", i) + fileExt)
		writers[i].PanicOnError()
		defer writers[i].Close()
	}

	wi := 0
	for swr.Next() {
		writers[wi].WriteStrings(swr.GetData())
		writers[wi].WriteEntryEnd()
		writers[wi].PanicOnError()
		wi++
		if wi == *nsplit {
			wi = 0
		}
	}

}
