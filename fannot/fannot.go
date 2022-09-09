package fannot

import (
	"strings"

	"github.com/hdevillers/go-blast"
	"github.com/hdevillers/go-fannot/ips"
	"github.com/hdevillers/go-fannot/refdb"
	"github.com/hdevillers/go-needle"
	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/utils"
)

type Fannot struct {
	Queries   []seq.Seq
	NQueries  int
	DBs       []refdb.Refdb
	DBi       int
	DBEntries map[string]seq.Seq
	Finished  []bool
	Results   []Result
	Param     Param
	BlastPar  blast.Param
	NeedlePar needle.Param
	Ips       ips.Ips
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
	fa.Results = make([]Result, fa.NQueries)

	/* INIT. Results later */
	//for i := 0; i < fa.NQueries; i++ {
	// Set default result (when no match is found)
	//	fa.Results[i] = *NewResult()
	//}

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

		// Init. the Result object for the current query
		fa.Results[i] = Result{}
}