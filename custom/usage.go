package custom

import (
	"flag"
	"fmt"
)

const (
	version = "1.0.0"
)

// Custom usage function
func Usage(flagSet flag.FlagSet, tool string, desc string, order *[]string, shand *map[string]string) {
	fmt.Printf("Program: %s\n", tool)
	if desc != "" {
		fmt.Printf("Description: %s.\n", desc)
	}
	fmt.Printf("Version: %s\n\n", version)
	fmt.Printf("Usage of %s:\n", tool)
	for _, name := range *order {
		f := flagSet.Lookup(name)
		sh, ok := (*shand)[name]
		if ok {
			// Argument has a shorthand
			fmt.Printf("\t-%s|-%s\n", sh, f.Name)
		} else {
			// No shorthand
			fmt.Printf("\t-%s\n", f.Name)
		}
		// Check default value
		dv := f.Value.String()
		if dv == "" {
			fmt.Printf("\t\t%s.\n", f.Usage)
		} else {
			fmt.Printf("\t\t%s [default: %s].\n", f.Usage, dv)
		}
	}

}
