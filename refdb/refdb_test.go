package refdb

import (
	"os"
	"testing"
)

func TestCreateRefdb(t *testing.T) {
	input := "../examples/uniprot/S288c_lipase.dat.gz"
	name := "S288c_lipase"
	outdir := "../examples/refdb/"
	desc := "Lipase proteins from Saccharomyces cerevisiae S288c."

	// Delete possible data from a previous duild
	_, err := os.Stat(outdir + name)
	if err == nil {
		err = os.RemoveAll(outdir + name)
		if err != nil {
			t.Fatal("Failed to remove data from previsous build.")
		}
	}

	// Init. Refdb object
	rdb := NewRefdb(outdir, name, input, desc, true, false, true)

	// Load the data
	rdb.LoadSource()

	// Save the json
	rdb.WriteJson()

	// Check if files have been created
	_, err = os.Stat("../examples/refdb/S288c_lipase/config.json")
	if err != nil {
		t.Fatal("The Refdb config file (json) does not exist.")
	}
}
