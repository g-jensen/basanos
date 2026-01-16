package main

import (
	"fmt"
	"os"

	"basanos/internal/assert"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: assert_lte <left> <right>")
		fmt.Fprintln(os.Stderr, "  Asserts left <= right")
		os.Exit(1)
	}

	result := assert.LessThanOrEqual(os.Args[1], os.Args[2])
	fmt.Print(result.Format())

	if !result.Passed {
		os.Exit(1)
	}
}
