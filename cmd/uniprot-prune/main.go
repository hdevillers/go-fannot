package main

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/hdevillers/go-fannot/custom"
	"github.com/hdevillers/go-fannot/uniprot"
)

// Initialize argument variables
var input string
var output string
var pmeth bool
var pdesc bool
var pfunc bool

// Init function that manages input arguments
func init() {
	// Define default values and usages
	const (
		toolName      = "uniprot-prune"
		toolDesc      = "Prune (discard) UniProt entries according to different criteria"
		inputDefault  = ""
		inputUsage    = "Input UniProt data file (dat(.gz))"
		outputDefault = ""
		outputUsage   = "Output pruned UniProt data file (dat(.gz))"
		pmethDefault  = false
		pmethUsage    = "Prune proteins that do not start by a methionine"
		pdescDefault  = false
		pdescUsage    = "Prune proteins without description"
		pfuncDefault  = false
		pfuncUsage    = "Prune proteins without function information"
	)

	// Init. flags
	flag.StringVar(&input, "input", inputDefault, inputUsage)
	flag.StringVar(&input, "i", inputDefault, inputUsage)
	flag.StringVar(&output, "output", outputDefault, outputUsage)
	flag.StringVar(&output, "o", outputDefault, outputUsage)
	flag.BoolVar(&pmeth, "methionine", pmethDefault, pmethUsage)
	flag.BoolVar(&pmeth, "m", pmethDefault, pmethUsage)
	flag.BoolVar(&pdesc, "description", pdescDefault, pdescUsage)
	flag.BoolVar(&pdesc, "d", pdescDefault, pdescUsage)
	flag.BoolVar(&pfunc, "function", pfuncDefault, pfuncUsage)
	flag.BoolVar(&pfunc, "f", pfuncDefault, pfuncUsage)

	// Shorthand associations
	shand := map[string]string{
		"input":       "i",
		"output":      "o",
		"methionine":  "m",
		"description": "d",
		"function":    "f",
	}

	// Usage print order
	order := []string{"input", "output", "methionine", "description", "function"}

	// Custom usage display
	flag.Usage = func() {
		custom.Usage(*flag.CommandLine, toolName, toolDesc, &order, &shand)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
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

	swr := uniprot.NewReader(input)
	swr.PanicOnError()
	defer swr.Close()

	sww := uniprot.NewWriter(output)
	sww.PanicOnError()
	defer sww.Close()

	tot := 0
	kpt := 0

	// Init. regex
	reMeth := regexp.MustCompile(`^M`)

	for swr.Next() {
		e := swr.Parse()
		tot++

		if pmeth {
			if !reMeth.MatchString(e.Sequence) {
				continue
			}
		}

		if pdesc {
			if e.Desc == "" {
				continue
			}
		}

		if pfunc {
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
