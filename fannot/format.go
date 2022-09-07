package fannot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	DEFAULT_FORMAT_DATA string = "Null"
)

type Format struct {
	Template   []string
	Fields     [][]string
	Empty      bool
	TransGnPn  bool
	TransToLwr bool
}

func NewFormat(input string) *Format {
	var f Format
	f.Empty = false

	// If the input string is empty, then return an empty format object
	if input == "" {
		f.Empty = true
		return &f
	}

	// Init. a Fields object to control entries
	fld := NewFields()

	// Look for possible 'transformer'
	tmp := strings.Split(input, "::")
	// No transformer by default
	f.TransGnPn = false
	f.TransToLwr = false
	for i := 1; i < len(tmp); i++ {
		if tmp[i] == "GnPn" {
			f.TransGnPn = true
		} else if tmp[i] == "ToLwr" {
			f.TransToLwr = true
		} else {
			panic(fmt.Sprintf(`Unsupported transformer: %s`, tmp[i]))
		}
	}

	// split the input into sub-inputs
	subs := strings.Split(tmp[0], "||")

	// Initialize Template and Fields attributes
	f.Template = make([]string, len(subs))
	f.Fields = make([][]string, len(subs))

	// Scan each sub-input
	re := regexp.MustCompile(`\{(\w+)\}`)
	for i, sub := range subs {
		dt := re.FindAllStringSubmatch(sub, -1)

		// If no field found then simply copy the sub-input in the template
		if len(dt) == 0 {
			f.Template[i] = sub
			f.Fields[i] = make([]string, 0)
		} else {
			// Identify fields
			f.Fields[i] = make([]string, len(dt))
			for j, fi := range dt {
				if fld.Exists(fi[1]) {
					su := regexp.MustCompile("{" + fi[1] + "}")
					sub = su.ReplaceAllString(sub, "#"+strconv.Itoa(j))
					f.Fields[i][j] = fi[1]
				} else {
					panic(fmt.Sprintf(`Unsupported field: %s.`, fi[1]))
				}
			}
			// Save the edited template
			f.Template[i] = sub
		}
	}

	return &f
}

func (f *Format) Compile(d *Description) string {
	var out = ""
	if f.Empty {
		return out
	}

SUBTEMP:
	for i, sub := range f.Template {
		// Check for fields to replace
		if len(f.Fields[i]) == 0 {
			// No field, just copy the current sub template
			out = out + sub
		} else {
			// Check if required fields are set (if not skip the sub)
			for j, field := range f.Fields[i] {
				if !d.IsSetField(field) {
					continue SUBTEMP
				}
				// Replace the value in the sub template
				tmp := strings.Split(sub, "#"+strconv.Itoa(j))
				sub = strings.Join(tmp, d.Data[field])
			}
			out = out + sub
		}
		// Check for transformers
		if f.TransToLwr {
			// Lower the first Character of the sentence except if it is an
			// acronym
			if regexp.MustCompile(`^[A-Z][a-z ]`).MatchString(out) {
				tmp := []rune(out)
				tmp[0] = unicode.ToLower(tmp[0])
				out = string(tmp)
			}
		}

		if f.TransGnPn {
			// Replace the gene name by protein name (if set)
			if d.IsSetField("GeneName") {
				re := regexp.MustCompile(d.Data["GeneName"])
				// We suppose that if there is a gene name, the protein name
				// has been automaticaly deduced
				out = re.ReplaceAllString(out, d.Data["ProteinName"])
			}
		}

	}
	return out
}
