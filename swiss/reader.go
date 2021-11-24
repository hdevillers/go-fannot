package swiss

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	gzip "github.com/klauspost/pgzip"
)

type Reader struct {
	closer  FileCloser
	scanner *bufio.Scanner
	data    []string
	err     error
	restart *regexp.Regexp
	reend   *regexp.Regexp
	reac    *regexp.Regexp
	regn    *regexp.Regexp
	relt    *regexp.Regexp
	rede    *regexp.Regexp
	rese    *regexp.Regexp
	resc    *regexp.Regexp
}

func NewReader(file string) *Reader {
	// Open the file
	f, err := os.Open(file)
	if err != nil {
		return &Reader{err: err}
	}

	// Setup regex to detecte beginning and end of an entry
	restart := regexp.MustCompile(`^ID   `)
	reend := regexp.MustCompile(`^\/\/$`)
	reac := regexp.MustCompile(`;\s?`)
	regn := regexp.MustCompile(`Name=(\w+)`)
	relt := regexp.MustCompile(`OrderedLocusNames=([\w\-]+)`)
	rede := regexp.MustCompile(`^RecName\: Full=`)
	rese := regexp.MustCompile(`\s`)
	resc := regexp.MustCompile(`\;`)

	// If the dat file has a 'gz' extention, then use zlib
	var testGZ = regexp.MustCompile(`\.gz$`)
	if testGZ.MatchString(file) {
		// Use zlib
		fgzip, err := gzip.NewReader(f)
		if err != nil {
			return &Reader{err: err}
		}
		return &Reader{
			closer:  fgzip,
			scanner: bufio.NewScanner(fgzip),
			restart: restart,
			reend:   reend,
			reac:    reac,
			regn:    regn,
			relt:    relt,
			rede:    rede,
			rese:    rese,
			resc:    resc,
		}
	} else {
		// Regular text file
		return &Reader{
			closer:  f,
			scanner: bufio.NewScanner(f),
			restart: restart,
			reend:   reend,
			reac:    reac,
			regn:    regn,
			relt:    relt,
			rede:    rede,
			rese:    rese,
			resc:    resc,
		}
	}
}

func (r *Reader) Next() bool {
	// Empty the current data
	r.data = nil

	// Scan a first line
	test := r.scanner.Scan()
	if !test {
		// The file is probable EOF
		return false
	}
	if r.scanner.Err() != nil {
		// The scan failed (or EOF)
		return false
	}

	line := r.scanner.Text()
	if r.restart.MatchString(line) {
		r.data = append(r.data, line)
		for r.scanner.Scan() {
			// We suppose no error if line 1 is ok
			line = r.scanner.Text()
			if r.reend.MatchString(line) {
				return true
			}
			// Do not append the // line!
			r.data = append(r.data, line)
		}
	}

	// Missing entry end => return false
	return false
}

func (r *Reader) Close() {
	r.closer.Close()
}

func (r *Reader) PanicOnError() {
	if r.err != nil {
		panic(r.err)
	}
}

func (r *Reader) GetData() *[]string {
	return &r.data
}

/*
	LightParse parse only essential elements for
	subset selection and do not extract the complete
	data from the entry.
*/
func (r *Reader) LightParse() *Entry {
	if len(r.data) == 0 {
		panic("No data read. You must call Next() method first.")
	}

	// Initialize the new entry
	var entry Entry

	// Split data by line types into a map
	mdata := make(map[string]string)
	for _, line := range r.data {
		key := line[0:2]
		mdata[key] += line[5:]
	}

	// Retrieve the length of the protein
	le := regexp.MustCompile(`(\d+) AA`).FindStringSubmatch(mdata["ID"])
	if len(le) != 2 {
		panic(fmt.Sprintf("Failed to retrieve the length of the protein (%s).", mdata["ID"]))
	}
	var err error
	entry.Length, err = strconv.Atoi(le[1])
	if err != nil {
		panic(err)
	}

	// Organisme and phylum
	entry.Organism = mdata["OS"]
	entry.Phylum = mdata["OC"]

	// Entry evidence level
	entry.Evidence = mdata["PE"][0:1]

	return &entry
}

func (r *Reader) Parse() *Entry {
	if len(r.data) == 0 {
		panic("No data read. You must call Next() method first.")
	}

	// Initialize the new entry
	var entry Entry

	// Split data by line types into a map
	mdata := make(map[string]string)
	for _, line := range r.data {
		key := line[0:2]
		mdata[key] += line[5:]
	}

	// Retrieve the length of the protein
	le := regexp.MustCompile(`(\d+) AA`).FindStringSubmatch(mdata["ID"])
	if len(le) != 2 {
		panic(fmt.Sprintf("Failed to retrieve the length of the protein (%s).", mdata["ID"]))
	}
	var err error
	entry.Length, err = strconv.Atoi(le[1])
	if err != nil {
		panic(err)
	}

	// Retrieve accession number
	// NOTE: For sake of simplicify, only the first accession will be kept
	ac := r.reac.Split(mdata["AC"], -1)
	entry.Access = ac[0]

	// Get information from CC lines
	if mdata["CC"] != "" {
		// Split by "-!- "
		cc := strings.Split(mdata["CC"], "-!- ")

		// Scan each cc type
	CCVAL:
		for _, ccv := range cc {
			// Look for FUNCTION
			if len(ccv) > 5 {
				if ccv[0:5] == "FUNCT" {
					ccv = commentFunctionCleanup(ccv)
					entry.Function = ccv
					break CCVAL
				}
			}
		}
	}

	// Retieve gene name and locus tag
	if mdata["GN"] != "" {
		// Retrieve the gene name (can be null)
		gn := r.regn.FindStringSubmatch(mdata["GN"])
		if gn != nil {
			entry.Name = gn[1]
		}

		// Retrieve the locus tag (can be null)
		lt := r.relt.FindStringSubmatch(mdata["GN"])
		if lt != nil {
			entry.Locus = lt[1]
		}
	}

	// Retrieve the functional annotation
	// NOTE: we suppose that all DE entries start with "RecName: Full="
	if mdata["DE"] != "" {
		de := r.resc.Split(mdata["DE"], 2)
		if r.rede.MatchString(de[0]) {
			// Delete useless accession numbers (ex. {ECO:XXXX})
			tmpDesc := strings.Split(de[0][14:], " {")
			entry.Desc = tmpDesc[0]
		}
	}

	// Organisme and phylum
	entry.Organism = mdata["OS"]
	entry.Phylum = mdata["OC"]

	// Protein sequence
	entry.Sequence = r.rese.ReplaceAllString(mdata["  "], "")

	// Entry evidence level
	entry.Evidence = mdata["PE"][0:1]

	return &entry
}

func commentFunctionCleanup(ccv string) string {
	// Delete duplicated spaces
	ccv = regexp.MustCompile(`    `).ReplaceAllString(ccv[10:], " ")

	// Delete pubmed indications
	ccv = regexp.MustCompile(` \(Pub[\w\d\:\, ]+\)`).ReplaceAllString(ccv, "")

	// Delete the possible {ECO:...} at the end
	ccv = strings.Split(ccv, " {")[0]

	// Split sentences
	sen := strings.Split(ccv, ". ")
	for i := range sen {
		if regexp.MustCompile(`[A-Z][a-z ]`).MatchString(sen[i][0:2]) {
			tmp := []rune(sen[i])
			tmp[0] = unicode.ToLower(tmp[0])
			sen[i] = string(tmp)
		}
	}

	// Join sentences
	out := strings.Join(sen, "; ")

	// Delete (By similarity) informations
	out = regexp.MustCompile(` \(By similarity\)`).ReplaceAllString(out, "")

	// Delete final point
	out = regexp.MustCompile(`\.$`).ReplaceAllString(out, "")

	return out
}
