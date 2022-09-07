package fannot

import (
	"testing"

	"github.com/hdevillers/go-seq/seq"
)

func TestFormatEmpty(t *testing.T) {
	// Generate an seq object
	s := seq.Seq{
		Id:       "P25641",
		Desc:     "lipase ATG15::ATG15::YCR068W::Saccharomyces cerevisiae (strain ATCC 204508 / S288c) (Baker's yeast).::lipase which is essential for lysis of subvacuolar cytoplasm to vacuole targeted bodies and intravacuolar autophagic bodies",
		Sequence: []byte("MLHKSPSRKRFASPLHLGCILTLTVLCLIAYYFALPDYLSVGKSSSRGAMDQKSDGTFRL"),
	}

	// Create the description
	d := NewDescription("UNIPROT", s)

	// Init. a new format object
	f := NewFormat("")

	if !f.Empty {
		t.Error("The value of the attribute 'Empty' should be 'false'")
	}

	// Compile the template
	o := f.Compile(d)

	if o != "" {
		t.Error("When 'Empty' is 'true', Complile should return an empty string.")
	}

}

func TestFormatGeneDescription(t *testing.T) {
	// Generate an seq object
	s := seq.Seq{
		Id:       "P25641",
		Desc:     "lipase ATG15::ATG15::YCR068W::Saccharomyces cerevisiae (strain ATCC 204508 / S288c) (Baker's yeast).::lipase which is essential for lysis of subvacuolar cytoplasm to vacuole targeted bodies and intravacuolar autophagic bodies",
		Sequence: []byte("MLHKSPSRKRFASPLHLGCILTLTVLCLIAYYFALPDYLSVGKSSSRGAMDQKSDGTFRL"),
	}

	// Create the description
	d := NewDescription("UNIPROT", s)

	// Init. an input template
	i := "{Prefix} ||{DbName}|{DbId} ||{Species}"

	// Init. a new format object
	f := NewFormat(i)

	// Compile the template
	o := f.Compile(d)

	// The expected out is "UNIPROT|P25641 Saccharomyces cerevisiae"
	e := "UNIPROT|P25641 Saccharomyces cerevisiae"
	if o != e {
		t.Errorf("Failed to format the gene description, expected: %s; obtained: %s", e, o)
	}
}

func TestFormatTransGeneName(t *testing.T) {
	// Generate an seq object
	s := seq.Seq{
		Id:       "P25641",
		Desc:     "lipase ATG15::ATG15::YCR068W::Saccharomyces cerevisiae (strain ATCC 204508 / S288c) (Baker's yeast).::lipase which is essential for lysis of subvacuolar cytoplasm to vacuole targeted bodies and intravacuolar autophagic bodies",
		Sequence: []byte("MLHKSPSRKRFASPLHLGCILTLTVLCLIAYYFALPDYLSVGKSSSRGAMDQKSDGTFRL"),
	}

	// Create the description
	d := NewDescription("UNIPROT", s)

	// Init. an input template
	i := "putative {ShortDesc}::GnPn"

	// Init. a new format object
	f := NewFormat(i)

	// Compile the template
	o := f.Compile(d)

	// The expected out is "UNIPROT|P25641 Saccharomyces cerevisiae"
	e := "putative lipase Atg15p"
	if o != e {
		t.Errorf("Failed to format the gene description, expected: %s; obtained: %s", e, o)
	}
}

func TestFormatBadField(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal(`Failed to detect bad field.`)
		}
	}()

	// Input pattern with a bad field
	i := "putative {BadField}"

	// Init. a new format object
	NewFormat(i)
}
