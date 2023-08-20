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

// Testing hit number
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

// Testing similarity and overwrite
func TestCreateRefdbSimilarityLevels(t *testing.T) {
	inputs := []string{
		"test_highly.fasta",
		"test_similar.fasta",
		"test_weakly.fasta",
	}
	names := []string{
		"TEST_HIGHLY",
		"TEST_SIMILAR",
		"TEST_WEAKLY",
	}
	descs := []string{
		"REFDB to test highly similar hit",
		"REFDB to test similar hit",
		"REFDB to test weakly similar hit",
	}
	indir := "../examples/testdb/"
	outdir := "../examples/refdb/"
	equal := false

	// Create all dbs
	for i := range inputs {
		dbdir_ow := outdir + names[i] + "_OW"
		dbdir_nw := outdir + names[i] + "_NW"

		// Delete possible data from a previous build
		_, err := os.Stat(dbdir_ow)
		if err == nil {
			err = os.RemoveAll(dbdir_ow)
		}
		_, err = os.Stat(dbdir_nw)
		if err == nil {
			err = os.RemoveAll(dbdir_nw)
		}

		// Init. Refdb objects
		rdb_ow := NewRefdb(outdir, names[i]+"_OW", indir+inputs[i], descs[i]+", overwrritable", equal, true, true)
		rdb_nw := NewRefdb(outdir, names[i]+"_NW", indir+inputs[i], descs[i]+", non-overwritable", equal, false, true)

		// Load the data
		rdb_ow.LoadFasta()
		rdb_nw.LoadFasta()

		// Save the json
		rdb_ow.WriteJson()
		rdb_nw.WriteJson()

		// Check if files have been created
		_, err = os.Stat(dbdir_ow + "/config.json")
		if err != nil {
			t.Fatal("The Refdb " + dbdir_ow + " config file (json) does not exist.")
		}
		_, err = os.Stat(dbdir_nw + "/config.json")
		if err != nil {
			t.Fatal("The Refdb " + dbdir_nw + " config file (json) does not exist.")
		}
	}
}

// Additional DB for testing overwrite
func TestCreateRefdbOverwrite(t *testing.T) {
	inputs := []string{
		"test_ow1.fasta",
		"test_ow2.fasta",
		"test_ow3.fasta",
	}
	names := []string{
		"TEST_OW1",
		"TEST_OW2",
		"TEST_OW3",
	}
	descs := []string{
		"REFDB to test ow similarity 51",
		"REFDB to test ow similarity 58",
		"REFDB to test ow similarity 74",
	}
	indir := "../examples/testdb/"
	outdir := "../examples/refdb/"
	equal := false

	// Create all dbs
	for i := range inputs {
		dbdir_ow := outdir + names[i]

		// Delete possible data from a previous build
		_, err := os.Stat(dbdir_ow)
		if err == nil {
			err = os.RemoveAll(dbdir_ow)
		}

		// Init. Refdb objects
		rdb_ow := NewRefdb(outdir, names[i], indir+inputs[i], descs[i], equal, true, true)

		// Load the data
		rdb_ow.LoadFasta()

		// Save the json
		rdb_ow.WriteJson()

		// Check if files have been created
		_, err = os.Stat(dbdir_ow + "/config.json")
		if err != nil {
			t.Fatal("The Refdb " + dbdir_ow + " config file (json) does not exist.")
		}
	}
}
