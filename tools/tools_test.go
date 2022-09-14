package tools

import (
	"os/exec"
	"regexp"
	"testing"
)

// Test the availability of makeblastdb
func TestMakeblastdb(t *testing.T) {
	out, err := exec.Command("makeblastdb", "-version").Output()
	if err != nil {
		t.Fatalf(`Cannot find makeblastdb tool. Please install BLAST+ suite and place binaries in your PATH.`)
	}

	re := regexp.MustCompile(`([\d\.]+\+)`)
	ver := re.Find(out)

	t.Logf(`Found makeblastdb version: %s`, ver)
}

// Test the availability of blastp
func TestBlastp(t *testing.T) {
	out, err := exec.Command("blastp", "-version").Output()
	if err != nil {
		t.Fatalf(`Cannot find blastp tool. Please install BLAST+ suite and place binaries in your PATH.`)
	}

	re := regexp.MustCompile(`([\d\.]+\+)`)
	ver := re.Find(out)

	t.Logf(`Found blastp version: %s`, ver)
}

// Test the availability of needle
func TestNeedle(t *testing.T) {
	out, err := exec.Command("needle", "-version").CombinedOutput()
	if err != nil {
		t.Fatalf(`Cannot find needle tool. Please install EMBOSS tool suite and place binaries in your PATH.`)
	}

	re := regexp.MustCompile(`([\d\.]+)`)
	ver := re.Find(out)

	t.Logf(`Found needle version: %s`, ver)
}
