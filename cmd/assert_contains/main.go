package main

import (
	"fmt"
	"os"

	"basanos/internal/assert"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: assert_contains <needle> <haystack>")
		fmt.Fprintln(os.Stderr, "  Arguments can be file paths or literal values")
		fmt.Fprintln(os.Stderr, "  If a file exists at the path, its contents are used")
		os.Exit(1)
	}

	needle, err := assert.ResolveValue(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading needle: %v\n", err)
		os.Exit(1)
	}

	haystack, err := assert.ResolveValue(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading haystack: %v\n", err)
		os.Exit(1)
	}

	result := assert.Contains(needle, haystack)
	fmt.Print(result.Format())

	if !result.Passed {
		os.Exit(1)
	}
}
