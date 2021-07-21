package swiss

import (
	"bufio"
	"os"
	"regexp"

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
			r.data = append(r.data, line)
			if r.reend.MatchString(line) {
				return true
			}
		}
	}

	// Missing entry end => return false
	return false
}

func (r *Reader) Close() {
	r.closer.Close()
}
