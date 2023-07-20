package fannot

import (
	"regexp"
	"testing"
)

// Test overwrite a weakly with a highly
func TestRefdbOverwrite01(t *testing.T) {
	// Initiate the fannot object
	fa := NewFannot("../examples/queries/query_01.fasta")

	// query_01.fasta contains 1 sequence
	if fa.NQueries != 1 {
		t.Fatalf(`Failed to initiate fannot with query_01.fasta, expected 1 sequence, found %d.`, fa.NQueries)
	}

	// Load the reference DB
	fa.GetDBs("TEST_WEAKLY_NW,TEST_HIGHLY_OW", "../examples/refdb/")

	// Load the JSON file
	fa.Param = *NewParamFromJson("../examples/three_levels.json")

	// Init. annotation templates/format
	fa.NoteFormat = *NewFormat(fa.Param.TemplateNote)
	fa.ProductFormat = *NewFormat(fa.Param.TemplateProduct)
	fa.GeneNameFormat = *NewFormat(fa.Param.TemplateGeneName)
	fa.FunctionFormat = *NewFormat(fa.Param.TemplateFunction)

REFDB:
	for fa.NextDB() {
		// fannot has been developed to use channels and multi-threading
		// Initiate channels
		queryChan := make(chan int)
		threadChan := make(chan int)

		// Launch a routine (mono-thread)
		go fa.FindFunction(queryChan, threadChan)

		// Push queries in channel
		nq := 0
		for i := 0; i < fa.NQueries; i++ {
			if !fa.Finished[i] {
				nq++
				queryChan <- i
			} else if fa.DBs[fa.DBi].OverWrite && fa.Results[i].Status <= fa.Param.MaxStatusOW {
				// Try to overwrite the annotation
				nq++
				queryChan <- i
			}
		}
		close(queryChan)

		// Wait end of computation
		<-threadChan

		if nq == 0 {
			break REFDB
		}
	}

	// Must find a function
	if !fa.Finished[0] {
		t.Fatal(`No function found for sequence 1 while it should find one.`)
	}

	// The final match must be highly
	re := regexp.MustCompile(`^highly`)
	if !re.MatchString(fa.Results[0].Note) {
		t.Fatalf(`First sequence note should start by "highly", obtained: %s`, fa.Results[0].Note)
	}

	if !fa.Results[0].HitOvrWrt {
		t.Fatalf(`Hit should have been overwritten, but it is not the case.`)
	}
}

// REFDB prevent from overwrite
func TestRefdbOverwrite02(t *testing.T) {
	// Initiate the fannot object
	fa := NewFannot("../examples/queries/query_01.fasta")

	// query_01.fasta contains 1 sequence
	if fa.NQueries != 1 {
		t.Fatalf(`Failed to initiate fannot with query_01.fasta, expected 1 sequence, found %d.`, fa.NQueries)
	}

	// Load the reference DB
	fa.GetDBs("TEST_WEAKLY_NW,TEST_HIGHLY_NW", "../examples/refdb/")

	// Load the JSON file
	fa.Param = *NewParamFromJson("../examples/three_levels.json")

	// Init. annotation templates/format
	fa.NoteFormat = *NewFormat(fa.Param.TemplateNote)
	fa.ProductFormat = *NewFormat(fa.Param.TemplateProduct)
	fa.GeneNameFormat = *NewFormat(fa.Param.TemplateGeneName)
	fa.FunctionFormat = *NewFormat(fa.Param.TemplateFunction)

	nrun := 0
REFDB:
	for fa.NextDB() {
		// fannot has been developed to use channels and multi-threading
		// Initiate channels
		queryChan := make(chan int)
		threadChan := make(chan int)

		// Launch a routine (mono-thread)
		go fa.FindFunction(queryChan, threadChan)

		// Push queries in channel
		nq := 0
		for i := 0; i < fa.NQueries; i++ {
			if !fa.Finished[i] {
				nq++
				queryChan <- i
				nrun++
			} else if fa.DBs[fa.DBi].OverWrite && fa.Results[i].Status <= fa.Param.MaxStatusOW {
				// Try to overwrite the annotation
				nq++
				queryChan <- i
				nrun++
			}
		}
		close(queryChan)

		// Wait end of computation
		<-threadChan

		if nq == 0 {
			break REFDB
		}
	}

	// Must find a function
	if !fa.Finished[0] {
		t.Fatal(`No function found for sequence 1 while it should find one.`)
	}

	// The final match must be highly
	re := regexp.MustCompile(`^weakly`)
	if !re.MatchString(fa.Results[0].Note) {
		t.Fatalf(`First sequence note should start by "weakly", obtained: %s`, fa.Results[0].Note)
	}

	if fa.Results[0].HitOvrWrt {
		t.Fatalf(`Hit should not be overwritten, but it is the case.`)
	}

	// The query should have been investigated only once
	if nrun != 1 {
		t.Fatalf(`The query should have been investigated only once, scanned %d times.`, nrun)
	}
}

// JSON (max status) prevents from overwriting
func TestRefdbOverwrite03(t *testing.T) {
	// Initiate the fannot object
	fa := NewFannot("../examples/queries/query_01.fasta")

	// query_01.fasta contains 1 sequence
	if fa.NQueries != 1 {
		t.Fatalf(`Failed to initiate fannot with query_01.fasta, expected 1 sequence, found %d.`, fa.NQueries)
	}

	// Load the reference DB
	fa.GetDBs("TEST_SIMILAR_NW,TEST_HIGHLY_OW", "../examples/refdb/")

	// Load the JSON file
	fa.Param = *NewParamFromJson("../examples/three_levels.json")

	// Init. annotation templates/format
	fa.NoteFormat = *NewFormat(fa.Param.TemplateNote)
	fa.ProductFormat = *NewFormat(fa.Param.TemplateProduct)
	fa.GeneNameFormat = *NewFormat(fa.Param.TemplateGeneName)
	fa.FunctionFormat = *NewFormat(fa.Param.TemplateFunction)

REFDB:
	for fa.NextDB() {
		// fannot has been developed to use channels and multi-threading
		// Initiate channels
		queryChan := make(chan int)
		threadChan := make(chan int)

		// Launch a routine (mono-thread)
		go fa.FindFunction(queryChan, threadChan)

		// Push queries in channel
		nq := 0
		for i := 0; i < fa.NQueries; i++ {
			if !fa.Finished[i] {
				nq++
				queryChan <- i
			} else if fa.DBs[fa.DBi].OverWrite && fa.Results[i].Status <= fa.Param.MaxStatusOW {
				// Try to overwrite the annotation
				nq++
				queryChan <- i
			}
		}
		close(queryChan)

		// Wait end of computation
		<-threadChan

		if nq == 0 {
			break REFDB
		}
	}

	// Must find a function
	if !fa.Finished[0] {
		t.Fatal(`No function found for sequence 1 while it should find one.`)
	}

	// The final match must be similar
	re := regexp.MustCompile(`^similar`)
	if !re.MatchString(fa.Results[0].Note) {
		t.Fatalf(`First sequence note should start by "similar", obtained: %s`, fa.Results[0].Note)
	}

	if fa.Results[0].HitOvrWrt {
		t.Fatalf(`Hit should not be overwritten, but it is the case.`)
	}
}

// JSON (max status) allows overwriting
func TestRefdbOverwrite04(t *testing.T) {
	// Initiate the fannot object
	fa := NewFannot("../examples/queries/query_01.fasta")

	// query_01.fasta contains 1 sequence
	if fa.NQueries != 1 {
		t.Fatalf(`Failed to initiate fannot with query_01.fasta, expected 1 sequence, found %d.`, fa.NQueries)
	}

	// Load the reference DB
	fa.GetDBs("TEST_SIMILAR_NW,TEST_HIGHLY_OW", "../examples/refdb/")

	// Load the JSON file
	fa.Param = *NewParamFromJson("../examples/other_maxStatusOW.json")

	// Init. annotation templates/format
	fa.NoteFormat = *NewFormat(fa.Param.TemplateNote)
	fa.ProductFormat = *NewFormat(fa.Param.TemplateProduct)
	fa.GeneNameFormat = *NewFormat(fa.Param.TemplateGeneName)
	fa.FunctionFormat = *NewFormat(fa.Param.TemplateFunction)

REFDB:
	for fa.NextDB() {
		// fannot has been developed to use channels and multi-threading
		// Initiate channels
		queryChan := make(chan int)
		threadChan := make(chan int)

		// Launch a routine (mono-thread)
		go fa.FindFunction(queryChan, threadChan)

		// Push queries in channel
		nq := 0
		for i := 0; i < fa.NQueries; i++ {
			if !fa.Finished[i] {
				nq++
				queryChan <- i
			} else if fa.DBs[fa.DBi].OverWrite && fa.Results[i].Status <= fa.Param.MaxStatusOW {
				// Try to overwrite the annotation
				nq++
				queryChan <- i
			}
		}
		close(queryChan)

		// Wait end of computation
		<-threadChan

		if nq == 0 {
			break REFDB
		}
	}

	// Must find a function
	if !fa.Finished[0] {
		t.Fatal(`No function found for sequence 1 while it should find one.`)
	}

	// The final match must be similar
	re := regexp.MustCompile(`^highly`)
	if !re.MatchString(fa.Results[0].Note) {
		t.Fatalf(`First sequence note should start by "highly", obtained: %s`, fa.Results[0].Note)
	}

	if !fa.Results[0].HitOvrWrt {
		t.Fatalf(`Hit should have been overwritten, but it is not the case.`)
	}
}
