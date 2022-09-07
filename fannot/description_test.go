package fannot

import (
	"testing"

	"github.com/hdevillers/go-seq/seq"
)

// Test description parsing from a refdb entry
func TestDescriptionFromRefdbEntry(t *testing.T) {
	// Generate an seq object
	s := seq.Seq{
		Id:       "P25641",
		Desc:     "Putative lipase ATG15::ATG15::YCR068W::Saccharomyces cerevisiae (strain ATCC 204508 / S288c) (Baker's yeast).::lipase which is essential for lysis of subvacuolar cytoplasm to vacuole targeted bodies and intravacuolar autophagic bodies",
		Sequence: []byte("MLHKSPSRKRFASPLHLGCILTLTVLCLIAYYFALPDYLSVGKSSSRGAMDQKSDGTFRL"),
	}

	d := NewDescription("UNIPROT", s)

	// List all possible fields
	fields := []string{
		"DbName", "DbId", "ShortDesc",
		"GeneName", "ProteinName", "LocusTag",
		"Species", "LongDesc",
	}

	// Test if all fields have been initialized
	for _, field := range fields {
		if !d.HasField(field) {
			t.Fatalf("Description data should contain the field %s but it was not initialized.", field)
		}
	}

	// With the present sequence, all fields should be filled
	for _, field := range fields {
		if !d.IsSetField(field) {
			t.Fatalf("Field %s should be set.", field)
		}
	}

	// Check Db name
	if d.Data["DbName"] != "UNIPROT" {
		t.Fatalf("The field DbName should be UNIPROT, found %s", d.Data["DbName"])
	}

	// Check Id
	if d.Data["DbId"] != s.Id {
		t.Fatalf("The field DbId should be %s, found %s.", s.Id, d.Data["DbId"])
	}

	// Check Gene name
	if d.Data["GeneName"] != "ATG15" {
		t.Fatalf("The field GeneName should be ATG15, found %s", d.Data["GeneName"])
	}

	// Check Protein name
	if d.Data["ProteinName"] != "Atg15p" {
		t.Fatalf("The field GeneName should be Atg15p, found %s", d.Data["ProteinName"])
	}

	// Check Species name
	if d.Data["Species"] != "Saccharomyces cerevisiae" {
		t.Fatalf("The field Species should be Saccharomyces cerevisiae, found %s", d.Data["Species"])
	}
}

// Test description parsing from a uncompleted refdb entry
func TestDescriptionFromPartialRefdbEntry(t *testing.T) {
	// Generate an seq object
	s := seq.Seq{
		Id:       "P25641",
		Desc:     "Putative lipase ATG15::::::Saccharomyces cerevisiae (strain ATCC 204508 / S288c) (Baker's yeast).::",
		Sequence: []byte("MLHKSPSRKRFASPLHLGCILTLTVLCLIAYYFALPDYLSVGKSSSRGAMDQKSDGTFRL"),
	}

	d := NewDescription("UNIPROT", s)

	// In that case, GeneName and LocusTag must have the default value
	if d.Data["GeneName"] != DEFAULT_DESCRIPTION {
		t.Fatalf("The field GeneName should not be set, but found %s.", d.Data["GeneName"])
	}
	if d.Data["LocusTag"] != DEFAULT_DESCRIPTION {
		t.Fatalf("The field LocusTag should not be set, but found %s.", d.Data["LocusTag"])
	}

	// Int that case, ShortDesc should be equal to LongDesc
	if d.Data["ShortDesc"] != d.Data["LongDesc"] {
		t.Fatalf("The fields ShortDesc and LongDesc should be equal, found the short: %s; and the long: %s.", d.Data["ShortDesc"], d.Data["LongDesc"])
	}
}

// Test description parsing from a simple fasta entry
func TestDescriptionFromSimpleEntry(t *testing.T) {
	// Generate an seq object
	s := seq.Seq{
		Id:       "P25641",
		Desc:     "Putative lipase ATG15",
		Sequence: []byte("MLHKSPSRKRFASPLHLGCILTLTVLCLIAYYFALPDYLSVGKSSSRGAMDQKSDGTFRL"),
	}

	d := NewDescription("UNIPROT", s)

	// In that case, GeneName and LocusTag must have the default value
	if d.Data["GeneName"] != DEFAULT_DESCRIPTION {
		t.Fatalf("The field GeneName should not be set, but found %s.", d.Data["GeneName"])
	}
	if d.Data["LocusTag"] != DEFAULT_DESCRIPTION {
		t.Fatalf("The field LocusTag should not be set, but found %s.", d.Data["LocusTag"])
	}

	// Int that case, ShortDesc and LongDesc should be equal to sequence description
	if d.Data["ShortDesc"] != s.Desc {
		t.Fatalf("The fields ShortDesc should be equal to %s, found %s.", s.Desc, d.Data["ShortDesc"])
	}
	if d.Data["LongDesc"] != s.Desc {
		t.Fatalf("The fields LongDesc should be equal to %s, found %s.", s.Desc, d.Data["LongDesc"])
	}
}

func TestDescriptionCorruptedEntry(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal(`Description parsing should have failed, but it did not.`)
		}
	}()

	// Generate an seq object
	s := seq.Seq{
		Id:       "P25641",
		Desc:     "Putative lipase ATG15::ATG15",
		Sequence: []byte("MLHKSPSRKRFASPLHLGCILTLTVLCLIAYYFALPDYLSVGKSSSRGAMDQKSDGTFRL"),
	}

	NewDescription("UNIPROT", s)
}
