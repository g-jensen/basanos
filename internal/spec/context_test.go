package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseContext_Name(t *testing.T) {
	ctx, err := ParseContext([]byte(`name: "Basic HTTP Tests"`))

	require.NoError(t, err)
	assert.Equal(t, "Basic HTTP Tests", ctx.Name)
}

func TestParseContext_Description(t *testing.T) {
	ctx, err := ParseContext([]byte(`description: "Tests for user authentication"`))

	require.NoError(t, err)
	assert.Equal(t, "Tests for user authentication", ctx.Description)
}

func TestParseContext_Env(t *testing.T) {
	yaml := `
env:
  PORT: "7654"
  TESTROOT: ./testroot
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Equal(t, "7654", ctx.Env["PORT"])
	assert.Equal(t, "./testroot", ctx.Env["TESTROOT"])
}

func TestParseContext_OnFailure(t *testing.T) {
	ctx, err := ParseContext([]byte(`on_failure: skip_children`))

	require.NoError(t, err)
	assert.Equal(t, "skip_children", ctx.OnFailure)
}

func TestParseContext_BeforeHook(t *testing.T) {
	yaml := `
before:
  run: echo "starting server"
  timeout: 10s
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Equal(t, `echo "starting server"`, ctx.Before.Run)
	assert.Equal(t, "10s", ctx.Before.Timeout)
}

func TestParseContext_AfterHook(t *testing.T) {
	yaml := `
after:
  run: kill $(cat server.pid)
  timeout: 5s
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Equal(t, "kill $(cat server.pid)", ctx.After.Run)
	assert.Equal(t, "5s", ctx.After.Timeout)
}

func TestParseContext_BeforeEachHook(t *testing.T) {
	yaml := `
before_each:
  run: rm -rf ${TESTROOT}/state/*
  timeout: 2s
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Equal(t, "rm -rf ${TESTROOT}/state/*", ctx.BeforeEach.Run)
	assert.Equal(t, "2s", ctx.BeforeEach.Timeout)
}

func TestParseContext_AfterEachHook(t *testing.T) {
	yaml := `
after_each:
  run: echo "cleanup"
  timeout: 1s
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Equal(t, `echo "cleanup"`, ctx.AfterEach.Run)
	assert.Equal(t, "1s", ctx.AfterEach.Timeout)
}

func TestParseContext_Scenarios(t *testing.T) {
	yaml := `
scenarios:
  - id: login_works
    name: "Login works"
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Len(t, ctx.Scenarios, 1)
	assert.Equal(t, "login_works", ctx.Scenarios[0].ID)
	assert.Equal(t, "Login works", ctx.Scenarios[0].Name)
}

func TestParseContext_ScenarioRun(t *testing.T) {
	yaml := `
scenarios:
  - id: fetch_index
    name: "Fetch index"
    run:
      command: curl -s http://localhost:8080/
      timeout: 30s
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Equal(t, "curl -s http://localhost:8080/", ctx.Scenarios[0].Run.Command)
	assert.Equal(t, "30s", ctx.Scenarios[0].Run.Timeout)
}

func TestParseContext_ScenarioAssertions(t *testing.T) {
	yaml := `
scenarios:
  - id: check_response
    name: "Check response"
    assertions:
      - command: assert_contains expected.fixture ${SCENARIO_OUTPUT}/stdout
        timeout: 1s
      - command: assert_equals 0 ${SCENARIO_OUTPUT}/exit_code
        timeout: 1s
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Len(t, ctx.Scenarios[0].Assertions, 2)
	assert.Equal(t, "assert_contains expected.fixture ${SCENARIO_OUTPUT}/stdout", ctx.Scenarios[0].Assertions[0].Command)
	assert.Equal(t, "1s", ctx.Scenarios[0].Assertions[0].Timeout)
}

func TestParseContext_ScenarioHooks(t *testing.T) {
	yaml := `
scenarios:
  - id: with_setup
    name: "With setup"
    before:
      run: curl -s http://localhost/setup
      timeout: 5s
    after:
      run: curl -s http://localhost/cleanup
      timeout: 3s
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Equal(t, "curl -s http://localhost/setup", ctx.Scenarios[0].Before.Run)
	assert.Equal(t, "5s", ctx.Scenarios[0].Before.Timeout)
	assert.Equal(t, "curl -s http://localhost/cleanup", ctx.Scenarios[0].After.Run)
	assert.Equal(t, "3s", ctx.Scenarios[0].After.Timeout)
}

func TestParseContext_NestedScenarios(t *testing.T) {
	yaml := `
scenarios:
  - id: user_sessions
    name: "User Sessions"
    scenarios:
      - id: login
        name: "Login works"
        run:
          command: curl http://localhost/login
          timeout: 10s
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Equal(t, "user_sessions", ctx.Scenarios[0].ID)
	assert.Equal(t, "login", ctx.Scenarios[0].Scenarios[0].ID)
	assert.Equal(t, "Login works", ctx.Scenarios[0].Scenarios[0].Name)
}

func TestParseContext_ScenarioGroupFields(t *testing.T) {
	yaml := `
scenarios:
  - id: user_sessions
    name: "User Sessions"
    env:
      SESSION_TIMEOUT: "3600"
    on_failure: continue
    before_each:
      run: reset_state.sh
      timeout: 3s
    after_each:
      run: cleanup.sh
      timeout: 2s
    scenarios:
      - id: login
        name: "Login"
`
	ctx, err := ParseContext([]byte(yaml))

	require.NoError(t, err)
	assert.Equal(t, "3600", ctx.Scenarios[0].Env["SESSION_TIMEOUT"])
	assert.Equal(t, "continue", ctx.Scenarios[0].OnFailure)
	assert.Equal(t, "reset_state.sh", ctx.Scenarios[0].BeforeEach.Run)
	assert.Equal(t, "cleanup.sh", ctx.Scenarios[0].AfterEach.Run)
}
