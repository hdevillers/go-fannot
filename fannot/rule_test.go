package fannot

import (
	"testing"
)

// Test the rule class
func TestRuleCheck(t *testing.T) {
	minSim := float64(66.8)
	minLra := float64(57.3)

	r := Rule{minSim, minLra, "Test", false, 0}

	// Testing parameter equality to threshold
	if !r.Test(minSim, minLra) {
		t.Fatal(`Testing parameter values equal to thresholds, it should return true.`)
	}

	// Testing higher values
	if !r.Test(minSim+10.0, minLra+10.0) {
		t.Fatal(`Testing parameter values higher than thresholds, it should return true.`)
	}

	// Testing lower similarity
	if r.Test(minSim-10.0, minLra) {
		t.Fatal(`Similarity lower than the threshold, it should return false.`)
	}

	// Testing lower lra
	if r.Test(minSim, minLra-10.0) {
		t.Fatal(`Length ratio lower than the threshold, it should return false.`)
	}

	// Testing both lower similarity and lra
	if r.Test(minSim-10.0, minLra-10.0) {
		t.Fatal(`Both parameter values lower than the thresholds, it should return false.`)
	}
}
