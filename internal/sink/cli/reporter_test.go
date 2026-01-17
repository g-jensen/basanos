package cli

import (
	"bytes"
	"testing"
	"time"

	"basanos/internal/event"

	"github.com/stretchr/testify/assert"
)

func TestSink_PrintsSummaryOnRunEnd(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := NewReporter(buffer, false, false)

	timestamp := time.Date(2026, 1, 15, 14, 30, 22, 0, time.UTC)
	sink.Emit(event.NewRunEndEvent("run-1", "fail", 3, 1, timestamp))

	assert.Equal(t, "\n\n3 passed, 1 failed\n", buffer.String())
}

func TestSink_PrintsFailuresBeforeSummary(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := NewReporter(buffer, false, false)

	timestamp := time.Date(2026, 1, 15, 14, 30, 22, 0, time.UTC)
	sink.Emit(event.NewScenarioExitEvent("run-1", "basic_http/health", "pass", timestamp))
	sink.Emit(event.NewScenarioExitEvent("run-1", "basic_http/login", "fail", timestamp))
	sink.Emit(event.NewScenarioExitEvent("run-1", "basic_http/status", "pass", timestamp))
	sink.Emit(event.NewRunEndEvent("run-1", "fail", 2, 1, timestamp))

	expected := `.F.

Failures:

  1) basic_http/login

2 passed, 1 failed
`
	assert.Equal(t, expected, buffer.String())
}

func TestSink_DisplaysStdoutForFailedScenario(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := NewReporter(buffer, false, false)

	timestamp := time.Date(2026, 1, 15, 14, 30, 22, 0, time.UTC)
	sink.Emit(event.NewScenarioEnterEvent("run-1", "basic_http/health", "Health Check", timestamp))
	sink.Emit(event.NewScenarioExitEvent("run-1", "basic_http/health", "pass", timestamp))
	sink.Emit(event.NewScenarioEnterEvent("run-1", "basic_http/login", "Login", timestamp))
	sink.Emit(event.NewOutputEvent("run-1", "stdout", "Login failed\n"))
	sink.Emit(event.NewScenarioExitEvent("run-1", "basic_http/login", "fail", timestamp))
	sink.Emit(event.NewRunEndEvent("run-1", "fail", 1, 1, timestamp))

	expected := `.F

Failures:

  1) basic_http/login
     stdout:
       Login failed

1 passed, 1 failed
`
	assert.Equal(t, expected, buffer.String())
}

func TestSink_DisplaysStderrForFailedScenario(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := NewReporter(buffer, false, false)

	timestamp := time.Date(2026, 1, 15, 14, 30, 22, 0, time.UTC)
	sink.Emit(event.NewScenarioEnterEvent("run-1", "basic_http/error", "Error Test", timestamp))
	sink.Emit(event.NewOutputEvent("run-1", "stdout", "Attempting request\n"))
	sink.Emit(event.NewOutputEvent("run-1", "stderr", "Connection refused\n"))
	sink.Emit(event.NewScenarioExitEvent("run-1", "basic_http/error", "fail", timestamp))
	sink.Emit(event.NewRunEndEvent("run-1", "fail", 0, 1, timestamp))

	expected := `F

Failures:

  1) basic_http/error
     stdout:
       Attempting request
     stderr:
       Connection refused

0 passed, 1 failed
`
	assert.Equal(t, expected, buffer.String())
}

func TestSink_FullVerboseRun(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := NewReporter(buffer, true, false)

	timestamp := time.Date(2026, 1, 15, 14, 30, 22, 0, time.UTC)
	sink.Emit(event.NewContextEnterEvent("run-1", "basic_http", "Basic HTTP", timestamp))
	sink.Emit(event.NewContextEnterEvent("run-1", "basic_http/user_sessions", "User Sessions", timestamp))
	sink.Emit(event.NewScenarioEnterEvent("run-1", "basic_http/user_sessions/login", "Login works", timestamp))
	sink.Emit(event.NewScenarioExitEvent("run-1", "basic_http/user_sessions/login", "pass", timestamp))
	sink.Emit(event.NewContextExitEvent("run-1", "basic_http/user_sessions", timestamp))
	sink.Emit(event.NewContextExitEvent("run-1", "basic_http", timestamp))
	sink.Emit(event.NewRunEndEvent("run-1", "pass", 1, 0, timestamp))

	expected := `Basic HTTP
  User Sessions
    Login works .

1 passed, 0 failed
`
	assert.Equal(t, expected, buffer.String())
}

func TestSink_FullColorRun(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := NewReporter(buffer, true, true)

	timestamp := time.Date(2026, 1, 15, 14, 30, 22, 0, time.UTC)
	sink.Emit(event.NewContextEnterEvent("run-1", "parent", "Parent", timestamp))
	sink.Emit(event.NewScenarioEnterEvent("run-1", "parent/pass", "Passes", timestamp))
	sink.Emit(event.NewScenarioExitEvent("run-1", "parent/pass", "pass", timestamp))
	sink.Emit(event.NewScenarioEnterEvent("run-1", "parent/fail", "Fails", timestamp))
	sink.Emit(event.NewScenarioExitEvent("run-1", "parent/fail", "fail", timestamp))
	sink.Emit(event.NewContextExitEvent("run-1", "parent", timestamp))
	sink.Emit(event.NewRunEndEvent("run-1", "fail", 1, 1, timestamp))

	assert.Contains(t, buffer.String(), "\033[32mPasses\033[0m")
	assert.Contains(t, buffer.String(), "\033[31mFails\033[0m")
}
