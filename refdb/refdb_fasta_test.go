package refdb

import (
	"os"
	"os/exec"
	"testing"
)

// Create a reference fasta from an annotated sequence file
func TestExtractRefFasta(t *testing.T) {
	if os.Getenv("DO_PYTHON_TESTS") == "yes" {
		input := "../examples/genbank/CALSDN010000023.1.gb"
		output := "../examples/from_fasta/test.fasta"

		// Remove possible data from previous build
		_, err := os.Stat(output)
		if err == nil {
			err = os.Remove(output)
			if err != nil {
				t.Fatal("Failed to remove data from previous tests.")
			}
		}

		cmd := exec.Command("python3", "../scripts/SeqFilesToRefFasta.py", "-s", input, "-f", "genbank", "-o", output)
		err = cmd.Run()
		if err != nil {
			t.Fatal("Failed to run python script SeqFilesToRefFasta.py: " + err.Error())
		}

		// Check if the output exists
		_, err = os.Stat(output)
		if err != nil {
			t.Fatal("Failed to find output reference fasta file.")
		}
	} else {
		// The test is skipped if DO_PYTHON_TEST is different from "yes" or unset
		t.Skip("Tests using Python are disabled by default, to enable it, use 'make test do_python_tests=yes'.")
	}
}

// Convert the reference fasta into a reference DB
func TestRefFastaToRefDB(t *testing.T) {
	if os.Getenv("DO_PYTHON_TESTS") == "yes" {
		input := "../examples/from_fasta/test.fasta"

		_, err := os.Stat(input)
		if err != nil {
			t.Fatal("Failed to find input fasta file.")
		}

	} else {
		// The test is skipped if DO_PYTHON_TEST is different from "yes" or unset
		t.Skip("Tests using Python are disabled by default, to enable it, use 'make test do_python_tests=yes'.")
	}
}
