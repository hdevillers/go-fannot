package refdb

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/hdevillers/go-fannot/uniprot"
	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/seqio"
)

const (
	FASTA_PATH   string = "protein.fasta"
	BLASTDB_PATH string = "blastdb"
	JSON_PATH    string = "config.json"
)

type Refdb struct {
	Id        string
	Desc      string
	Root      string
	Source    string
	Blastdb   string
	Fasta     string
	Nprot     int
	Equal     bool // Indicate if the DB contain proteins of the query
	OverWrite bool // Indicate if annotations from the DB can overwrite "similar" annotations
	Reviewed  bool // Indicate if the DB is reviewed (Uniprot) or not (TrEmbl)
	GeneName  bool // Indicate if we can transfer gene name in the query feature
}

func NewRefdb(outdir, id, source, desc string, equal bool, ow bool, re bool, gn bool) *Refdb {
	var rdb Refdb

	// Check if the source exist
	_, err := os.Stat(source)
	if os.IsNotExist(err) {
		panic("The input source data file does not exists.")
	}

	// Check if the output directory exists
	_, err = os.Stat(outdir)
	if os.IsNotExist(err) {
		err = os.Mkdir(outdir, 770)
		if err == nil {
			panic(err)
		}
	}

	// Turn outdir into ablsolute path (if necessary)
	if !filepath.IsAbs(outdir) {
		apath, err := filepath.Abs(outdir)
		if err != nil {
			panic(err)
		}
		outdir = apath
	}

	// Prepare the root directory
	rootdir := outdir + "/" + id
	_, err = os.Stat(rootdir)
	if os.IsNotExist(err) {
		os.Mkdir(rootdir, 0770)
	} else {
		panic("The refdb name is already used in the output directory.")
	}

	// Setup path values
	rdb.Id = id
	rdb.Root = rootdir
	rdb.Source = source
	rdb.Desc = desc
	rdb.Equal = equal
	rdb.OverWrite = ow
	rdb.Reviewed = re
	rdb.GeneName = gn

	return &rdb
}

func (r *Refdb) LoadSource() {
	// Init. the uniprot reader
	swr := uniprot.NewReader(r.Source)
	swr.PanicOnError()
	defer swr.Close()

	// Init. the fasta writer
	r.Fasta = r.Root + "/" + FASTA_PATH
	fw := seqio.NewWriter(r.Fasta, "fasta", false)

	// Scan entries
	ne := 0
	for swr.Next() {
		e := swr.Parse()
		swr.PanicOnError()
		ne++

		desc := e.Desc + "::" + e.Name + "::" + e.Locus + "::" + e.Organism + "::" + e.Function
		nseq := seq.NewSeq(e.Access)
		nseq.Desc = desc
		nseq.Sequence = []byte(e.Sequence)

		fw.Write(*nseq)
	}
	r.Nprot = ne

	// Prepare the BLASTDB
	r.Blastdb = r.Root + "/" + BLASTDB_PATH
	err := exec.Command("makeblastdb",
		"-in", r.Fasta,
		"-out", r.Blastdb,
		"-dbtype", "prot",
		"-input_type", "fasta",
	).Run()
	if err != nil {
		panic(err)
	}
}

func (r *Refdb) PrintInfoHeader() {
	fmt.Println("ID\t#Proteins\tDescription")
}

func (r *Refdb) PrintInfo() {
	fmt.Printf("%s\t%d\t%s\n", r.Id, r.Nprot, r.Desc)
}

// Create a json file from an existing object
func (r *Refdb) WriteJson() {
	// Create the output file
	f, err := os.Create(r.Root + "/" + JSON_PATH)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create the writer object
	fw := bufio.NewWriter(f)

	// Create the json encoder
	jw := json.NewEncoder(fw)

	// encode
	err = jw.Encode(r)
	if err != nil {
		panic(err)
	}

	fw.Flush()
}

// create a Refdb object from a json file
func ReadJson(file string) *Refdb {
	var refdb Refdb

	// Open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create the reader
	fr := bufio.NewReader(f)

	// Create the json decoder
	jr := json.NewDecoder(fr)

	// Decode the entry
	err = jr.Decode(&refdb)
	if err != nil {
		panic(err)
	}

	return &refdb
}

// Create a Refdb object from an id and a directory
func FindRefDB(id, dir string) *Refdb {
	// Check if the provided id is a JSON file
	tjson := regexp.MustCompile(`\.json$`)
	if tjson.MatchString(id) {
		_, err := os.Stat(id)
		json := id
		if os.IsNotExist(err) {
			_, err := os.Stat(dir + "/" + id)
			if os.IsNotExist(err) {
				panic(fmt.Sprintf("Failed to find the JSON file (%s/)%s.", dir, id))
			} else {
				json = dir + "/" + id
			}
		}
		return ReadJson(json)
	} else {
		// The ID is probably a real ID
		json := dir + "/" + id + "/" + JSON_PATH
		_, err := os.Stat(json)
		if os.IsNotExist(err) {
			// Try without the directory
			json = id + "/" + JSON_PATH
			_, err := os.Stat(json)
			if os.IsNotExist(err) {
				panic(fmt.Sprintf("Failed to find the DB with ID: %s (directory: %s).", id, dir))
			}
		}
		return ReadJson(json)
	}
}
