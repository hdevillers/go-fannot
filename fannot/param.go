package fannot

import (
	"bufio"
	"encoding/json"
	"os"
)

// Default thresholds
const (
	N_BEST_HITS  int    = 3
	UNKNOWN_FUNC string = "hypothetical protein"
)

// Global parameter object
type Param struct {
	Unk_ann string
	Nbh_chk int
	Rules   []Rule
}

// Create a new parameter object with default values
func NewParam() *Param {
	var p Param

	// Main values
	p.Nbh_chk = N_BEST_HITS
	p.Unk_ann = UNKNOWN_FUNC

	// Prepare rules
	p.Rules = make([]Rule, 2)
	p.Rules[0] = *NewRuleHighlySimilar()
	p.Rules[1] = *NewRuleSimilar()

	// Return
	return &p
}

// Create a new parameter object from a JSON
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

	// Return
	return &p
}
