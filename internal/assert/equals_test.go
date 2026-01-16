package assert

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals_IdenticalStrings_Pass(t *testing.T) {
	result := Equals("hello", "hello")

	assert.True(t, result.Passed)
}

func TestEquals_DifferentStrings_Fail(t *testing.T) {
	result := Equals("hello", "world")

	assert.False(t, result.Passed)
	assert.Equal(t, "hello", result.Expected)
	assert.Equal(t, "world", result.Actual)
}

func TestEquals_GeneratesDiff(t *testing.T) {
	result := Equals("line1\nline2\nline3", "line1\nchanged\nline3")

	assert.NotEmpty(t, result.Diff)
	assert.Contains(t, result.Diff, "-line2")
	assert.Contains(t, result.Diff, "+changed")
}

func TestResult_Format_Pass(t *testing.T) {
	result := Result{Passed: true}

	output := result.Format()

	assert.Contains(t, output, "PASS")
}

func TestResult_Format_Fail(t *testing.T) {
	result := Result{
		Passed:   false,
		Expected: "hello",
		Actual:   "world",
		Diff:     "-hello\n+world\n",
	}

	output := result.Format()

	assert.Contains(t, output, "FAIL")
	assert.Contains(t, output, "Expected:")
	assert.Contains(t, output, "hello")
	assert.Contains(t, output, "Actual:")
	assert.Contains(t, output, "world")
	assert.Contains(t, output, "Diff:")
}

func TestEquals_DiffHasNoFileHeaders(t *testing.T) {
	result := Equals("hello", "world")

	assert.NotContains(t, result.Diff, "--- expected")
	assert.NotContains(t, result.Diff, "+++ actual")
	assert.Contains(t, result.Diff, "@@ ")
	assert.Contains(t, result.Diff, "-hello")
	assert.Contains(t, result.Diff, "+world")
}

func TestResolveValue_ReturnsLiteralWhenFileDoesNotExist(t *testing.T) {
	result, err := ResolveValue("hello")

	assert.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestResolveValue_ReadsFileWhenExists(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(tempFile, []byte("file content"), 0644)

	result, err := ResolveValue(tempFile)

	assert.NoError(t, err)
	assert.Equal(t, "file content", result)
}
