package ips

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	D_MAX_EVALUE float64 = 1e-10
)

type IpsEntry struct {
	GeneName string
	KeyValue map[string]string
	Nkeys    int
}

func NewIpsEntry(gn string) *IpsEntry {
	var ie IpsEntry
	ie.GeneName = gn
	ie.KeyValue = make(map[string]string)
	ie.Nkeys = 0
	return &ie
}

type Ips struct {
	Data   map[string]IpsEntry
	NGenes int
	Evalue float64
}

func NewIps() *Ips {
	var i Ips
	i.Data = make(map[string]IpsEntry)
	i.NGenes = 0
	return &i
}

func CleanUpAnnot(a string) string {
	// Trim spaces
	b := []byte(strings.TrimSpace(a))

	// Lower the first cap (if not an acronyme)
	if regexp.MustCompile(`^[A-Z][a-z ]`).Match(b) {
		b[0] = []byte(strings.ToLower(a))[0]
	}

	return string(b)
}

func (i *Ips) LoadIpsData(f string) error {
	fh, err := os.Open(f)
	if err != nil {
		return err
	}
	fb := bufio.NewScanner(fh)

	geneId := ""

	for fb.Scan() {
		err = fb.Err()
		if err != nil {
			return err
		}

		line := fb.Text()
		elem := strings.Split(line, "\t")

		// Only lines with 13 elements contain a IPR ID
		if len(elem) == 13 {
			// Check the E-value (or score)
			eval, err := strconv.ParseFloat(elem[8], 64)
			if err != nil {
				return err
			}

			if eval <= i.Evalue {
				if geneId != elem[0] {
					// Create a new entry for the gene
					geneId = elem[0]
					i.Data[geneId] = *NewIpsEntry(geneId)
					i.NGenes++
				}

				// Check if the IPR is already stored
				_, set := i.Data[geneId].KeyValue[elem[11]]

				if !set {
					i.Data[geneId].KeyValue[elem[11]] = CleanUpAnnot(elem[12])
					i.Data[geneId].Nkeys++
				}
			}
		}

	}

	return nil
}
