package runner

import (
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
