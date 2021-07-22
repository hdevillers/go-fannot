package swiss

import (
	"bufio"
	"os"
	"regexp"
	"strconv"

	gzip "github.com/klauspost/pgzip"
)

type FileCloser interface {
	Close() error
}

type Reader struct {
	closer  FileCloser
	scanner *bufio.Scanner
	data    []string
	err     error
	restart *regexp.Regexp
	reend   *regexp.Regexp
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
		}
	} else {
		// Regular text file
		return &Reader{
			closer:  f,
			scanner: bufio.NewScanner(f),
			restart: restart,
			reend:   reend,
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

	// Retrieve accession number
	// NOTE: For sake of simplicify, only the first accession will be kept
	ac := regexp.MustCompile(`;\s?`).Split(mdata["AC"], -1)
	entry.Access = ac[0]

	// Retieve gene name and locus tag
	if mdata["GN"] != "" {
		// Retrieve the gene name (can be null)
		gn := regexp.MustCompile(`Name=(\w+)`).FindStringSubmatch(mdata["GN"])
		if gn != nil {
			entry.Name = gn[1]
		}

		// Retrieve the locus tag (can be null)
		lt := regexp.MustCompile(`OrderedLocusNames=([\w\-]+)`).FindStringSubmatch(mdata["GN"])
		if lt != nil {
			entry.Locus = lt[1]
		}
	}

	// Retrieve the functional annotation
	// NOTE: we suppose that all DE entries start with "RecName: Full="
	if mdata["DE"] != "" {
		de := regexp.MustCompile(`;`).Split(mdata["DE"], 2)
		if regexp.MustCompile(`^RecName\: Full=`).MatchString(de[0]) {
			entry.Desc = de[0][14:]
		}
	}

	// Organisme and phylum
	entry.Organism = mdata["OS"]
	entry.Phylum = mdata["OC"]

	// Protein sequence
	entry.Sequence = regexp.MustCompile(`\s`).ReplaceAllString(mdata["  "], "")

	// Entry evidence level
	entry.Evidence, _ = strconv.Atoi(mdata["PE"][0:1])

	return &entry
}
