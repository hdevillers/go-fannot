package swiss

import (
	"fmt"
)

type Entry struct {
	Access   string
	Name     string
	Locus    string
	Desc     string
	Organism string
	Phylum   string
	Sequence string
	Evidence int
}

func (e *Entry) Info() {
	fmt.Printf("Accession:\t%s\n", e.Access)
	fmt.Printf("Gene name:\t%s\n", e.Name)
	fmt.Printf("Locus tag:\t%s\n", e.Locus)
	fmt.Printf("Description:\t%s\n", e.Desc)
	fmt.Printf("Organism:\t%s\n", e.Organism)
	fmt.Printf("Phylum: \t%s\n", e.Phylum)
	fmt.Printf("Evidence:\t%d\n", e.Evidence)
	fmt.Printf("Sequence:\t%s\n", e.Sequence)
}
