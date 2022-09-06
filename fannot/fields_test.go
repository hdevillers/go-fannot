package fannot

import (
	"testing"
)

func TestFieldsCheck(t *testing.T) {
	f := NewFields()

	// The field 'GeneName' should exist
	if !f.Exists("GeneName") {
		t.Fatal("The field 'GeneName' should be defined in field labels.")
	}

	// The field 'BlaBla' does not exist
	if f.Exists("BlaBla") {
		t.Fatal("The field 'BlaBla' should not be defined in field labels.")
	}
}
