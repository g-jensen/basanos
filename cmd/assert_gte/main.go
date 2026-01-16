package main

import (
	"os"

	"basanos/internal/assert"
)

func main() {
	os.Exit(assert.RunCLI(os.Args[1:], os.Stdin, os.Stdout,
		assert.ResolveLiterals, wrapGreaterThanOrEqual))
}

func wrapGreaterThanOrEqual(first, second string) assert.AssertResult {
	return assert.GreaterThanOrEqual(first, second)
}
