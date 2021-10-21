package fannot

import (
	"fmt"
	"strings"

	"github.com/hdevillers/go-blast"
	"github.com/hdevillers/go-fannot/refdb"
	"github.com/hdevillers/go-needle"
	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/utils"
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
		"Null",
		"Null",
	}
}

// Functional annotation main structure
type Fannot struct {
	Queries    []seq.Seq
	NQueries   int
	DBs        []refdb.Refdb
	DBi        int
	DBEntries  map[string]seq.Seq
	Continue   []bool
	Results    []FAResult
	BlastPar   blast.Param
	NeedlePar  needle.Param
	NBestHits  int
	MinCovHigh float64
	MinSimHigh float64
	MinCovNorm float64
	MinSimNorm float64
}

func NewFannot(i string) *Fannot {
	var fa Fannot

	// Load the query sequences
	fa.NQueries = utils.LoadSeqInArray(i, "fasta", &fa.Queries)

	// Init. BLAST and NEEDLE parameter setings
	fa.BlastPar = *blast.NewParam()
	fa.NeedlePar = *needle.NewParam()

	// Init. results and continue variables
	fa.Continue = make([]bool, fa.NQueries)
	fa.Results = make([]FAResult, fa.NQueries)
	for i := 0; i < fa.NQueries; i++ {
		// Set default result (when no match is found)
		fa.Results[i] = *NewFAResult()
	}

	// Setup default threshold
	fa.NBestHits = N_BEST_HITS
	fa.MinCovHigh = MIN_COV_HIGH
	fa.MinSimHigh = MIN_SIM_HIGH
	fa.MinCovNorm = MIN_COV_NORM
	fa.MinSimNorm = MIN_SIM_NORM

	return &fa
}

func (fa *Fannot) GetDBs(i, d string) {
	// split ids
	ids := strings.Split(i, ",")

	// Empty current DBs if necessary
	fa.DBs = make([]refdb.Refdb, 0)
	fa.DBi = -1

	// Fill with found DB
	for _, id := range ids {
		newDB := refdb.FindRefDB(id, d)
		fa.DBs = append(fa.DBs, *newDB)
	}
}

func (fa *Fannot) LoadDBEntries() {
	// Load DB entries (FASTA)
	utils.LoadSeqInMap(fa.DBs[fa.DBi].Fasta, "fasta", &fa.DBEntries)
}

func (fa *Fannot) NextDB() bool {
	fa.DBEntries = make(map[string]seq.Seq)
	fa.DBi++
	if fa.DBi < len(fa.DBs) {
		fa.LoadDBEntries()
		return true
	} else {
		return false
	}
}

// Go-routine that treat one given gene
func (fa *Fannot) FindFunction(queryChan chan int, threadChan chan int) {
	// Init. search tools
	blt := blast.NewBlast()
	blt.Par = &fa.BlastPar
	blt.Db = fa.DBs[fa.DBi].Blastdb

	// Get the query id(s) from the chan
	for qi := range queryChan {
		/* First step: BLAST */

		// Add the query and run blast
		blt.AddQuery(fa.Queries[qi])
		err := blt.Search()
		if err != nil {
			panic(err)
		}

		// Parse blast result
		results := blt.Rst.Iterations[0]
		if len(results.Hits) > 0 {
			chkhit := 0 // Number if hit checked
			bestHitId := "NULL"
			bestHitSim := 0.0
			var bestHitNdl needle.Result
		HITS:
			for _, hit := range results.Hits {
				// For each hit, compute the global alignment and extract the similarity
				hitId := hit.HitDef
				hitSeq, test := fa.DBEntries[hitId]
				if !test {
					panic(fmt.Sprintf("Failed to find the hit %s in the reference DB (%s).", hitId, fa.DBs[fa.DBi].Id))
				}
				ndl := needle.NewNeedle(fa.Queries[qi], hitSeq)
				ndl.Par = &fa.NeedlePar
				err = ndl.Align()
				if err != nil {
					panic(err)
				}

				if ndl.Rst.GetSimilarityPct() > bestHitSim {
					bestHitId = hitId
					bestHitNdl = *ndl.Rst
					bestHitSim = bestHitNdl.GetSimilarityPct()
				}

				chkhit++
				if chkhit >= fa.NBestHits {
					break HITS
				}
			}
		}
		// Else do nothing

		// Remove the query
		blt.ResetQuery()
	}

	// Terminate the thread
	threadChan <- 1
}
