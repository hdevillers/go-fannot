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

// Create the different db for functional tests
// Testing equality
func TestCreateRefdbEqual(t *testing.T) {
	input := "../examples/testdb/test_equal.fasta"
	name := "TEST_EQUAL"
	outdir := "../examples/refdb/"
	desc := "REFDB to test equality."
	equal := true
	overw := false

	// Delete possible data from a previous build
	_, err := os.Stat(outdir + name)
	if err == nil {
		err = os.RemoveAll(outdir + name)
		if err != nil {
			t.Fatal("Failed to remove data from previous build.")
		}
	}

	// Init. Refdb object
	rdb := NewRefdb(outdir, name, input, desc, equal, overw, true)

	// Load the data
	rdb.LoadFasta()

	// Save the json
	rdb.WriteJson()

	// Check if files have been created
	_, err = os.Stat("../examples/refdb/TEST_EQUAL/config.json")
	if err != nil {
		t.Fatal("The Refdb config file (json) does not exist.")
	}
}

// Testing non-equality
func TestCreateRefdbNotEqual(t *testing.T) {
	input := "../examples/testdb/test_equal.fasta"
	name := "TEST_NOT_EQUAL"
	outdir := "../examples/refdb/"
	desc := "REFDB to test non-equality."
	equal := false
	overw := false

	// Delete possible data from a previous build
	_, err := os.Stat(outdir + name)
	if err == nil {
		err = os.RemoveAll(outdir + name)
		if err != nil {
			t.Fatal("Failed to remove data from previous build.")
		}
	}

	// Init. Refdb object
	rdb := NewRefdb(outdir, name, input, desc, equal, overw, true)

	// Load the data
	rdb.LoadFasta()

	// Save the json
	rdb.WriteJson()

	// Check if files have been created
	_, err = os.Stat("../examples/refdb/TEST_NOT_EQUAL/config.json")
	if err != nil {
		t.Fatal("The Refdb config file (json) does not exist.")
	}
}

// Testing non-equality
func TestCreateRefdbHitNumber(t *testing.T) {
	input := "../examples/testdb/test_hitnumber.fasta"
	name := "TEST_HIT_NUMBER"
	outdir := "../examples/refdb/"
	desc := "REFDB to test hit number."
	equal := false
	overw := false

	// Delete possible data from a previous build
	_, err := os.Stat(outdir + name)
	if err == nil {
		err = os.RemoveAll(outdir + name)
		if err != nil {
			t.Fatal("Failed to remove data from previous build.")
		}
	}

	// Init. Refdb object
	rdb := NewRefdb(outdir, name, input, desc, equal, overw, true)

	// Load the data
	rdb.LoadFasta()

	// Save the json
	rdb.WriteJson()

	// Check if files have been created
	_, err = os.Stat("../examples/refdb/TEST_HIT_NUMBER/config.json")
	if err != nil {
		t.Fatal("The Refdb config file (json) does not exist.")
	}
}
