package fannot

import (
	"regexp"
	"testing"
)

// Test a simple highly similar match
func TestRefdbHighlySimilar(t *testing.T) {
	// Initiate the fannot object
	fa := NewFannot("../examples/queries/query_01.fasta")

	// query_01.fasta contains 1 sequence
	if fa.NQueries != 1 {
		t.Fatalf(`Failed to initiate fannot with query_01.fasta, expected 1 sequence, found %d.`, fa.NQueries)
	}

	// Load the reference DB
	fa.GetDBs("TEST_HIGHLY_NW", "../examples/refdb/")

	// There is only on DB provided, select it
	if !fa.NextDB() {
		t.Fatal(`Failed to load reference database.`)
	}

	// Check database name
	if fa.DBs[fa.DBi].Id != "TEST_HIGHLY_NW" {
		t.Fatalf(`Bad reference DB id, expected TEST_HIGHLY_NW, obtained %s.`, fa.DBs[fa.DBi].Id)
	}

	// Load the JSON file
	fa.Param = *NewParamFromJson("../examples/three_levels.json")

	// Init. annotation templates/format
	fa.NoteFormat = *NewFormat(fa.Param.TemplateNote)
	fa.ProductFormat = *NewFormat(fa.Param.TemplateProduct)
	fa.GeneNameFormat = *NewFormat(fa.Param.TemplateGeneName)
	fa.FunctionFormat = *NewFormat(fa.Param.TemplateFunction)

	// fannot has been developed to use channels and multi-threading
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

	// Must find a function
	if !fa.Finished[0] {
		t.Fatal(`No function found for sequence 1 while it should find one.`)
	}

	// Check similarity level
	re := regexp.MustCompile(`^highly`)
	if !re.MatchString(fa.Results[0].Note) {
		t.Fatalf(`First sequence note should start by "highly", obtained: %s`, fa.Results[0].Note)
	}
	if fa.Results[0].Status != 3 {
		t.Fatalf(`Hit status should be 3, obtained: %d`, fa.Results[0].Status)
	}

	// Check similarity value
	if fa.Results[0].Similarity < fa.Param.Rules[0].Min_sim {
		t.Fatalf(`First sequence similarity is too low for a highly similar hit, obtained: %.02f`, fa.Results[0].Similarity)
	}
}

// Test a simple similar match
func TestRefdbSimilar(t *testing.T) {
	// Initiate the fannot object
	fa := NewFannot("../examples/queries/query_01.fasta")

	// query_01.fasta contains 1 sequence
	if fa.NQueries != 1 {
		t.Fatalf(`Failed to initiate fannot with query_01.fasta, expected 1 sequence, found %d.`, fa.NQueries)
	}

	// Load the reference DB
	fa.GetDBs("TEST_SIMILAR_NW", "../examples/refdb/")

	// There is only on DB provided, select it
	if !fa.NextDB() {
		t.Fatal(`Failed to load reference database.`)
	}

	// Check database name
	if fa.DBs[fa.DBi].Id != "TEST_SIMILAR_NW" {
		t.Fatalf(`Bad reference DB id, expected TEST_SIMILAR_NW, obtained %s.`, fa.DBs[fa.DBi].Id)
	}

	// Load the JSON file
	fa.Param = *NewParamFromJson("../examples/three_levels.json")

	// Init. annotation templates/format
	fa.NoteFormat = *NewFormat(fa.Param.TemplateNote)
	fa.ProductFormat = *NewFormat(fa.Param.TemplateProduct)
	fa.GeneNameFormat = *NewFormat(fa.Param.TemplateGeneName)
	fa.FunctionFormat = *NewFormat(fa.Param.TemplateFunction)

	// fannot has been developed to use channels and multi-threading
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

	// Must find a function
	if !fa.Finished[0] {
		t.Fatal(`No function found for sequence 1 while it should find one.`)
	}

	// Check similarity level
	re := regexp.MustCompile(`^similar`)
	if !re.MatchString(fa.Results[0].Note) {
		t.Fatalf(`First sequence note should start by "similar", obtained: %s`, fa.Results[0].Note)
	}
	if fa.Results[0].Status != 2 {
		t.Fatalf(`Hit status should be 2, obtained: %d`, fa.Results[0].Status)
	}

	// Check similarity value
	if fa.Results[0].Similarity < fa.Param.Rules[1].Min_sim {
		t.Fatalf(`First sequence similarity is too low for a similar hit, obtained: %.02f`, fa.Results[0].Similarity)
	}
	if fa.Results[0].Similarity >= fa.Param.Rules[0].Min_sim {
		t.Fatalf(`First sequence similarity is too high for a similar hit, obtained: %.02f`, fa.Results[0].Similarity)
	}
}

// Test a simple weakly similar match
func TestRefdbWeaklySimilar(t *testing.T) {
	// Initiate the fannot object
	fa := NewFannot("../examples/queries/query_01.fasta")

	// query_01.fasta contains 1 sequence
	if fa.NQueries != 1 {
		t.Fatalf(`Failed to initiate fannot with query_01.fasta, expected 1 sequence, found %d.`, fa.NQueries)
	}

	// Load the reference DB
	fa.GetDBs("TEST_WEAKLY_NW", "../examples/refdb/")

	// There is only on DB provided, select it
	if !fa.NextDB() {
		t.Fatal(`Failed to load reference database.`)
	}

	// Check database name
	if fa.DBs[fa.DBi].Id != "TEST_WEAKLY_NW" {
		t.Fatalf(`Bad reference DB id, expected TEST_WEAKLY_NW, obtained %s.`, fa.DBs[fa.DBi].Id)
	}

	// Load the JSON file
	fa.Param = *NewParamFromJson("../examples/three_levels.json")

	// Init. annotation templates/format
	fa.NoteFormat = *NewFormat(fa.Param.TemplateNote)
	fa.ProductFormat = *NewFormat(fa.Param.TemplateProduct)
	fa.GeneNameFormat = *NewFormat(fa.Param.TemplateGeneName)
	fa.FunctionFormat = *NewFormat(fa.Param.TemplateFunction)

	// fannot has been developed to use channels and multi-threading
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

	if !fa.Finished[0] {
		t.Fatal(`No function found for sequence 1 while it should find one.`)
	}

	// Check similarity level found
	re := regexp.MustCompile(`^weakly`)
	if !re.MatchString(fa.Results[0].Note) {
		t.Fatalf(`First sequence note should start by "weakly", obtained: %s`, fa.Results[0].Note)
	}
	if fa.Results[0].Status != 1 {
		t.Fatalf(`Hit status should be 1, obtained: %d`, fa.Results[0].Status)
	}

	// Check similarity values
	if fa.Results[0].Similarity < fa.Param.Rules[2].Min_sim {
		t.Fatalf(`First sequence similarity is too low for a weakly similar hit, obtained: %.02f`, fa.Results[0].Similarity)
	}
	if fa.Results[0].Similarity >= fa.Param.Rules[1].Min_sim {
		t.Fatalf(`First sequence similarity is too high for a weakly similar hit, obtained: %.02f`, fa.Results[0].Similarity)
	}
}
