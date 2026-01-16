package sink

import (
	"testing"

	"basanos/internal/event"
	"basanos/internal/testutil/fs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSink_WritesOutput(t *testing.T) {
	memFS := fs.NewMemoryFS()
	sink := NewFileSink(memFS, "2026-01-15_143022")

	sink.Emit(event.NewScenarioRunStartEvent("2026-01-15_143022", "basic_http/login"))
	sink.Emit(event.NewOutputEvent("2026-01-15_143022", "stdout", "hello\n"))

	content, err := memFS.ReadFile("2026-01-15_143022/basic_http/login/_run/stdout")
	require.NoError(t, err)
	assert.Equal(t, "hello\n", string(content))
}

func TestFileSink_WritesHookOutput(t *testing.T) {
	memFS := fs.NewMemoryFS()
	runID := "2026-01-15_143022"
	sink := NewFileSink(memFS, runID)

	sink.Emit(event.NewHookStartEvent(runID, "basic_http", "before", ""))
	sink.Emit(event.NewOutputEvent(runID, "stdout", "starting server\n"))
	sink.Emit(event.NewHookEndEvent(runID, "basic_http", "before", "", 0))

	stdoutContent, err := memFS.ReadFile(runID + "/basic_http/before/stdout")
	require.NoError(t, err)
	assert.Equal(t, "starting server\n", string(stdoutContent))

	exitCodeContent, err := memFS.ReadFile(runID + "/basic_http/before/exit_code")
	require.NoError(t, err)
	assert.Equal(t, "0", string(exitCodeContent))
}

func TestFileSink_WritesAssertionOutput(t *testing.T) {
	memFS := fs.NewMemoryFS()
	runID := "2026-01-15_143022"
	sink := NewFileSink(memFS, runID)

	sink.Emit(event.NewAssertionStartEvent(runID, "basic_http/login", 0, "assert_equals 0 exit_code"))
	sink.Emit(event.NewOutputEvent(runID, "stdout", "PASS\n"))
	sink.Emit(event.NewAssertionEndEvent(runID, "basic_http/login", 0, 0))

	stdoutContent, err := memFS.ReadFile(runID + "/basic_http/login/_assertions/0/stdout")
	require.NoError(t, err)
	assert.Equal(t, "PASS\n", string(stdoutContent))

	exitCodeContent, err := memFS.ReadFile(runID + "/basic_http/login/_assertions/0/exit_code")
	require.NoError(t, err)
	assert.Equal(t, "0", string(exitCodeContent))
}

func TestFileSink_WritesScenarioRunExitCode(t *testing.T) {
	memFS := fs.NewMemoryFS()
	runID := "2026-01-15_143022"
	sink := NewFileSink(memFS, runID)

	sink.Emit(event.NewScenarioRunStartEvent(runID, "basic_http/login"))
	sink.Emit(event.NewScenarioRunEndEvent(runID, "basic_http/login", 0))

	exitCodeContent, err := memFS.ReadFile(runID + "/basic_http/login/_run/exit_code")
	require.NoError(t, err)
	assert.Equal(t, "0", string(exitCodeContent))
}

func TestFileSink_AppendsOutput(t *testing.T) {
	memFS := fs.NewMemoryFS()
	runID := "2026-01-15_143022"
	sink := NewFileSink(memFS, runID)

	sink.Emit(event.NewScenarioRunStartEvent(runID, "basic_http/login"))
	sink.Emit(event.NewOutputEvent(runID, "stdout", "line1\n"))
	sink.Emit(event.NewOutputEvent(runID, "stdout", "line2\n"))

	content, err := memFS.ReadFile(runID + "/basic_http/login/_run/stdout")
	require.NoError(t, err)
	assert.Equal(t, "line1\nline2\n", string(content))
}

func TestFileSink_RunDirectoryIsPrefixedWithUnderscore(t *testing.T) {
	memFS := fs.NewMemoryFS()
	runID := "2026-01-15_143022"
	sink := NewFileSink(memFS, runID)

	sink.Emit(event.NewScenarioRunStartEvent(runID, "basic_http/login"))
	sink.Emit(event.NewOutputEvent(runID, "stdout", "hello\n"))

	content, err := memFS.ReadFile(runID + "/basic_http/login/_run/stdout")
	require.NoError(t, err)
	assert.Equal(t, "hello\n", string(content))
}
