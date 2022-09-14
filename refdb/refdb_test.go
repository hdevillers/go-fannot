package refdb

import (
	"os"
	"testing"
)

func TestCreateRefdb(t *testing.T) {
	input := "../examples/swiss/S288c_lipase.dat.gz"
	name := "S288c_lipase"
	outdir := "../examples/refdb/"
	desc := "Lipase proteins from Saccharomyces cerevisiae S288c."

	// Init. Refdb object
	rdb := NewRefdb(outdir, name, input, desc, true, false, true, true)

	// Load the data
	rdb.LoadSource()

	// Save the json
	rdb.WriteJson()

	// Check if files have been created
	_, err := os.Stat("../examples/refdb/S288c_lipase/config.json")
	if err != nil {
		t.Fatal("The Refdb config file (json) does not exist.")
	}
}
