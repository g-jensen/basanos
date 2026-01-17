package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDotPrinter_PrintScenarioResult_Pass(t *testing.T) {
	buf := &bytes.Buffer{}
	dot := &dotPrinter{writer: buf, color: noopColorizer{}}

	dot.printScenarioResult("pass")

	assert.Equal(t, ".", buf.String())
}

func TestDotPrinter_PrintScenarioResult_Fail(t *testing.T) {
	buf := &bytes.Buffer{}
	dot := &dotPrinter{writer: buf, color: noopColorizer{}}

	dot.printScenarioResult("fail")

	assert.Equal(t, "F", buf.String())
}

func TestDotPrinter_PrintScenarioResult_WithColor(t *testing.T) {
	buf := &bytes.Buffer{}
	dot := &dotPrinter{writer: buf, color: ansiColorizer{}}

	dot.printScenarioResult("pass")

	assert.Equal(t, "\033[32m.\033[0m", buf.String())
}

func TestDotPrinter_ContextMethods_AreNoops(t *testing.T) {
	buf := &bytes.Buffer{}
	dot := &dotPrinter{writer: buf, color: noopColorizer{}}

	dot.printContextEnter("Some Context")
	dot.printContextExit()
	dot.printScenarioEnter("Some Scenario")

	assert.Equal(t, "", buf.String())
}

func TestVerbosePrinter_PrintContextEnter(t *testing.T) {
	buf := &bytes.Buffer{}
	verbose := &verbosePrinter{writer: buf, color: noopColorizer{}}

	verbose.printContextEnter("My Context")

	assert.Equal(t, "My Context\n", buf.String())
}

func TestVerbosePrinter_PrintContextEnter_IncrementsDepth(t *testing.T) {
	buf := &bytes.Buffer{}
	verbose := &verbosePrinter{writer: buf, color: noopColorizer{}}

	verbose.printContextEnter("Parent")
	verbose.printContextEnter("Child")

	assert.Equal(t, "Parent\n  Child\n", buf.String())
}

func TestVerbosePrinter_PrintContextExit_DecrementsDepth(t *testing.T) {
	buf := &bytes.Buffer{}
	verbose := &verbosePrinter{writer: buf, color: noopColorizer{}}

	verbose.printContextEnter("Parent")
	verbose.printContextEnter("Child")
	verbose.printContextExit()
	verbose.printContextEnter("Sibling")

	assert.Equal(t, "Parent\n  Child\n  Sibling\n", buf.String())
}

func TestVerbosePrinter_PrintScenarioEnter_StoresName(t *testing.T) {
	buf := &bytes.Buffer{}
	verbose := &verbosePrinter{writer: buf, color: noopColorizer{}}

	verbose.printScenarioEnter("My Scenario")

	assert.Equal(t, "My Scenario", verbose.pendingName)
	assert.Equal(t, "", buf.String())
}

func TestVerbosePrinter_PrintScenarioResult_OutputsNameAndStatus(t *testing.T) {
	buf := &bytes.Buffer{}
	verbose := &verbosePrinter{writer: buf, color: noopColorizer{}}

	verbose.printScenarioEnter("Login works")
	verbose.printScenarioResult("pass")

	assert.Equal(t, "Login works .\n", buf.String())
}

func TestVerbosePrinter_PrintScenarioResult_RespectsDepth(t *testing.T) {
	buf := &bytes.Buffer{}
	verbose := &verbosePrinter{writer: buf, color: noopColorizer{}}

	verbose.printContextEnter("Parent")
	verbose.printScenarioEnter("Child scenario")
	verbose.printScenarioResult("pass")

	assert.Contains(t, buf.String(), "  Child scenario .\n")
}

func TestVerbosePrinter_PrintScenarioResult_WithColor(t *testing.T) {
	buf := &bytes.Buffer{}
	verbose := &verbosePrinter{writer: buf, color: ansiColorizer{}}

	verbose.printScenarioEnter("Login works")
	verbose.printScenarioResult("pass")

	assert.Equal(t, "\033[32mLogin works\033[0m\n", buf.String())
}
