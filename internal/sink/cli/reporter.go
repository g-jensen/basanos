package cli

import (
	"fmt"
	"io"
	"strings"

	"basanos/internal/event"
	"basanos/internal/sink"
)

type failure struct {
	path   string
	stdout string
	stderr string
}

type Reporter struct {
	writer        io.Writer
	printer       printer
	failures      []failure
	currentStdout strings.Builder
	currentStderr strings.Builder
}

func NewReporter(writer io.Writer, verbose bool, color bool) sink.Sink {
	colors := newColorizer(color)
	var output printer
	if verbose {
		output = &verbosePrinter{writer: writer, color: colors}
	} else {
		output = &dotPrinter{writer: writer, color: colors}
	}
	return &Reporter{writer: writer, printer: output}
}

func (reporter *Reporter) Emit(incoming any) error {
	switch typed := incoming.(type) {
	case *event.ContextEnterEvent:
		reporter.printer.printContextEnter(typed.Name)
	case *event.ContextExitEvent:
		reporter.printer.printContextExit()
	case *event.ScenarioEnterEvent:
		reporter.printer.printScenarioEnter(typed.Name)
		reporter.currentStdout.Reset()
		reporter.currentStderr.Reset()
	case *event.OutputEvent:
		reporter.handleOutput(typed)
	case *event.ScenarioExitEvent:
		reporter.handleScenarioExit(typed)
	case *event.RunEndEvent:
		reporter.printer.finish()
		fmt.Fprintf(reporter.writer, "\n")
		reporter.printFailures()
		reporter.printSummary(typed.Passed, typed.Failed)
	}
	return nil
}

func (reporter *Reporter) handleOutput(output *event.OutputEvent) {
	switch output.Stream {
	case "stdout":
		reporter.currentStdout.WriteString(output.Data)
	case "stderr":
		reporter.currentStderr.WriteString(output.Data)
	}
}

func (reporter *Reporter) handleScenarioExit(exit *event.ScenarioExitEvent) {
	reporter.printer.printScenarioResult(exit.Status)
	if exit.Status == "fail" {
		reporter.failures = append(reporter.failures, failure{
			path:   exit.Path,
			stdout: reporter.currentStdout.String(),
			stderr: reporter.currentStderr.String(),
		})
	}
}

func (reporter *Reporter) printFailures() {
	if len(reporter.failures) == 0 {
		return
	}
	fmt.Fprintf(reporter.writer, "Failures:\n\n")
	for index, fail := range reporter.failures {
		reporter.printFailure(index+1, fail)
	}
	fmt.Fprintf(reporter.writer, "\n")
}

func (reporter *Reporter) printFailure(index int, fail failure) {
	fmt.Fprintf(reporter.writer, "  %d) %s\n", index, fail.path)
	reporter.printIndentedOutput("stdout", fail.stdout)
	reporter.printIndentedOutput("stderr", fail.stderr)
}

func (reporter *Reporter) printIndentedOutput(label, content string) {
	if content == "" {
		return
	}
	fmt.Fprintf(reporter.writer, "     %s:\n", label)
	for _, line := range strings.Split(strings.TrimSuffix(content, "\n"), "\n") {
		fmt.Fprintf(reporter.writer, "       %s\n", line)
	}
}

func (reporter *Reporter) printSummary(passed, failed int) {
	fmt.Fprintf(reporter.writer, "%d passed, %d failed\n", passed, failed)
}
