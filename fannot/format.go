package fannot

import "regexp"

const (
	DEFAULT_FORMAT_DATA string = "Null"
)

type Format struct {
	Template string
	Fields   [][]string
}

func NewFormat(input string) *Format {
	// input should start by [ and end by ]
	re := regexp.MustCompile(`^\[.*\]$`)
	if re.MatchString(input) {
		input = input[1:(len(input) - 1)]
	} else {
		panic("The input format string should start with a '[' and end with a ']' ")
	}

	var f Format

	if input == "" {
		f.Template = DEFAULT_FORMAT_DATA
	} else {

	}

	return &f
}
