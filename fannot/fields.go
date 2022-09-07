package fannot

type Fields struct {
	Labels map[string]int
}

func NewFields() *Fields {
	var f Fields

	// List all alloxed fields
	af := []string{
		"DbName", "DbId", "ShortDesc",
		"LongDesc", "GeneName", "ProteinName",
		"LocusTag", "Species", "Prefix",
		"Putative", "Unreviewed",
	}

	// Init. Labels attribute
	f.Labels = make(map[string]int)
	for _, l := range af {
		f.Labels[l] = 1
	}

	return &f
}

func (f *Fields) Exists(s string) bool {
	_, test := f.Labels[s]
	return test
}
