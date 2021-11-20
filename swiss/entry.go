package swiss

import (
	"fmt"
	"regexp"
)

type Entry struct {
	Access   string
	Name     string
	Locus    string
	Desc     string
	Length   int
	Organism string
	Phylum   string
	Sequence string
	Evidence string
}

func (e *Entry) Info() {
	fmt.Printf("Accession:\t%s\n", e.Access)
	fmt.Printf("Gene name:\t%s\n", e.Name)
	fmt.Printf("Locus tag:\t%s\n", e.Locus)
	fmt.Printf("Description:\t%s\n", e.Desc)
	fmt.Printf("Length (aa):\t%d\n", e.Length)
	fmt.Printf("Organism:\t%s\n", e.Organism)
	fmt.Printf("Phylum: \t%s\n", e.Phylum)
	fmt.Printf("Evidence:\t%s\n", e.Evidence)
	fmt.Printf("Sequence:\t%s\n", e.Sequence)
}

func (e *Entry) TestEvidence(re string) bool {
	retest, err := regexp.Compile(re)
	if err != nil {
		panic("[TestEvidence]: Cannot compile regex.")
	}

	if retest.MatchString(e.Evidence) {
		return true
	}

	return false
}

func (e *Entry) TestTaxonomy(re string) bool {
	retest, err := regexp.Compile(re)
	if err != nil {
		panic("[TestTaxonomy]: Cannot compile regex.")
	}

	// First test the Organism string
	if retest.MatchString(e.Organism) {
		return true
	}

	// Then test the phylum (if previous test is false)
	if retest.MatchString(e.Phylum) {
		return true
	}

	return false
}

func (e *Entry) Test(tk, ts, ek, es string) bool {
	// Skip direct
	return false
}
