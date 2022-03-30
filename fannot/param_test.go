package fannot

import (
	"testing"
)

// Test the default parameter settings
func TestParamDefault(t *testing.T) {
	p := NewParam()

	// Check the default best hit number
	if p.Nbh_chk != N_BEST_HITS {
		t.Errorf("Default number of hits should be %d, found %d", N_BEST_HITS, p.Nbh_chk)
	}

	// Check the default annotation for unknown function
	if p.Unk_ann != UNKNOWN_FUNC {
		t.Errorf("Default unknown function annotation should be %s, found %s", UNKNOWN_FUNC, p.Unk_ann)
	}

	// Check the number of rules (=2)
	if len(p.Rules) != 2 {
		t.Errorf("Expecting 2 default rules, found %d", len(p.Rules))
	}
}

// Test to read the three rules JSON (in example)
func TestParamThreeRules(t *testing.T) {
	p := NewParamFromJson("../examples/three_levels.json")

	// Check the number of rules
	if len(p.Rules) != 3 {
		t.Errorf("Expecting 3 default rules, found %d", len(p.Rules))
	}
}
