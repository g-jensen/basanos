package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_ValidSpec_ReturnsEmptySlice(t *testing.T) {
	ctx := &Context{Name: "Valid Spec"}

	errors := Validate(ctx, "context.yaml")

	assert.Equal(t, []ValidationError{}, errors)
}

func TestValidate_ScenarioWithoutID_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name:      "Test Spec",
		Scenarios: []Scenario{{Name: "Missing ID Scenario"}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].id", errors[0].Path)
	assert.Contains(t, errors[0].Message, "required")
}

func TestValidate_InvalidOnFailure_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name:      "Test Spec",
		OnFailure: "invalid_value",
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "on_failure", errors[0].Path)
	assert.Contains(t, errors[0].Message, "skip_children")
	assert.Contains(t, errors[0].Message, "continue")
	assert.Contains(t, errors[0].Message, "abort_run")
}

func TestValidate_ScenarioInvalidOnFailure_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name:      "Test Spec",
		Scenarios: []Scenario{{ID: "test", OnFailure: "bad_value"}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].on_failure", errors[0].Path)
}

func TestValidate_BeforeHookWithoutRun_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name:   "Test Spec",
		Before: &Hook{Timeout: "5s"},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "before.run", errors[0].Path)
	assert.Contains(t, errors[0].Message, "required")
}

func TestValidate_ScenarioRunBlockWithoutCommand_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name:      "Test Spec",
		Scenarios: []Scenario{{ID: "test", Run: &RunBlock{Timeout: "10s"}}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].run.command", errors[0].Path)
	assert.Contains(t, errors[0].Message, "required")
}

func TestValidate_AssertionWithoutCommand_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name: "Test Spec",
		Scenarios: []Scenario{{
			ID:         "test",
			Run:        &RunBlock{Command: "echo hello"},
			Assertions: []Assertion{{Timeout: "1s"}},
		}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].assertions[0].command", errors[0].Path)
	assert.Contains(t, errors[0].Message, "required")
}

func TestValidate_HookInvalidTimeout_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name:   "Test Spec",
		Before: &Hook{Run: "echo test", Timeout: "invalid"},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "before.timeout", errors[0].Path)
	assert.Contains(t, errors[0].Message, "duration")
}

func TestValidate_RunBlockInvalidTimeout_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name:      "Test Spec",
		Scenarios: []Scenario{{ID: "test", Run: &RunBlock{Command: "echo test", Timeout: "bad"}}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].run.timeout", errors[0].Path)
	assert.Contains(t, errors[0].Message, "duration")
}

func TestValidate_AssertionInvalidTimeout_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name: "Test Spec",
		Scenarios: []Scenario{{
			ID:         "test",
			Run:        &RunBlock{Command: "echo hello"},
			Assertions: []Assertion{{Command: "assert_equals", Timeout: "xyz"}},
		}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].assertions[0].timeout", errors[0].Path)
	assert.Contains(t, errors[0].Message, "duration")
}

func TestValidate_DuplicateScenarioIDs_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name: "Test Spec",
		Scenarios: []Scenario{
			{ID: "duplicate", Run: &RunBlock{Command: "echo first"}},
			{ID: "duplicate", Run: &RunBlock{Command: "echo second"}},
		},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[1].id", errors[0].Path)
	assert.Contains(t, errors[0].Message, "duplicate")
}

func TestValidate_ScenarioWithRunAndNestedScenarios_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name: "Test Spec",
		Scenarios: []Scenario{{
			ID:  "hybrid",
			Run: &RunBlock{Command: "echo leaf"},
			Scenarios: []Scenario{
				{ID: "child", Run: &RunBlock{Command: "echo child"}},
			},
		}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].run", errors[0].Path)
	assert.Contains(t, errors[0].Message, "group")
}

func TestValidate_LeafScenarioWithBeforeEach_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name: "Test Spec",
		Scenarios: []Scenario{{
			ID:         "leaf",
			Run:        &RunBlock{Command: "echo leaf"},
			BeforeEach: &Hook{Run: "echo before each"},
		}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].before_each", errors[0].Path)
	assert.Contains(t, errors[0].Message, "leaf")
}

func TestValidate_LeafScenarioWithAfterEach_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name: "Test Spec",
		Scenarios: []Scenario{{
			ID:        "leaf",
			Run:       &RunBlock{Command: "echo leaf"},
			AfterEach: &Hook{Run: "echo after each"},
		}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].after_each", errors[0].Path)
	assert.Contains(t, errors[0].Message, "leaf")
}

func TestValidate_AfterHookWithoutRun_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name:  "Test Spec",
		After: &Hook{Timeout: "5s"},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "after.run", errors[0].Path)
	assert.Contains(t, errors[0].Message, "required")
}

func TestValidate_NestedScenarioWithoutID_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name: "Test Spec",
		Scenarios: []Scenario{{
			ID:        "group",
			Scenarios: []Scenario{{Name: "Missing ID"}},
		}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].scenarios[0].id", errors[0].Path)
	assert.Contains(t, errors[0].Message, "required")
}

func TestValidate_ContextBeforeEachWithoutRun_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name:       "Test Spec",
		BeforeEach: &Hook{Timeout: "2s"},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "before_each.run", errors[0].Path)
	assert.Contains(t, errors[0].Message, "required")
}

func TestValidate_ContextAfterEachWithoutRun_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name:      "Test Spec",
		AfterEach: &Hook{Timeout: "1s"},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "after_each.run", errors[0].Path)
	assert.Contains(t, errors[0].Message, "required")
}

func TestValidate_ScenarioBeforeHookWithoutRun_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name: "Test Spec",
		Scenarios: []Scenario{{
			ID:     "test",
			Run:    &RunBlock{Command: "echo hello"},
			Before: &Hook{Timeout: "5s"},
		}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].before.run", errors[0].Path)
	assert.Contains(t, errors[0].Message, "required")
}

func TestValidate_ScenarioAfterHookWithoutRun_ReturnsError(t *testing.T) {
	ctx := &Context{
		Name: "Test Spec",
		Scenarios: []Scenario{{
			ID:    "test",
			Run:   &RunBlock{Command: "echo hello"},
			After: &Hook{Timeout: "3s"},
		}},
	}

	errors := Validate(ctx, "context.yaml")

	require.Len(t, errors, 1)
	assert.Equal(t, "scenarios[0].after.run", errors[0].Path)
	assert.Contains(t, errors[0].Message, "required")
}
