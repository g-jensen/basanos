package runner

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveAssertionArgs_LogicalStdout(t *testing.T) {
	captured := CapturedOutput{
		Stdout:   "hello world",
		Stderr:   "",
		ExitCode: 0,
	}
	env := map[string]string{"RUN_OUTPUT": "/path/to/run"}

	first, second, err := resolveAssertionArgs("assert_contains expected.txt ${RUN_OUTPUT}/stdout", captured, env)

	require.NoError(t, err)
	assert.Equal(t, "expected.txt", first)
	assert.Equal(t, "hello world", second)
}

func TestResolveAssertionArgs_LogicalExitCode(t *testing.T) {
	captured := CapturedOutput{
		Stdout:   "",
		Stderr:   "",
		ExitCode: 42,
	}
	env := map[string]string{"RUN_OUTPUT": "/path/to/run"}

	first, second, err := resolveAssertionArgs("assert_equals 0 ${RUN_OUTPUT}/exit_code", captured, env)

	require.NoError(t, err)
	assert.Equal(t, "0", first)
	assert.Equal(t, "42", second)
}

func TestResolveAssertionArgs_LiteralValues(t *testing.T) {
	captured := CapturedOutput{}
	env := map[string]string{}

	first, second, err := resolveAssertionArgs("assert_equals expected actual", captured, env)

	require.NoError(t, err)
	assert.Equal(t, "expected", first)
	assert.Equal(t, "actual", second)
}

func TestResolveAssertionArgs_MixedLiteralAndLogical(t *testing.T) {
	captured := CapturedOutput{
		Stdout:   "",
		Stderr:   "",
		ExitCode: 0,
	}
	env := map[string]string{"RUN_OUTPUT": "/path/to/run"}

	first, second, err := resolveAssertionArgs("assert_equals 0 ${RUN_OUTPUT}/exit_code", captured, env)

	require.NoError(t, err)
	assert.Equal(t, "0", first)
	assert.Equal(t, "0", second)
}

func TestResolveAssertionArgs_ReadsFileContents(t *testing.T) {
	tempFile, err := os.CreateTemp("", "expected-*.fixture")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("expected content")
	require.NoError(t, err)
	tempFile.Close()

	captured := CapturedOutput{
		Stdout:   "expected content",
		Stderr:   "",
		ExitCode: 0,
	}
	env := map[string]string{"RUN_OUTPUT": "/path/to/run"}

	first, second, err := resolveAssertionArgs("assert_equals "+tempFile.Name()+" ${RUN_OUTPUT}/stdout", captured, env)

	require.NoError(t, err)
	assert.Equal(t, "expected content", first)
	assert.Equal(t, "expected content", second)
}

func TestParseCommandArgs_SimpleArgs(t *testing.T) {
	executable, args := parseCommandArgs("cmd arg1 arg2")

	assert.Equal(t, "cmd", executable)
	assert.Equal(t, []string{"arg1", "arg2"}, args)
}

func TestParseCommandArgs_DoubleQuotedArg(t *testing.T) {
	executable, args := parseCommandArgs(`cmd "200 OK" arg2`)

	assert.Equal(t, "cmd", executable)
	assert.Equal(t, []string{"200 OK", "arg2"}, args)
}

func TestParseCommandArgs_SingleQuotedArg(t *testing.T) {
	executable, args := parseCommandArgs(`cmd '200 OK' arg2`)

	assert.Equal(t, "cmd", executable)
	assert.Equal(t, []string{"200 OK", "arg2"}, args)
}

func TestParseCommandArgs_MixedQuotes(t *testing.T) {
	executable, args := parseCommandArgs(`cmd "first arg" 'second arg'`)

	assert.Equal(t, "cmd", executable)
	assert.Equal(t, []string{"first arg", "second arg"}, args)
}

func TestParseCommandArgs_QuotedWithSpecialChars(t *testing.T) {
	executable, args := parseCommandArgs(`cmd "hello \"world\"" arg`)

	assert.Equal(t, "cmd", executable)
	assert.Equal(t, []string{`hello "world"`, "arg"}, args)
}

func TestParseCommandArgs_EmptyQuotes(t *testing.T) {
	executable, args := parseCommandArgs(`cmd "" arg`)

	assert.Equal(t, "cmd", executable)
	assert.Equal(t, []string{"", "arg"}, args)
}
