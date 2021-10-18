package fannot

import (
	"github.com/hdevillers/go-fannot/refdb"
	"github.com/hdevillers/go-seq/seq"
)

// Initialize default thresholds
const (
	N_BEST_HITS  int     = 3
	MIN_COV_HIGH float64 = 90.0
	MIN_SIM_HIGH float64 = 80.0
	MIN_COV_NORM float64 = 70.0
	MIN_SIM_NORM float64 = 50.0
)

// DEFINING STRUCTURES

// Functional Annotation Results
type FAResult struct {
	Product string
	Note    string
	GeneID  string
	RefID   string
}

func NewFAResult() *FAResult {
	return &FAResult{
		"Hypothetical protein",
		"Hypothetical protein",
		"",
		"",
	}
}

// Single reference entry
type DBEntry struct {
	Seq      seq.Seq
	Species  string
	Function string
}

// Functional annotation main structure
type Fannot struct {
	Queries []seq.Seq
	DBs     []refdb.Refdb
}
