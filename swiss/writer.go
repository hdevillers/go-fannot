package swiss

import (
	"bufio"
	"os"
	"regexp"

	gzip "github.com/klauspost/pgzip"
)

type Writer struct {
	closer FileCloser
	writer FileWriter
	err    error
}

func NewWriter(file string) *Writer {
	f, err := os.Create(file)
	if err != nil {
		return &Writer{err: err}
	}

	// If the provided file has a 'gz' extention, then compress
	testGZ := regexp.MustCompile(`\.gz$`)
	if testGZ.MatchString(file) {
		fgz := gzip.NewWriter(f)
		return &Writer{
			closer: fgz,
			writer: fgz,
		}
	} else {
		return &Writer{
			closer: f,
			writer: bufio.NewWriter(f),
		}
	}
}

func (w *Writer) Close() {
	// Flush before closing
	w.err = w.writer.Flush()
	w.PanicOnError()
	w.closer.Close()
}

func (w *Writer) PanicOnError() {
	if w.err != nil {
		panic(w.err)
	}
}

func (w *Writer) WriteStrings(s *[]string) {
	for i := range *s {
		_, w.err = w.writer.Write([]byte((*s)[i]))
		_, w.err = w.writer.Write([]byte{'\n'})
	}
}

func (w *Writer) WriteEntryEnd() {
	_, w.err = w.writer.Write([]byte{'/', '/', '\n'})
}
