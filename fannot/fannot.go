package fannot

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hdevillers/go-blast"
	"github.com/hdevillers/go-fannot/ips"
	"github.com/hdevillers/go-fannot/refdb"
	"github.com/hdevillers/go-needle"
	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/utils"
)

type Fannot struct {
	Queries        []seq.Seq
	NQueries       int
	DBs            []refdb.Refdb
	DBi            int
	DBEntries      map[string]seq.Seq
	Finished       []bool
	Results        []Result
	Param          Param
	BlastPar       blast.Param
	NeedlePar      needle.Param
	Ips            ips.Ips
	NoteFormat     Format
	ProductFormat  Format
	GeneNameFormat Format
	FunctionFormat Format
}

func NewFannot(i string) *Fannot {
	var fa Fannot

	// Load the query sequences
	fa.NQueries = utils.LoadSeqInArray(i, "fasta", &fa.Queries)

	// Init. BLAST and NEEDLE parameter settings
	fa.BlastPar = *blast.NewParam()
	fa.NeedlePar = *needle.NewParam()

	// Init. the IPS object
	fa.Ips = *ips.NewIps()

	// Init. results and Finished variables
	fa.Finished = make([]bool, fa.NQueries)
	fa.Results = make([]Result, fa.NQueries)

	// Setup default threshold
	fa.Param = *NewParam()

	return &fa
}

/*
Get reference database from input arguments
i is the list of DB ids (coma sep)
d is the directory path that contain DBs
*/
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

// Load all the sequences from the current DB
func (fa *Fannot) LoadDBEntries() {
	// Load DB entries (FASTA)
	utils.LoadSeqInMap(fa.DBs[fa.DBi].Fasta, "fasta", &fa.DBEntries)
}

// Select the next DB and load its data
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

		// Init. the Result object for the current query if not finished yet
		if !fa.Finished[qi] {
			fa.Results[qi] = Result{
				Note:     fa.Param.DefaultNote,
				Product:  fa.Param.DefaultProduct,
				GeneName: fa.Param.DefaultGeneName,
				Function: fa.Param.DefaultFunction,
			}
		}

		// Init. a BestHit object
		bestHit := NewBestHit(fa.Queries[qi], fa.NeedlePar)

		// Parse blast result
		results := blt.Rst.Iterations[0]
		if len(results.Hits) > 0 {

			// Scan each hit
		HITS:
			for _, hit := range results.Hits {
				// Retrieved the hit seq.Seq object
				hitId := hit.GetHitId()
				hitSeq, test := fa.DBEntries[hitId]
				if !test {
					panic(fmt.Sprintf("Failed to find the hit %s in the reference DB (%s).", hitId, fa.DBs[fa.DBi].Id))
				}

				// Check if it is the best hit
				bestHit.CheckHit(hitSeq)

				if bestHit.NumHits >= fa.Param.NbHitCheck {
					break HITS
				}
			}

			// Check if the best hit satisfy a rule
		CHECK:
			for ri, rule := range fa.Param.Rules {
				if rule.Test(bestHit.Similarity, bestHit.LengthRatio) {
					bestHit.IdRule = ri
					break CHECK
				}
			}

			// If found a satified rule, extract annotation
			if bestHit.IdRule != -1 {
				// Extract and copy annotation?
				copy := false
				ow := false
				if !fa.Finished[qi] {
					fa.Finished[qi] = true
					copy = true
				} else if fa.DBs[fa.DBi].OverWrite && fa.Param.Rules[bestHit.IdRule].Ovr_wrt {
					// The current DB allows overwrite
					// An overwrite is possible only if the stored
					// annotation has a hit status lower than the
					// hit status of the retained rule
					if fa.Results[qi].Status < fa.Param.Rules[bestHit.IdRule].Hit_sta {
						copy = true
						ow = true
					} else if fa.Results[qi].Status == fa.Param.Rules[bestHit.IdRule].Hit_sta {
						// In that case, similarity level must be really improved
						// Compute similarity difference
						simdiff := bestHit.Similarity - fa.Results[qi].Similarity
						if simdiff >= fa.Param.MinSimDiffOW {
							copy = true
							ow = true
						}
					}
				}

				// If annotation copy is required
				if copy {
					// Prepare description
					desc := NewDescription("uniprot", bestHit.HitSeq)

					// (re)set description according to DB and rule
					desc.Unreviewed = !fa.DBs[fa.DBi].Reviewed
					desc.SetField("Prefix", fa.Param.Rules[bestHit.IdRule].Pre_ann)
					if fa.DBs[fa.DBi].Equal && bestHit.Similarity == 100.0 {
						desc.SetField("Prefix", "")
					}

					// Complete Result
					fa.Results[qi].Note = fa.NoteFormat.Compile(desc)
					fa.Results[qi].Product = fa.ProductFormat.Compile(desc)
					fa.Results[qi].GeneName = fa.GeneNameFormat.Compile(desc)
					fa.Results[qi].Function = fa.FunctionFormat.Compile(desc)
					fa.Results[qi].Similarity = bestHit.Similarity
					fa.Results[qi].LengthRatio = bestHit.LengthRatio
					fa.Results[qi].HitOvrWrt = ow
					fa.Results[qi].HitNum = bestHit.NumHits
					fa.Results[qi].HitId = bestHit.HitSeq.Id
					fa.Results[qi].HitLocus = desc.GetField("LocusTag")
					fa.Results[qi].HitSpecies = desc.GetField("Species")
					fa.Results[qi].Status = fa.Param.Rules[bestHit.IdRule].Hit_sta
					fa.Results[qi].DbId = fa.DBs[fa.DBi].Id
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
			// Sort IPS keys in order to avoid random IPS order
			ipsids := make([]string, 0)
			for ipsid := range ips.KeyValue {
				ipsids = append(ipsids, ipsid)
			}
			sort.Strings(ipsids)

			for _, ipsid := range ipsids {
				fa.Results[qi].IpsId = append(fa.Results[qi].IpsId, ipsid)
				fa.Results[qi].IpsAnnot = append(fa.Results[qi].IpsAnnot, ips.KeyValue[ipsid])
			}

			// If no homology found, then add IpsAnnot to /note qualifier
			if fa.Results[qi].Status == 0 {
				fa.Results[qi].Note += ", InterProScan predictions: " + strings.Join(fa.Results[qi].IpsAnnot, "; ")
			}
		}
	}
}

// Write out
func (fa *Fannot) WriteOut(o string) {
	if o == "" {
		// No output path provided, print to stdout
		fmt.Print(Header())
		for i := 0; i < fa.NQueries; i++ {
			fmt.Print(fa.Results[i].ToString(fa.Queries[i].Id))
		}
	} else {
		// Save to the output path provided
		f, err := os.Create(o)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		fw := bufio.NewWriter(f)

		fw.WriteString(Header())
		for i := 0; i < fa.NQueries; i++ {
			fw.WriteString(fa.Results[i].ToString(fa.Queries[i].Id))
		}
		fw.Flush()
	}
}
