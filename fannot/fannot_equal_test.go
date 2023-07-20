package fannot

import (
	"regexp"
	"testing"
)

// Functional test: input01.fasta, one thread
func TestRefdbHasEqual(t *testing.T) {
	// Initiate the fannot object
	fa := NewFannot("../examples/queries/query_01.fasta")

	// query_01.fasta contains 1 sequence
	if fa.NQueries != 1 {
		t.Fatalf(`Failed to initiate fannot with query_01.fasta, expected 1 sequence, found %d.`, fa.NQueries)
	}

	// Load the reference DB
	fa.GetDBs("TEST_EQUAL", "../examples/refdb/")

	// There is only on DB provided, select it
	if !fa.NextDB() {
		t.Fatal(`Failed to load reference database.`)
	}

	// This database is name S288c_lipase
	if fa.DBs[fa.DBi].Id != "TEST_EQUAL" {
		t.Fatalf(`Bad reference DB id, expected TEST_EQUAL, obtained %s.`, fa.DBs[fa.DBi].Id)
	}

	// Setup format (with default)
	fa.NoteFormat = *NewFormat(TPL_NOTE)
	fa.ProductFormat = *NewFormat(TPL_PRODUCT)
	fa.GeneNameFormat = *NewFormat(TPL_GENENAME)
	fa.FunctionFormat = *NewFormat(TPL_FUNCTION)

	// fannot has been developped to use channels and multi-threading
	// Initiate channels
	queryChan := make(chan int)
	threadChan := make(chan int)

	// Launch a routine (mono-thread)
	go fa.FindFunction(queryChan, threadChan)

	// Push queries in channel
	for i := 0; i < fa.NQueries; i++ {
		queryChan <- i
	}
	close(queryChan)

	// Wait end of computation
	<-threadChan

	// In this toy example, sequence 1 and 2 should have a function, not the two other
	if !fa.Finished[0] {
		t.Fatal(`No function found for sequence 1 while it should find one.`)
	}

	// The first sequence is contained in the refdb, hence not similarity
	// level should be indicated (just equality)
	re := regexp.MustCompile(`^uniprot`)
	if !re.MatchString(fa.Results[0].Note) {
		t.Fatalf(`First sequence note should start by "uniprot", obtained: %s`, fa.Results[0].Note)
	}

	// The first sequence has a similarity of 100%
	if fa.Results[0].Similarity != 100.0 {
		t.Fatalf(`First sequence should have a hit similarity of 100, obtained: %.02f`, fa.Results[0].Similarity)
	}
}

// Functional test: input01.fasta, one thread
func TestRefdbHasNotEqual(t *testing.T) {
	// Initiate the fannot object
	fa := NewFannot("../examples/queries/query_01.fasta")

	// query_01.fasta contains 1 sequence
	if fa.NQueries != 1 {
		t.Fatalf(`Failed to initiate fannot with query_01.fasta, expected 1 sequence, found %d.`, fa.NQueries)
	}

	// Load the reference DB
	fa.GetDBs("TEST_NOT_EQUAL", "../examples/refdb/")

	// There is only on DB provided, select it
	if !fa.NextDB() {
		t.Fatal(`Failed to load reference database.`)
	}

	// This database is name S288c_lipase
	if fa.DBs[fa.DBi].Id != "TEST_NOT_EQUAL" {
		t.Fatalf(`Bad reference DB id, expected TEST_NOT_EQUAL, obtained %s.`, fa.DBs[fa.DBi].Id)
	}

	// Setup format (with default)
	fa.NoteFormat = *NewFormat(TPL_NOTE)
	fa.ProductFormat = *NewFormat(TPL_PRODUCT)
	fa.GeneNameFormat = *NewFormat(TPL_GENENAME)
	fa.FunctionFormat = *NewFormat(TPL_FUNCTION)

	// fannot has been developped to use channels and multi-threading
	// Initiate channels
	queryChan := make(chan int)
	threadChan := make(chan int)

	// Launch a routine (mono-thread)
	go fa.FindFunction(queryChan, threadChan)

	// Push queries in channel
	for i := 0; i < fa.NQueries; i++ {
		queryChan <- i
	}
	close(queryChan)

	// Wait end of computation
	<-threadChan

	// In this toy example, sequence 1 and 2 should have a function, not the two other
	if !fa.Finished[0] {
		t.Fatal(`No function found for sequence 1 while it should find one.`)
	}

	// The first sequence is contained in the refdb, hence not similarity
	// level should be indicated (just equality)
	re := regexp.MustCompile(`^highly`)
	if !re.MatchString(fa.Results[0].Note) {
		t.Fatalf(`First sequence note should start by "highly", obtained: %s`, fa.Results[0].Note)
	}

	// The first sequence has a similarity of 100%
	if fa.Results[0].Similarity != 100.0 {
		t.Fatalf(`First sequence should have a hit similarity of 100, obtained: %.02f`, fa.Results[0].Similarity)
	}
}