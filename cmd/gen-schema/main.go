package main

import (
	"fmt"
	"os"
)

func main() {
	source, err := os.ReadFile("internal/event/event.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading event.go: %v\n", err)
		os.Exit(1)
	}

	schema, err := GenerateSchema(string(source))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(schema))
}
