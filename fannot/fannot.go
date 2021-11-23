package fannot

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/hdevillers/go-blast"
	"github.com/hdevillers/go-fannot/ips"
	"github.com/hdevillers/go-fannot/refdb"
	"github.com/hdevillers/go-needle"
	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/utils"
)

// Initialize default thresholds
const (
	N_BEST_HITS  int     = 3
	MIN_LRA_HIGH float64 = 0.8
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
	HitNum   int
	HitOW    bool
	IpsId    []string
	IpsAnnot []string
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
		0,
		false,
		make([]string, 0),
		make([]string, 0),
	}
}

func ParseHitDesc(hd string, hid string, rid string, hs int, eq bool) *FAResult {
	var far FAResult
	values := strings.Split(hd, "::")

	far.Product = values[0]
	far.Status = hs
	far.GeneID = hid
	far.RefID = rid

	// Clean up Product
	if regexp.MustCompile(`^[A-Z][a-z ]`).MatchString(far.Product) {
		tmp := []rune(far.Product)
		tmp[0] = unicode.ToLower(tmp[0])
		far.Product = string(tmp)
	}

	// Keep only species name (delete strain data)
	tmpOrg := strings.Split(values[3], " (")
	far.Organism = tmpOrg[0]

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
		// Clean up product: remove locus tag references
		far.Product = regexp.MustCompile(" "+far.Locus).ReplaceAllString(far.Product, "")
	}
	if values[1] != "" {
		far.Note += " " + values[1]
		far.Name = values[1]
		reName := regexp.MustCompile(" " + far.Name)
		if reName.MatchString(far.Product) {
			protName := strings.Title(strings.ToLower(far.Name)) + "p"
			far.Product = reName.ReplaceAllString(far.Product, protName)
		}
	}

	if values[4] != "" {
		far.Note += ", " + values[4]
	} else {
		far.Note += ", " + far.Product
	}

	return &far
}

func (far *FAResult) PrintFAResult(gid string) {
	fmt.Printf(
		"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%.03f\t%.03f\t%s\t%d\t%t\n",
		gid, far.Product, far.Note, far.Organism,
		far.GeneID, far.Locus, far.Name,
		strings.Join(far.IpsId, ","), strings.Join(far.IpsAnnot, "; "), far.Status,
		far.HitSim, far.HitLR, far.RefID, far.HitNum,
		far.HitOW,
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
	fmt.Println("GeneID\tProduct\tNote\tOrganism\tRefID\tRefLocus\tRefName\tIPSID\tIPSAnnot\tStatus\tSimilarity\tLengthRatio\tDBID\tHitNum\tOverWritten")
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
	Ips        ips.Ips
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

	// Init. the IPS object
	fa.Ips = *ips.NewIps()

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
			bestHitNum := 0
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
					panic(fmt.Sprintf("Failed to align query %s againt ref %s, error: %s.", fa.Queries[qi].Id, hitId, err.Error()))
				}

				if ndl.Rst.GetSimilarityPct() > bestHitSim {
					bestHitId = hitId
					bestHitDesc = hitSeq.Desc
					bestHitLen = hitSeq.Length()
					bestHitSim = ndl.Rst.GetSimilarityPct()
					bestHitNum = chkhit + 1
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
				// If no annotation yet
				if !fa.Finished[qi] {
					// Set an annotation to this protein
					fa.Finished[qi] = true
					fa.Results[qi] = *ParseHitDesc(bestHitDesc, bestHitId, fa.DBs[fa.DBi].Id, bestHitStatus, hitIsQuery)
					fa.Results[qi].HitSim = bestHitSim
					fa.Results[qi].HitLR = bestHitLenRatio
					fa.Results[qi].HitNum = bestHitNum
				} else if fa.DBs[fa.DBi].OverWrite {
					// The current DB allows overwrite
					// An overwrite is possible only if the stored
					// annotation is "similar" and the new hit is
					// better.
					if fa.Results[qi].Status == 1 {
						if bestHitSim > fa.Results[qi].HitSim && bestHitLenRatio > fa.Results[qi].HitLR {
							fa.Results[qi] = *ParseHitDesc(bestHitDesc, bestHitId, fa.DBs[fa.DBi].Id, bestHitStatus, hitIsQuery)
							fa.Results[qi].HitSim = bestHitSim
							fa.Results[qi].HitLR = bestHitLenRatio
							fa.Results[qi].HitNum = bestHitNum
							fa.Results[qi].HitOW = true
						}
					}
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

func (fa *Fannot) AddIpsAnnot() {
	// Check each results
	for qi := 0; qi < fa.NQueries; qi++ {
		// Get gene ID
		gid := fa.Queries[qi].Id

		// Get ips entries for this gene (if exists)
		ips, ok := fa.Ips.Data[gid]
		if ok {
			for ipr, ann := range ips.KeyValue {
				fa.Results[qi].IpsId = append(fa.Results[qi].IpsId, ipr)
				fa.Results[qi].IpsAnnot = append(fa.Results[qi].IpsAnnot, ann)
			}

			// If no homology found, then add IpsAnnot to /note qualifier
			if fa.Results[qi].Status == 0 {
				fa.Results[qi].Note += ", InterProScan predictions: " + strings.Join(fa.Results[qi].IpsAnnot, "; ")
			}
		}
	}
}
