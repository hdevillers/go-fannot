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
	MIN_LRA_HIGH float64 = 0.9
	MIN_SIM_HIGH float64 = 80.0
	MIN_LRA_NORM float64 = 0.7
	MIN_SIM_NORM float64 = 50.0
	UNKNOWN_FUNC string  = "hypothetical protein"
	PRE_SIM_HIGH string  = "higly similar to"
	PRE_SIM_NORM string  = "similar to"
)

// DEFINING STRUCTURES

// Functional Annotation Results
type FAResult struct {
	Product  string
	Note     string
	Locus    string
	Name     string
	Status   int
	Organism string
	GeneID   string
	RefID    string
	HitSim   float64
	HitLR    float64
}

func NewFAResult() *FAResult {
	return &FAResult{
		UNKNOWN_FUNC,
		UNKNOWN_FUNC,
		"Null",
		"Null",
		0,
		"Null",
		"Null",
		"Null",
		0.0,
		0.0,
	}
}

func ParseHitDesc(hd string, hid string, rid string, hs int, eq bool) *FAResult {
	var far FAResult
	values := strings.Split(hd, "::")

	far.Product = values[0]
	far.Organism = values[3]
	far.Status = hs
	far.GeneID = hid
	far.RefID = rid

	if eq {
		far.Note = fmt.Sprintf("uniprot|%s %s", far.GeneID, far.Organism)
	} else if hs == 2 {
		far.Note = fmt.Sprintf("%s uniprot|%s %s", PRE_SIM_HIGH, far.GeneID, far.Organism)
	} else {
		far.Note = fmt.Sprintf("%s uniprot|%s %s", PRE_SIM_NORM, far.GeneID, far.Organism)
	}
	if values[2] != "" {
		far.Note += " " + values[2]
		far.Locus = values[2]
	}
	if values[1] != "" {
		far.Note += " " + values[1]
		far.Name = values[1]
	}
	far.Note += " " + far.Product

	return &far
}

func (far *FAResult) PrintFAResult(gid string) {
	fmt.Printf(
		"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%.03f\t%.03f\t%s\n",
		gid, far.Product, far.Note, far.Organism,
		far.GeneID, far.Locus, far.Name, far.Status,
		far.HitSim, far.HitLR, far.RefID,
	)
}

// UTILS

// Return the minimal length ratio
func getMinLengthRatio(l1, l2 int) float64 {
	if l1 < l2 {
		return float64(l1) / float64(l2)
	} else {
		return float64(l2) / float64(l1)
	}
}

// Print functional annotation table header
func PrintFAResultsHeader() {
	fmt.Println("GeneID\tProduct\tNote\tOrganism\tRefID\tRefLocus\tRefName\tStatus\tSimilarity\tLengthRatio\tDBID")
}

// Functional annotation main structure
type Fannot struct {
	Queries    []seq.Seq
	NQueries   int
	DBs        []refdb.Refdb
	DBi        int
	DBEntries  map[string]seq.Seq
	Finished   []bool
	Results    []FAResult
	BlastPar   blast.Param
	NeedlePar  needle.Param
	NBestHits  int
	MinLraHigh float64
	MinSimHigh float64
	MinLraNorm float64
	MinSimNorm float64
}

func NewFannot(i string) *Fannot {
	var fa Fannot

	// Load the query sequences
	fa.NQueries = utils.LoadSeqInArray(i, "fasta", &fa.Queries)

	// Init. BLAST and NEEDLE parameter setings
	fa.BlastPar = *blast.NewParam()
	fa.NeedlePar = *needle.NewParam()

	// Init. results and Finished variables
	fa.Finished = make([]bool, fa.NQueries)
	fa.Results = make([]FAResult, fa.NQueries)
	for i := 0; i < fa.NQueries; i++ {
		// Set default result (when no match is found)
		fa.Results[i] = *NewFAResult()
	}

	// Setup default threshold
	fa.NBestHits = N_BEST_HITS
	fa.MinLraHigh = MIN_LRA_HIGH
	fa.MinSimHigh = MIN_SIM_HIGH
	fa.MinLraNorm = MIN_LRA_NORM
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
			bestHitDesc := ""
			bestHitSim := 0.0
			bestHitLen := 0
			bestHitStatus := 0
		HITS:
			for _, hit := range results.Hits {
				// For each hit, compute the global alignment and extract the similarity
				hitId := hit.GetHitId()
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
					bestHitDesc = hitSeq.Desc
					bestHitLen = hitSeq.Length()
					bestHitSim = ndl.Rst.GetSimilarityPct()
				}

				chkhit++
				if chkhit >= fa.NBestHits {
					break HITS
				}
			}

			// Validate the best Hit
			bestHitLenRatio := getMinLengthRatio(bestHitLen, fa.Queries[qi].Length())
			if bestHitSim >= fa.MinSimHigh && bestHitLenRatio >= fa.MinLraHigh {
				bestHitStatus = 2
			} else if bestHitSim >= fa.MinSimNorm && bestHitLenRatio >= fa.MinLraNorm {
				bestHitStatus = 1
			}

			// Get the annotation if the best hit is good enough
			hitIsQuery := false
			if fa.DBs[fa.DBi].Equal && bestHitSim == 100.0 {
				hitIsQuery = true
			}
			if bestHitStatus > 0 {
				// Do not investigate this gene again
				fa.Finished[qi] = true
				fa.Results[qi] = *ParseHitDesc(bestHitDesc, bestHitId, fa.DBs[fa.DBi].Id, bestHitStatus, hitIsQuery)
				fa.Results[qi].HitSim = bestHitSim
				fa.Results[qi].HitLR = bestHitLenRatio
			}

		}
		// Else do nothing

		// Remove the query
		blt.ResetQuery()
	}

	// Terminate the thread
	threadChan <- 1
}
