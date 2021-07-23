package refdb

import (
	"bufio"
	"encoding/json"
	"os"
)

type Refdb struct {
	Name    string
	Source  string
	Blastdb string
	Fasta   string
	Details string
	Nprot   int
}

func NewRefdb(name string, source string) *Refdb {
	return &Refdb{
		Name:   name,
		Source: source,
	}
}

// Create a json file from an existing object
func (r *Refdb) WriteJson(file string) {
	// Create the output file
	f, err := os.Create(file)
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
