package fannot

import (
	"bufio"
	"encoding/json"
	"os"
)

// Default thresholds
const (
	N_BEST_HITS  int     = 3
	MIN_LRA_HIGH float64 = 0.8
	MIN_SIM_HIGH float64 = 80.0
	MIN_LRA_NORM float64 = 0.7
	MIN_SIM_NORM float64 = 50.0
	UNKNOWN_FUNC string  = "hypothetical protein"
	PRE_SIM_HIGH string  = "highly similar to"
	PRE_SIM_NORM string  = "similar to"
	CPY_GEN_HIGH bool    = true
	CPY_GEN_NORM bool    = false
)

type Rule struct {
	Min_sim float64
	Min_lra float64
	Pre_ann string
	Cpy_gen bool
}

type Param struct {
	Unk_ann string
	Nbh_chk int
	Rules   []Rule
}

func NewParam() *Param {
	var p Param

	// Main values
	p.Nbh_chk = N_BEST_HITS
	p.Unk_ann = UNKNOWN_FUNC

	// Prepare rules
	rule_high := Rule{MIN_SIM_HIGH, MIN_LRA_HIGH, PRE_SIM_HIGH, CPY_GEN_HIGH}
	rule_norm := Rule{MIN_SIM_NORM, MIN_LRA_NORM, PRE_SIM_NORM, CPY_GEN_NORM}
	p.Rules = make([]Rule, 2)
	p.Rules[0] = rule_high
	p.Rules[1] = rule_norm

	// Return
	return &p
}

func NewParamFromJson(file string) *Param {
	var p Param

	// Open the file
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create the reader
	fr := bufio.NewReader(f)

	// Create the json decoder
	jr := json.NewDecoder(fr)

	// Decode the entry
	err = jr.Decode(&p)
	if err != nil {
		panic(err)
	}

	return &p
}
