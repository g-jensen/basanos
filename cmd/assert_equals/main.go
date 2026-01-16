package main

import (
	"fmt"
	"os"

	"basanos/internal/assert"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: assert_equals <expected> <actual>")
		fmt.Fprintln(os.Stderr, "  Arguments can be file paths or literal values")
		fmt.Fprintln(os.Stderr, "  If a file exists at the path, its contents are used")
		os.Exit(1)
	}

	expected, err := assert.ResolveValue(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading expected: %v\n", err)
		os.Exit(1)
	}

	actual, err := assert.ResolveValue(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading actual: %v\n", err)
		os.Exit(1)
	}

	result := assert.Equals(expected, actual)
	fmt.Print(result.Format())

	if !result.Passed {
		os.Exit(1)
	}
}
