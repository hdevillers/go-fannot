package fannot

const (
	MIN_LRA_HIGH float64 = 0.8
	MIN_SIM_HIGH float64 = 80.0
	MIN_LRA_NORM float64 = 0.7
	MIN_SIM_NORM float64 = 50.0
	PRE_SIM_HIGH string  = "highly similar to "
	PRE_SIM_NORM string  = "similar to "
	CPY_GEN_HIGH bool    = true
	CPY_GEN_NORM bool    = false
	OVR_WRT_HIGH bool    = true
	OVR_WRT_NORM bool    = false
	HIT_STA_HIGH int     = 2
	HIT_STA_NORM int     = 1
)

// Single rule object
type Rule struct {
	Min_sim float64 // Minimal similarity threshold
	Min_lra float64 // Minimal length ratio threshold
	Pre_ann string  // Annotation prefix
	Cpy_gen bool    // Copy the gene name in the annotation
	Ovr_wrt bool    // Can overwrite a previous annotation
	Hit_sta int     // Hit status (integer)
}

// Default rule builder: Highly similar
func NewRuleHighlySimilar() *Rule {
	return &Rule{
		Min_sim: MIN_SIM_HIGH,
		Min_lra: MIN_LRA_HIGH,
		Pre_ann: PRE_SIM_HIGH,
		Cpy_gen: CPY_GEN_HIGH,
		Ovr_wrt: OVR_WRT_HIGH,
		Hit_sta: HIT_STA_HIGH,
	}
}

// Default rule builder: Highly similar
func NewRuleSimilar() *Rule {
	return &Rule{
		Min_sim: MIN_SIM_NORM,
		Min_lra: MIN_LRA_NORM,
		Pre_ann: PRE_SIM_NORM,
		Cpy_gen: CPY_GEN_NORM,
		Ovr_wrt: OVR_WRT_NORM,
		Hit_sta: HIT_STA_NORM,
	}
}

// Test if a query respect the rule
func (r *Rule) Test(s float64, l float64) bool {
	if s >= r.Min_sim && l >= r.Min_lra {
		return true
	}
	return false
}
