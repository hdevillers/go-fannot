package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hdevillers/go-fannot/custom"
	"github.com/hdevillers/go-fannot/uniprot"
)

// Initialize argument variables
var src string
var listSrc bool
var mir string
var listMir bool
var div string
var listDiv bool
var dir string
var skipSum bool

func init() {
	// Define default values and usages
	const (
		toolName       = "uniprot-download"
		toolDesc       = "Download UniProt database for a given taxonomic division"
		srcDefault     = "sprot"
		srcUsage       = "Source database"
		listSrcDefault = false
		listSrcUsage   = "List available sources"
		mirDefault     = "us"
		mirUsage       = "Download mirror"
		listMirDefault = false
		listMirUsage   = "List available download mirrors"
		divDefault     = ""
		divUsage       = "Taxonomic division"
		listDivDefault = false
		listDivUsage   = "List available taxonomic divisions"
		dirDefault     = "."
		dirUsage       = "Output directory"
		skipSumDefault = false
		skipSumUsage   = "Do not make the checksum validation after download"
	)

	// Init. flags
	flag.StringVar(&src, "source", srcDefault, srcUsage)
	flag.StringVar(&src, "s", srcDefault, srcUsage)
	flag.StringVar(&mir, "mirror", mirDefault, mirUsage)
	flag.StringVar(&mir, "m", mirDefault, mirUsage)
	flag.StringVar(&div, "division", divDefault, divUsage)
	flag.StringVar(&div, "d", divDefault, divUsage)
	flag.StringVar(&dir, "output-dir", dirDefault, dirUsage)
	flag.StringVar(&dir, "o", dirDefault, dirUsage)
	flag.BoolVar(&listSrc, "list-sources", listSrcDefault, listSrcUsage)
	flag.BoolVar(&listSrc, "S", listSrcDefault, listSrcUsage)
	flag.BoolVar(&listMir, "list-mirrors", listMirDefault, listMirUsage)
	flag.BoolVar(&listMir, "M", listMirDefault, listMirUsage)
	flag.BoolVar(&listDiv, "list-divisions", listDivDefault, listDivUsage)
	flag.BoolVar(&listDiv, "D", listDivDefault, listDivUsage)
	flag.BoolVar(&skipSum, "skip-sum", skipSumDefault, skipSumUsage)

	// Shorthand associations
	shand := map[string]string{
		"source":         "s",
		"mirror":         "m",
		"division":       "d",
		"output-dir":     "o",
		"list-sources":   "S",
		"list-mirrors":   "M",
		"list-divisions": "D",
	}

	order := []string{"source", "mirror", "division", "output-dir", "list-sources", "list-mirrors", "list-divisions", "skip-sum"}

	// Custom usage display
	flag.Usage = func() {
		custom.Usage(*flag.CommandLine, toolName, toolDesc, &order, &shand)
	}
}

func main() {
	flag.Parse()

	// Check directory
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0750)
		if err != nil {
			panic(err)
		}
	}

	// Init. metalink object
	metalink := uniprot.NewMetalink()

	// List databases if required
	if listSrc {
		fmt.Println("Available databases:")
		for k := range metalink.Databases {
			fmt.Println("  - ", k)
		}
		return
	}
	// List mirrors if required
	if listMir {
		fmt.Println("Available mirrors:")
		for k := range metalink.Mirrors {
			fmt.Println("  - ", k)
		}
		return
	}

	// Look for missing division
	if div == "" && !listDiv {
		panic("You must provide a taxonomic division. Use -list-div to get the list of available divisions.")
	}

	// Check if the required database is available
	if !metalink.CheckDatabase(src) {
		panic("The required database is not available. Use -list-src to get the list of available databases.")
	}
	// Check if the required mirror is available
	if !metalink.CheckMirror(mir) {
		panic("The required mirror is not available. Use -list-mir to get the list of available mirrors.")
	}

	// NOTE: Divisions are retrieved from the metalink file,
	// hence, data must be retrieved before checking arguments.
	fmt.Println("Retrieving Metlink file from UniProt...")
	metalink.Retrieve(mir)

	// List division if required
	if listDiv {
		fmt.Println("Available taxonomic divisions:")
		for k := range metalink.Divisions {
			fmt.Println("  - ", k)
		}
		return
	}

	// Check if the required division is available
	if !metalink.CheckDivision(div) {
		panic("The required division is not available. Use -list-div to get the list of available divisions.")
	}

	// Prepare output
	odir := fmt.Sprintf("%s/uniprot_v%s/", dir, metalink.Version)
	if _, err := os.Stat(odir); os.IsNotExist(err) {
		err = os.Mkdir(odir, 0750)
		if err != nil {
			panic(err)
		}
	}

	// Find the URL
	fmt.Println("Finding the proper URL...")
	url, sum := metalink.GetUrl(src, div, mir)

	// Destination file
	dest := fmt.Sprintf("%suniprot_%s_%s.dat.gz", odir, src, div)

	// Check if the file already exists
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		fmt.Printf("The destination file (%s) already exists. Delete it or change output directory.\n", dest)
		return
	}

	// Launch FTP download
	uniprot.FtpDownload(url, sum, dest, skipSum)
}
