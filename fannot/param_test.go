package fannot

import (
	"testing"
)

// Test the default parameter settings
func TestParamDefault(t *testing.T) {
	p := NewParam()

	// Check the default best hit number
	if p.NbHitCheck != NB_HIT_CHECK {
		t.Errorf("Default number of hits should be %d, found %d", NB_HIT_CHECK, p.NbHitCheck)
	}

	// Check the default annotation for unknown function
	if p.DefaultNote != DFT_NOTE {
		t.Errorf("Default unknown function annotation should be %s, found %s", DFT_NOTE, p.DefaultNote)
	}

	// Check the number of rules (=2)
	if len(p.Rules) != 2 {
		t.Errorf("Expecting 2 default rules, found %d", len(p.Rules))
	}

	// Check a given rule
	if p.Rules[0].Min_sim != MIN_SIM_HIGH {
		t.Errorf("Expected a minimal similarity of %.02f, found a similarity of %.02f.", MIN_SIM_HIGH, p.Rules[0].Min_sim)
	}
	if p.Rules[0].Min_lra != MIN_LRA_HIGH {
		t.Errorf("Expected a minimal length ration  of %.02f, found a similarity of %.02f.", MIN_LRA_HIGH, p.Rules[0].Min_lra)
	}
	if p.Rules[0].Hit_sta != HIT_STA_HIGH {
		t.Error("Wrong hit status for the first default rule.")
	}
	if p.Rules[1].Hit_sta != HIT_STA_NORM {
		t.Error("Wrong hit status for the second default rule.")
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
