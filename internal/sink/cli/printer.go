package cli

import (
	"fmt"
	"io"
	"strings"
)

type printer interface {
	printContextEnter(name string)
	printContextExit()
	printScenarioEnter(name string)
	printScenarioResult(status string)
	finish()
}

type dotPrinter struct {
	writer io.Writer
	color  colorizer
}

func (dot *dotPrinter) printContextEnter(name string) {}

func (dot *dotPrinter) printContextExit() {}

func (dot *dotPrinter) printScenarioEnter(name string) {}

func (dot *dotPrinter) printScenarioResult(status string) {
	char := dot.formatChar(status)
	dot.writer.Write([]byte(char))
}

func (dot *dotPrinter) formatChar(status string) string {
	if status == "pass" {
		return dot.color.green(".")
	}
	return dot.color.red("F")
}

func (dot *dotPrinter) finish() {
	dot.writer.Write([]byte("\n"))
}

type verbosePrinter struct {
	writer      io.Writer
	depth       int
	color       colorizer
	pendingName string
}

func (verbose *verbosePrinter) printContextEnter(name string) {
	fmt.Fprintf(verbose.writer, "%s%s\n", strings.Repeat("  ", verbose.depth), name)
	verbose.depth++
}

func (verbose *verbosePrinter) printContextExit() {
	verbose.depth--
}

func (verbose *verbosePrinter) printScenarioEnter(name string) {
	verbose.pendingName = name
}

func (verbose *verbosePrinter) printScenarioResult(status string) {
	indent := strings.Repeat("  ", verbose.depth)
	formattedName := verbose.color.formatName(verbose.pendingName, status)
	fmt.Fprintf(verbose.writer, "%s%s\n", indent, formattedName)
}

func (verbose *verbosePrinter) finish() {}
