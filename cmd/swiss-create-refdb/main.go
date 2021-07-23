package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hdevillers/go-fannot/refdb"
)

func main() {
	input := flag.String("i", "", "Input swissProt data file.")
	name := flag.String("n", "", "Name of the reference database.")
	outdir := flag.String("d", ".", "Output directory.")
	flag.Parse()

	if *input == "" {
		panic("You must provide a SwissProt data file.")
	}
	if *name == "" {
		panic("You must provide a name for the new reference database.")
	}

	// Get the absolute path of the output directory
	if !filepath.IsAbs(*outdir) {
		apath, err := filepath.Abs(*outdir)
		if err != nil {
			panic(err)
		}
		*outdir = apath
	}

	// Check if the output directory exists
	_, err := os.Stat(*outdir)
	if os.IsNotExist(err) {
		os.Mkdir(*outdir, 0770)
	}

	// Prepare the refdb root directory
	rootdir := *outdir + "/" + *name
	_, err = os.Stat(rootdir)
	if os.IsNotExist(err) {
		os.Mkdir(rootdir, 0770)
	} else {
		panic("The refdb name is already used in the output directory.")
	}

	rdb := refdb.NewRefdb(*name, *input)
	rdb.WriteJson("Test.json")

	rdb2 := refdb.ReadJson("Test.json")
	fmt.Println(rdb2.Name)
}
