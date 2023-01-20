package uniprot

import (
	"bufio"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
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

// FTP download
func FtpDownload(url, esum, dest string) {
	// Prepare the url
	url = strings.Replace(url, "ftp://", "", 1)
	path := strings.SplitN(url, "/", 2)

	// Init. FTP connection
	fmt.Println("Connecting to UniProt FTP...")
	cnn, err := ftp.Dial(fmt.Sprintf("%s:21", path[0]), ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		panic(err)
	}

	// Anonymous identification
	err = cnn.Login("anonymous", "anonymous")
	if err != nil {
		panic(err)
	}

	// Get the file
	fmt.Println("Downloading the file...")
	resp, err := cnn.Retr(path[1])
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	// CloseFTP connection
	err = cnn.Quit()
	if err != nil {
		panic(err)
	}

	// Create destination file
	out, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	//defer out.Close()

	// Create an output buffer
	bufo := bufio.NewWriter(out)

	// Complete download and write out the file
	_, err = io.Copy(bufo, resp)
	//_, err = bufo.Write(resp)
	if err != nil {
		panic(err)
	}

	err = bufo.Flush()
	if err != nil {
		panic(err)
	}
	out.Close()

	// Create a byte reader
	in, err := os.Open(dest)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	// Check MD5
	fmt.Printf("Calculating check sum (expecting: %s) of the downloaded file...\n", esum)
	hash := md5.New()
	io.Copy(hash, in)

	// Obtained sum
	osum := fmt.Sprintf("%x", hash.Sum(nil))
	if osum != esum {
		panic(fmt.Sprintf("Check sum failed,\nexpected: %s,\n obtained: %s\n", esum, osum))
	} else {
		fmt.Printf("Check sum (%s) is OK.\n", osum)
	}
}
