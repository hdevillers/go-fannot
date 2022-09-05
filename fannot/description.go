package fannot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hdevillers/go-seq/seq"
)

const (
	DEFAULT_DATA string = "Null"
)

type Description struct {
	Data map[string]string
}

func NewDescription(dn string, s seq.Seq) *Description {
	var d Description
	d.Data = make(map[string]string)

	// Fill db name
	d.Data["DbName"] = dn

	// Parse sequence Id/Description
	d.Data["DbId"] = s.Id
	values := strings.Split(s.Desc, "::")

	// Expected a regular description (without split) or a
	// refdb description (5 splits)
	if len(values) == 1 {
		d.Data["ShortDesc"] = s.Desc
		d.Data["GeneName"] = DEFAULT_DATA
		d.Data["ProteinName"] = DEFAULT_DATA
		d.Data["LocusTag"] = DEFAULT_DATA
		d.Data["Species"] = DEFAULT_DATA
		d.Data["LongDesc"] = s.Desc
	} else if len(values) == 5 {
		// First value: the short desctription
		d.Data["ShortDesc"] = values[0] // Supposed never nil

		// Second value: the gene name (and protein name)
		if values[1] != "" {
			d.Data["GeneName"] = values[1]
			d.Data["ProteinName"] = strings.Title(strings.ToLower(values[1])) + "p"
		} else {
			d.Data["GeneName"] = DEFAULT_DATA
			d.Data["ProteinName"] = DEFAULT_DATA
		}

		// Thrid value: the locus tag
		if values[2] != "" {
			d.Data["LocusTag"] = values[2]
		} else {
			d.Data["LocusTag"] = DEFAULT_DATA
		}

		// Fourth value: the organism name
		if values[3] != "" {
			// Clean-up organism name (delete comments)
			tmpOrg := strings.Split(values[3], " (")
			tmpOrg[0] = regexp.MustCompile(`\.$`).ReplaceAllString(tmpOrg[0], "")
			d.Data["Species"] = tmpOrg[0]
		} else {
			d.Data["Species"] = DEFAULT_DATA
		}

		// Fifth value: the long description
		if values[4] != "" {
			d.Data["LongDesc"] = values[4]
		} else {
			// If no long description provided, return the short one
			d.Data["LongDesc"] = values[0]
		}

	} else {
		panic(fmt.Sprintf("Failed to parse sequence description, expected 1 or 5 elements, found %d.", len(values)))
	}

	return &d
}

func (d *Description) HasField(f string) bool {
	_, test := d.Data[f]
	return (test)
}

func (d *Description) IsSetField(f string) bool {
	val, test := d.Data[f]
	if test {
		if val == DEFAULT_DATA {
			return false
		} else {
			return true
		}
	}
	return false
}
