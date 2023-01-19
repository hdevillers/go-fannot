package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hdevillers/go-fannot/uniprot"
)

func main() {
	db := flag.String("db", "sprot", "Source database.")
	listDb := flag.Bool("list-db", false, "List available databases.")
	mir := flag.String("mir", "us", "Download mirror.")
	listMir := flag.Bool("list-mir", false, "List available mirrors.")
	div := flag.String("div", "", "Taxonomic division.")
	listDiv := flag.Bool("list-div", false, "List available taxonomic divisions.")
	dir := flag.String("dir", ".", "Output directory.")
	flag.Parse()

	// Check directory
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		err = os.Mkdir(*dir, 0750)
		if err != nil {
			panic(err)
		}
	}

	// Init. metalink object
	metalink := uniprot.NewMetalink()

	// List databases if required
	if *listDb {
		fmt.Println("Available databases:")
		for k := range metalink.Databases {
			fmt.Println("  - ", k)
		}
		return
	}
	// List mirrors if required
	if *listMir {
		fmt.Println("Available mirrors:")
		for k := range metalink.Mirrors {
			fmt.Println("  - ", k)
		}
		return
	}

	// Look for missing division
	if *div == "" && !*listDiv {
		panic("You must provide a taxonomic division. Use -list-div to get the list of available divisions.")
	}

	// Check if the required database is available
	if !metalink.CheckDatabase(*db) {
		panic("The required database is not available. Use -list-db to get the list of available databases.")
	}
	// Check if the required mirror is available
	if !metalink.CheckMirror(*mir) {
		panic("The required mirror is not available. Use -list-mir to get the list of available mirrors.")
	}

	// NOTE: Divisions are retrieved from the metalink file,
	// hence, data must be retrieved before checking arguments.
	fmt.Println("Retrieving Metlink file from UniProt...")
	metalink.Retrieve()

	// List division if required
	if *listDiv {
		fmt.Println("Available taxonomic divisions:")
		for k := range metalink.Divisions {
			fmt.Println("  - ", k)
		}
		return
	}

	// Check if the required division is available
	if !metalink.CheckDivision(*div) {
		panic("The required division is not available. Use -list-div to get the list of available divisions.")
	}

	// Prepare output
	odir := fmt.Sprintf("%s/uniprot_v%s/", *dir, metalink.Version)
	if _, err := os.Stat(odir); os.IsNotExist(err) {
		err = os.Mkdir(odir, 0750)
		if err != nil {
			panic(err)
		}
	}

	// Find the URL
	fmt.Println("Finding the proper URL...")
	url, sum := metalink.GetUrl(*db, *div, *mir)

	// Replace ftp by https

	// Destination file
	dest := fmt.Sprintf("%suniprot_%s_%s.dat.gz", odir, *db, *div)

	// Check if the file already exists
	if _, err := os.Stat(dest); os.IsExist(err) {
		fmt.Printf("The destination file (%s) already exists. Delete it or change output directory.\n", dest)
		return
	}

	// Create it
	out, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Download the data
	fmt.Println("Downloading the required file...")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("The request returned a bad status: %s\n", resp.Status))
	}

	// Check MD5
	fmt.Println("Check sum (MD5) the downloaded file...")
	hash := md5.New()
	io.Copy(hash, resp.Body)
	// Obtained sum
	osum := hash.Sum(nil)

	if string(osum) != sum {
		panic(fmt.Sprintf("Check sum failed,\nexpected: %s,\n obtained: %s", sum, string(osum)))
	}

	// Write out the file
	fmt.Println("Writing out the file...")
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
