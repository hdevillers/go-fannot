package fannot

import (
	"testing"
)

// Functional test: input01.fasta, one thread
func TestHitNumber(t *testing.T) {
	// Initiate the fannot object
	fa := NewFannot("../examples/queries/query_hitnumber.fasta")

	// Input01 contains 4 sequences
	if fa.NQueries != 3 {
		t.Fatalf(`Failed to initiate fannot with query_hitnumber.fasta, expected 3 sequences, found %d.`, fa.NQueries)
	}

	// Load the reference DB
	fa.GetDBs("TEST_HIT_NUMBER", "../examples/refdb/")

	// There is only on DB provided, select it
	if !fa.NextDB() {
		t.Fatal(`Failed to load reference database.`)
	}

	// This database is name TEST_HIT_NUMBER
	if fa.DBs[fa.DBi].Id != "TEST_HIT_NUMBER" {
		t.Fatalf(`Bad reference DB id, expected TEST_HIT_NUMBER, obtained %s.`, fa.DBs[fa.DBi].Id)
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

	if !fa.Finished[1] {
		t.Fatal(`No function found for sequence 2 while it should find one.`)
	}

	if !fa.Finished[2] {
		t.Fatal(`No function found for sequence 2 while it should find one.`)
	}

	// The first query should have a hit number of 1
	if fa.Results[0].HitNum != 1 {
		t.Fatalf(`First query should have a hit number of 1, obtained: %d.`, fa.Results[0].HitNum)
	}

	// The second query should have a hit number of 2
	if fa.Results[1].HitNum != 2 {
		t.Fatalf(`Second query should have a hit number of 2, obtained: %d.`, fa.Results[1].HitNum)
	}

	// The third query should have a hit number of 3
	if fa.Results[2].HitNum != 3 {
		t.Fatalf(`Third query should have a hit number of 3, obtained: %d.`, fa.Results[2].HitNum)
	}

}
