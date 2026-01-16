package main

import (
	"fmt"
	"os"

	"basanos/internal/assert"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: assert_matches <pattern> <target>")
		fmt.Fprintln(os.Stderr, "  pattern: regex pattern (always literal)")
		fmt.Fprintln(os.Stderr, "  target: file path or literal value")
		os.Exit(1)
	}

	pattern := os.Args[1]

	target, err := assert.ResolveValue(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading target: %v\n", err)
		os.Exit(1)
	}

	result := assert.Matches(pattern, target)
	fmt.Print(result.Format())

	if !result.Passed {
		os.Exit(1)
	}
}
