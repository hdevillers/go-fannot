package uniprot

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// Metalink path
const (
	MetalinkPath = "https://ftp.uniprot.org/pub/databases/uniprot/current_release/knowledgebase/taxonomic_divisions/RELEASE.metalink"
)

// Metalink xml structure
type Metalink struct {
	XMLName   xml.Name `xml:"metalink"`
	Version   string   `xml:"version"`
	Files     []File   `xml:"files>file"`
	Mirrors   map[string]int
	Databases map[string]int
	Divisions map[string]int
}

// Single file structure
type File struct {
	XMLName xml.Name `xml:"file"`
	Name    string   `xml:"name,attr"`
	Md5     string   `xml:"verification>hash"`
	Urls    []struct {
		Link     string `xml:",innerxml"`
		Location string `xml:"location,attr"`
	} `xml:"resources>url"`
}

func NewMetalink() *Metalink {
	var m Metalink

	// Init. databases
	m.Databases = make(map[string]int)
	m.Databases["sprot"] = 1
	m.Databases["trembl"] = 1

	// Init. mirrors
	m.Mirrors = make(map[string]int)
	m.Mirrors["us"] = 1
	m.Mirrors["uk"] = 1
	m.Mirrors["ch"] = 1

	// Init. divisions
	m.Divisions = make(map[string]int)

	return &m
}

// Load metalink data from url
func (m *Metalink) Retrieve() {
	// HTTP request
	resp, err := http.Get(MetalinkPath)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Bad status: %s\n", resp.Status))
	}

	// Get xml data
	xmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Parse it
	xml.Unmarshal(xmlData, m)

	// Extract the list of available divisions
	fspl := regexp.MustCompile(`[_\.]`)
	for _, file := range m.Files {
		tmp := fspl.Split(file.Name, -1)
		// Skip README files
		if tmp[0] == "uniprot" {
			m.Divisions[tmp[2]] = 1
		}
	}
}

// Check Database value
func (m *Metalink) CheckDatabase(db string) bool {
	_, test := m.Databases[db]
	return test
}

// Check Mirror value
func (m *Metalink) CheckMirror(mir string) bool {
	_, test := m.Mirrors[mir]
	return test
}

// Check Division value
func (m *Metalink) CheckDivision(div string) bool {
	_, test := m.Divisions[div]
	return test
}

// Get an URL according to a DB, a Division and a Mirror
// This function return the URL and the expected MD5 of the file
func (m *Metalink) GetUrl(db, div, mir string) (string, string) {
	link := ""
	md5 := ""
	if !m.CheckDatabase(db) {
		panic(fmt.Sprintf("The database %s is not available.", db))
	}
	if !m.CheckDivision(div) {
		panic(fmt.Sprintf("The taxonomic division %s is not available.", div))
	}
	if !m.CheckMirror(mir) {
		panic(fmt.Sprintf("The mirror %s is not available.", mir))
	}

	// Expected basename
	bn := fmt.Sprintf("uniprot_%s_%s.dat.gz", db, div)

FILES:
	for _, file := range m.Files {
		if file.Name == bn {
			for _, url := range file.Urls {
				if url.Location == mir {
					md5 = file.Md5
					link = url.Link
					break FILES
				}
			}
		}
	}

	return link, md5
}
