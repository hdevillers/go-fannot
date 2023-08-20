package fannot

import (
	"bufio"
	"encoding/json"
	"os"
)

// Default thresholds
const (
	NB_HIT_CHECK int     = 3
	UNKNOWN_FUNC string  = "hypothetical protein"
	DFT_NOTE     string  = "hypothetical protein"
	DFT_PRODUCT  string  = "hypothetical protein"
	DFT_GENENAME string  = ""
	DFT_FUNCTION string  = ""
	DFT_MAXSTSOW int     = 1
	DFT_MINSIDOW float64 = 5.0
	TPL_NOTE     string  = "{Prefix}||{DbName}|{DbId} ||{Species} ||{LocusTag} ||{GeneName} ||{LongDesc}"
	TPL_PRODUCT  string  = "{ShortDesc}::ToLwr::GnPn"
	TPL_GENENAME string  = "{GeneName}"
	TPL_FUNCTION string  = ""
)

// Global parameter object
type Param struct {
	DefaultNote      string
	DefaultProduct   string
	DefaultGeneName  string
	DefaultFunction  string
	TemplateNote     string
	TemplateProduct  string
	TemplateGeneName string
	TemplateFunction string
	NbHitCheck       int
	Rules            []Rule
	MaxStatusOW      int
	MinSimDiffOW     float64
}

// Create a new parameter object with default values
func NewParam() *Param {
	var p Param

	// Init. default annotations
	p.DefaultNote = DFT_NOTE
	p.DefaultProduct = DFT_PRODUCT
	p.DefaultGeneName = DFT_GENENAME
	p.DefaultFunction = DFT_FUNCTION
	p.MaxStatusOW = DFT_MAXSTSOW
	p.MinSimDiffOW = DFT_MINSIDOW

	// Init. templates
	p.TemplateNote = TPL_NOTE
	p.TemplateProduct = TPL_PRODUCT
	p.TemplateGeneName = TPL_GENENAME
	p.TemplateFunction = TPL_FUNCTION

	// Nb of best hit to keep/check
	p.NbHitCheck = NB_HIT_CHECK

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
