package fannot

// Single rule object
type Rule struct {
	Min_sim float64 // Minimal similarity threshold
	Min_lra float64 // Minimal length ratio threshold
	Pre_ann string  // Annotation prefix
	Cpy_gen bool    // Copy the gene name in the annotation
	Ovr_wrt bool    // Can overwrite a previous annotation
	Hit_sta int     // Hit status (integer)
}

// Test if a query respect the rule
func (r *Rule) Test(s float64, l float64) bool {
	if s >= r.Min_sim && l >= r.Min_lra {
		return true
	}
	return false
}
