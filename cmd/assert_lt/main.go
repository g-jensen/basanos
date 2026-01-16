package main

import (
	"os"

	"basanos/internal/assert"
)

func main() {
	os.Exit(assert.RunCLI(os.Args[1:], os.Stdin, os.Stdout,
		assert.ResolveLiterals, wrapLessThan))
}

func wrapLessThan(first, second string) assert.AssertResult {
	return assert.LessThan(first, second)
}
