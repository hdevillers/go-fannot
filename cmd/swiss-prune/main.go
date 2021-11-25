package main

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/hdevillers/go-fannot/swiss"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := flag.String("i", "", "Input SwissProt data file.")
	output := flag.String("o", "", "Output pruned SwissPort data file.")
	pmeth := flag.Bool("m", false, "Prune proteins that do not start by a Methionine.")
	pdesc := flag.Bool("d", false, "Prune proteins without description.")
	pfunc := flag.Bool("f", false, "Prune proteins without function information.")
	flag.Parse()

	if *input == "" {
		panic("You must provide a SwissProt data file.")
	}

	if *output == "" {
		panic("You must provide an output file name.")
	}

	swr := swiss.NewReader(*input)
	swr.PanicOnError()
	defer swr.Close()

	sww := swiss.NewWriter(*output)
	sww.PanicOnError()
	defer sww.Close()

	tot := 0
	kpt := 0

	// Init. regex
	reMeth := regexp.MustCompile(`^M`)

	for swr.Next() {
		e := swr.Parse()
		tot++

		if *pmeth {
			if !reMeth.MatchString(e.Sequence) {
				continue
			}
		}

		if *pdesc {
			if e.Desc == "" {
				continue
			}
		}

		if *pfunc {
			if e.Function == "" {
				continue
			}
		}
		kpt++
		sww.WriteStrings(swr.GetData())
		sww.WriteEntryEnd()
		sww.PanicOnError()
	}

	fmt.Println("Scan", tot, "entries and kept", kpt, "ones.")
}
