package event

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunStartEvent_JSON(t *testing.T) {
	timestamp := time.Date(2026, 1, 15, 14, 30, 22, 0, time.UTC)

	event := NewRunStartEvent("2026-01-15_143022", timestamp)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "run_start", result["event"])
	assert.Equal(t, "2026-01-15_143022", result["run_id"])
	assert.Equal(t, "2026-01-15T14:30:22Z", result["timestamp"])
}

func TestContextEnterEvent_JSON(t *testing.T) {
	timestamp := time.Date(2026, 1, 15, 14, 30, 22, 0, time.UTC)

	event := NewContextEnterEvent("run-123", "basic_http", "Basic HTTP", timestamp)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "context_enter", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http", result["path"])
	assert.Equal(t, "Basic HTTP", result["name"])
	assert.Equal(t, "2026-01-15T14:30:22Z", result["timestamp"])
}

func TestContextExitEvent_JSON(t *testing.T) {
	timestamp := time.Date(2026, 1, 15, 14, 35, 45, 0, time.UTC)

	event := NewContextExitEvent("run-123", "basic_http", timestamp)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "context_exit", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http", result["path"])
	assert.Equal(t, "2026-01-15T14:35:45Z", result["timestamp"])
}

func TestHookStartEvent_JSON(t *testing.T) {
	event := NewHookStartEvent("run-123", "basic_http", "before", "")

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "hook_start", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http", result["path"])
	assert.Equal(t, "before", result["hook"])
	assert.NotContains(t, string(data), "from")
}

func TestHookStartEvent_WithFrom_JSON(t *testing.T) {
	event := NewHookStartEvent("run-123", "basic_http/login", "before_each", "basic_http")

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "hook_start", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http/login", result["path"])
	assert.Equal(t, "before_each", result["hook"])
	assert.Equal(t, "basic_http", result["from"])
}

func TestHookEndEvent_JSON(t *testing.T) {
	event := NewHookEndEvent("run-123", "basic_http", "before", "", 0)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "hook_end", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http", result["path"])
	assert.Equal(t, "before", result["hook"])
	assert.Equal(t, float64(0), result["exit_code"])
}

func TestScenarioEnterEvent_JSON(t *testing.T) {
	timestamp := time.Date(2026, 1, 15, 14, 40, 10, 0, time.UTC)

	event := NewScenarioEnterEvent("run-123", "basic_http/login", "Login works", timestamp)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "scenario_enter", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http/login", result["path"])
	assert.Equal(t, "Login works", result["name"])
	assert.Equal(t, "2026-01-15T14:40:10Z", result["timestamp"])
}

func TestScenarioExitEvent_JSON(t *testing.T) {
	timestamp := time.Date(2026, 1, 15, 14, 42, 30, 0, time.UTC)

	event := NewScenarioExitEvent("run-123", "basic_http/login", "pass", timestamp)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "scenario_exit", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http/login", result["path"])
	assert.Equal(t, "pass", result["status"])
	assert.Equal(t, "2026-01-15T14:42:30Z", result["timestamp"])
}

func TestOutputEvent_JSON(t *testing.T) {
	event := NewOutputEvent("run-123", "stdout", "Hello world\n")

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "output", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "stdout", result["stream"])
	assert.Equal(t, "Hello world\n", result["data"])
}

func TestAssertionStartEvent_JSON(t *testing.T) {
	event := NewAssertionStartEvent("run-123", "basic_http/login", 0, "assert_equals 0 exit_code")

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "assertion_start", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http/login", result["path"])
	assert.Equal(t, float64(0), result["index"])
	assert.Equal(t, "assert_equals 0 exit_code", result["command"])
}

func TestAssertionEndEvent_JSON(t *testing.T) {
	event := NewAssertionEndEvent("run-123", "basic_http/login", 0, 0)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "assertion_end", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http/login", result["path"])
	assert.Equal(t, float64(0), result["index"])
	assert.Equal(t, float64(0), result["exit_code"])
}

func TestTimeoutEvent_JSON(t *testing.T) {
	event := NewTimeoutEvent("run-123", "basic_http/slow", "run", "30s")

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "timeout", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http/slow", result["path"])
	assert.Equal(t, "run", result["phase"])
	assert.Equal(t, "30s", result["limit"])
}

func TestRunEndEvent_JSON(t *testing.T) {
	timestamp := time.Date(2026, 1, 15, 14, 45, 0, 0, time.UTC)

	event := NewRunEndEvent("2026-01-15_143022", "fail", 12, 2, timestamp)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "run_end", result["event"])
	assert.Equal(t, "2026-01-15_143022", result["run_id"])
	assert.Equal(t, "fail", result["status"])
	assert.Equal(t, float64(12), result["passed"])
	assert.Equal(t, float64(2), result["failed"])
	assert.Equal(t, "2026-01-15T14:45:00Z", result["timestamp"])
}

func TestScenarioRunStartEvent_JSON(t *testing.T) {
	event := NewScenarioRunStartEvent("run-123", "basic_http/login")

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "run_start", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http/login", result["path"])
}

func TestScenarioRunEndEvent_JSON(t *testing.T) {
	event := NewScenarioRunEndEvent("run-123", "basic_http/login", 0)

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "run_end", result["event"])
	assert.Equal(t, "run-123", result["run_id"])
	assert.Equal(t, "basic_http/login", result["path"])
	assert.Equal(t, float64(0), result["exit_code"])
}
