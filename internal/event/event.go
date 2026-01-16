package event

import "time"

type RunStartEvent struct {
	Event     string    `json:"event"`
	RunID     string    `json:"run_id"`
	Timestamp time.Time `json:"timestamp"`
}

func NewRunStartEvent(runID string, timestamp time.Time) *RunStartEvent {
	return &RunStartEvent{
		Event:     "run_start",
		RunID:     runID,
		Timestamp: timestamp,
	}
}

type ContextEnterEvent struct {
	Event     string    `json:"event"`
	RunID     string    `json:"run_id,omitempty"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func NewContextEnterEvent(runID, path, name string, timestamp time.Time) *ContextEnterEvent {
	return &ContextEnterEvent{
		Event:     "context_enter",
		RunID:     runID,
		Path:      path,
		Name:      name,
		Timestamp: timestamp,
	}
}

type ContextExitEvent struct {
	Event     string    `json:"event"`
	RunID     string    `json:"run_id,omitempty"`
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
}

func NewContextExitEvent(runID, path string, timestamp time.Time) *ContextExitEvent {
	return &ContextExitEvent{
		Event:     "context_exit",
		RunID:     runID,
		Path:      path,
		Timestamp: timestamp,
	}
}

type HookStartEvent struct {
	Event string `json:"event"`
	RunID string `json:"run_id,omitempty"`
	Path  string `json:"path"`
	Hook  string `json:"hook"`
	From  string `json:"from,omitempty"`
}

func NewHookStartEvent(runID, path, hook, from string) *HookStartEvent {
	return &HookStartEvent{
		Event: "hook_start",
		RunID: runID,
		Path:  path,
		Hook:  hook,
		From:  from,
	}
}

type HookEndEvent struct {
	Event    string `json:"event"`
	RunID    string `json:"run_id,omitempty"`
	Path     string `json:"path"`
	Hook     string `json:"hook"`
	From     string `json:"from,omitempty"`
	ExitCode int    `json:"exit_code"`
}

func NewHookEndEvent(runID, path, hook, from string, exitCode int) *HookEndEvent {
	return &HookEndEvent{
		Event:    "hook_end",
		RunID:    runID,
		Path:     path,
		Hook:     hook,
		From:     from,
		ExitCode: exitCode,
	}
}

type ScenarioEnterEvent struct {
	Event     string    `json:"event"`
	RunID     string    `json:"run_id,omitempty"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func NewScenarioEnterEvent(runID, path, name string, timestamp time.Time) *ScenarioEnterEvent {
	return &ScenarioEnterEvent{
		Event:     "scenario_enter",
		RunID:     runID,
		Path:      path,
		Name:      name,
		Timestamp: timestamp,
	}
}

type ScenarioExitEvent struct {
	Event     string    `json:"event"`
	RunID     string    `json:"run_id,omitempty"`
	Path      string    `json:"path"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func NewScenarioExitEvent(runID, path, status string, timestamp time.Time) *ScenarioExitEvent {
	return &ScenarioExitEvent{
		Event:     "scenario_exit",
		RunID:     runID,
		Path:      path,
		Status:    status,
		Timestamp: timestamp,
	}
}

type ScenarioRunStartEvent struct {
	Event string `json:"event"`
	RunID string `json:"run_id,omitempty"`
	Path  string `json:"path"`
}

func NewScenarioRunStartEvent(runID, path string) *ScenarioRunStartEvent {
	return &ScenarioRunStartEvent{
		Event: "run_start",
		RunID: runID,
		Path:  path,
	}
}

type ScenarioRunEndEvent struct {
	Event    string `json:"event"`
	RunID    string `json:"run_id,omitempty"`
	Path     string `json:"path"`
	ExitCode int    `json:"exit_code"`
}

func NewScenarioRunEndEvent(runID, path string, exitCode int) *ScenarioRunEndEvent {
	return &ScenarioRunEndEvent{
		Event:    "run_end",
		RunID:    runID,
		Path:     path,
		ExitCode: exitCode,
	}
}

type OutputEvent struct {
	Event  string `json:"event"`
	RunID  string `json:"run_id,omitempty"`
	Stream string `json:"stream"`
	Data   string `json:"data"`
}

func NewOutputEvent(runID, stream, data string) *OutputEvent {
	return &OutputEvent{
		Event:  "output",
		RunID:  runID,
		Stream: stream,
		Data:   data,
	}
}

type AssertionStartEvent struct {
	Event   string `json:"event"`
	RunID   string `json:"run_id,omitempty"`
	Path    string `json:"path"`
	Index   int    `json:"index"`
	Command string `json:"command"`
}

func NewAssertionStartEvent(runID, path string, index int, command string) *AssertionStartEvent {
	return &AssertionStartEvent{
		Event:   "assertion_start",
		RunID:   runID,
		Path:    path,
		Index:   index,
		Command: command,
	}
}

type AssertionEndEvent struct {
	Event    string `json:"event"`
	RunID    string `json:"run_id,omitempty"`
	Path     string `json:"path"`
	Index    int    `json:"index"`
	ExitCode int    `json:"exit_code"`
}

func NewAssertionEndEvent(runID, path string, index int, exitCode int) *AssertionEndEvent {
	return &AssertionEndEvent{
		Event:    "assertion_end",
		RunID:    runID,
		Path:     path,
		Index:    index,
		ExitCode: exitCode,
	}
}

type TimeoutEvent struct {
	Event string `json:"event"`
	RunID string `json:"run_id,omitempty"`
	Path  string `json:"path"`
	Phase string `json:"phase"`
	Limit string `json:"limit"`
}

func NewTimeoutEvent(runID, path, phase, limit string) *TimeoutEvent {
	return &TimeoutEvent{
		Event: "timeout",
		RunID: runID,
		Path:  path,
		Phase: phase,
		Limit: limit,
	}
}

type RunEndEvent struct {
	Event     string    `json:"event"`
	RunID     string    `json:"run_id"`
	Status    string    `json:"status"`
	Passed    int       `json:"passed"`
	Failed    int       `json:"failed"`
	Timestamp time.Time `json:"timestamp"`
}

func NewRunEndEvent(runID, status string, passed, failed int, timestamp time.Time) *RunEndEvent {
	return &RunEndEvent{
		Event:     "run_end",
		RunID:     runID,
		Status:    status,
		Passed:    passed,
		Failed:    failed,
		Timestamp: timestamp,
	}
}
