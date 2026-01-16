package main

import (
	"os"

	"basanos/internal/assert"
)

func main() {
	os.Exit(assert.RunCLI(os.Args[1:], os.Stdin, os.Stdout,
		assert.ResolveBothValues, wrapEquals))
}

func wrapEquals(first, second string) assert.AssertResult {
	return assert.Equals(first, second)
}
