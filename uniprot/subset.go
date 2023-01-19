package uniprot

type Subset struct {
	Ekeep string
	Eskip string
	Tkeep string
	Tskip string
	Dkeep string
	Dskip string
	Lmin  int
}

type SubsetWriter struct {
	Writer *Writer
}

func NewSubsetWriter(o string) *SubsetWriter {
	return &SubsetWriter{NewWriter(o)}
}

// Recorder routine
func (sww *SubsetWriter) RecordEntry(ec chan *[]string, re chan int) {
	nrec := 0

	for e := range ec {
		sww.Writer.WriteStrings(e)
		sww.Writer.WriteEntryEnd()
		sww.Writer.PanicOnError()
		nrec++
	}

	// Throw the number of recorded entries
	re <- nrec
}

func (s *Subset) ParseFile(ec chan *[]string, th chan int, in string) {
	// Create a reader
	swr := NewReader(in)
	swr.PanicOnError()
	defer swr.Close()

	ntot := 0

	for swr.Next() {
		// Parse the entry
		e := swr.Parse()
		ntot++

		if e.Length < s.Lmin {
			continue
		}

		if s.Eskip != "" {
			if e.TestEvidence(s.Eskip) {
				continue
			}
		}

		if s.Tskip != "" {
			if e.TestTaxonomy(s.Tskip) {
				continue
			}
		}

		if s.Dskip != "" {
			if e.TestDescription(s.Dskip) {
				continue
			}
		}

		if s.Ekeep != "" {
			if !e.TestEvidence(s.Ekeep) {
				continue
			}
		}

		if s.Tkeep != "" {
			if !e.TestTaxonomy(s.Tkeep) {
				continue
			}
		}

		if s.Dkeep != "" {
			if !e.TestDescription(s.Dkeep) {
				continue
			}
		}

		// Copy the pointer (otherwize it is lost before writing...)
		var tmp []string
		tmp = *swr.GetData()

		ec <- &tmp
	}

	// Throw the number of scanned entries
	th <- ntot
}

func (s *Subset) LightParseFile(ec chan *[]string, th chan int, in string) {
	// Create a reader
	swr := NewReader(in)
	swr.PanicOnError()
	defer swr.Close()

	ntot := 0

	for swr.Next() {
		// Parse the entry
		e := swr.LightParse()
		ntot++

		if e.Length < s.Lmin {
			continue
		}

		if s.Eskip != "" {
			if e.TestEvidence(s.Eskip) {
				continue
			}
		}

		if s.Tskip != "" {
			if e.TestTaxonomy(s.Tskip) {
				continue
			}
		}

		if s.Ekeep != "" {
			if !e.TestEvidence(s.Ekeep) {
				continue
			}
		}

		if s.Tkeep != "" {
			if !e.TestTaxonomy(s.Tkeep) {
				continue
			}
		}

		// Copy the pointer (otherwize it is lost before writing...)
		var tmp []string
		tmp = *swr.GetData()

		ec <- &tmp
	}

	// Throw the number of scanned entries
	th <- ntot
}
