# Basanos

*βάσανος — the touchstone*

An acceptance test framework for agentic orchestration.

---

## Vision

A hierarchical acceptance testing framework that:
1. Defines system behavior as a tree of executable specifications
2. Runs deterministically against live systems
3. Emits structured events that can drive agentic TDD workflows
4. Provides full observability into every execution

The name **Basanos** (βάσανος) — the ancient Greek touchstone used to test if gold is genuine. The system under test is the metal; the acceptance suite is the basanos.

---

## Core Architecture

### Two Layers

**Framework (static, built once):**
- Spec schema and directory conventions
- Event-emitting test runner
- Pluggable sinks for output
- Assertion executables
- Reporters that transform captured data

**Agent-generated (per project):**
- Spec content (directory structure + YAML)
- Fixture files
- Helper scripts
- Iteration as requirements evolve

Key principle: **Spec generation is agentic; execution is deterministic.** No agent involvement at runtime.

---

## Spec Structure

The directory tree *is* the spec tree:

```
spec/
  context.yaml
  cli_commands/
    context.yaml
    expected_help.fixture
  basic_http/
    context.yaml
    expected_index.fixture
    user_sessions/
      context.yaml
      ...
```

### context.yaml Schema

```yaml
name: "Human-readable name"
description: "What this context tests"

env:
  PORT: 7654
  TESTROOT: ./testroot

on_failure: skip_children        # skip_children | continue | abort_run

before:
  run: |
    ${SERVER_CMD} -p ${PORT} &
    echo $! > ${CONTEXT_OUTPUT}/server.pid
  timeout: 10s

after:
  run: kill $(cat ${CONTEXT_OUTPUT}/server.pid)
  timeout: 5s

before_each:
  run: rm -rf ${TESTROOT}/state/*
  timeout: 2s

after_each:
  run: echo "cleanup"
  timeout: 1s

scenarios:
  - id: unique_identifier
    name: "Scenario name"
    
    before:
      run: curl -s http://localhost:${PORT}/setup
      timeout: 5s
    
    run:
      command: curl -s http://localhost:${PORT}/
      timeout: 30s
    
    assertions:
      - command: assert_contains expected_index.fixture ${SCENARIO_OUTPUT}/stdout
        timeout: 1s
      - command: assert_equals 0 ${SCENARIO_OUTPUT}/exit_code
    
    after:
      run: curl -s http://localhost:${PORT}/cleanup
      timeout: 5s

  - id: nested_group
    name: "Group of scenarios"
    
    env:
      EXTRA_VAR: "merged with parent env"
    
    on_failure: continue
    
    before_each:
      run: reset_state.sh
      timeout: 3s
    
    scenarios:
      - id: leaf_scenario
        run:
          command: some_command
          timeout: 10s
        assertions:
          - command: assert_equals 0 ${SCENARIO_OUTPUT}/exit_code
```

---

## Lifecycle Hooks

| Hook | Applies to | When |
|------|-----------|------|
| `before` | Groups or leaves | Once when entering this node |
| `after` | Groups or leaves | Once when exiting this node |
| `before_each` | Groups only | Before each descendant leaf |
| `after_each` | Groups only | After each descendant leaf |

**Execution order for a leaf:**
1. Ancestor `before`s (root → leaf)
2. Ancestor `before_each`s (root → leaf)
3. Self's `before`
4. Self's `run`
5. Self's `assertions`
6. Self's `after`
7. Ancestor `after_each`s (leaf → root)
8. Ancestor `after`s (leaf → root)

---

## Inheritance

| Property | Behavior |
|----------|----------|
| `env` | Merges down; child values override parent for same key |
| `on_failure` | Inherits from parent unless overridden |
| `timeout` | Per-command, explicit on each, no inheritance |

---

## Failure Modes

Configured via `on_failure`:

| Mode | Behavior |
|------|----------|
| `skip_children` | Skip remaining scenarios in this context, continue siblings |
| `continue` | Log failure, continue executing |
| `abort_run` | Stop entire test run immediately |

---

## Variables

| Variable | Scope | Meaning |
|----------|-------|---------|
| `${CONTEXT_OUTPUT}` | Any hook | Output dir for this context |
| `${SCENARIO_OUTPUT}` | Leaf scenario | Output dir for this scenario |
| `${SPEC_ROOT}` | Any | Root of spec/ directory |
| User `env:` vars | Inherited + merged | Custom variables |

---

## Event-Driven Execution

The runner emits a stream of events. Sinks subscribe and handle them. Execution is sequential — no parallelization.

```
┌─────────┐     events      ┌─────────────┐
│ Runner  │────────────────▶│ Sink: files │──▶ runs/
└─────────┘        │        └─────────────┘
                   │
                   │        ┌──────────────────┐
                   ├───────▶│ Sink: json-stream│──▶ stdout (ndjson)
                   │        └──────────────────┘
                   │
                   │        ┌──────────────────┐
                   ├───────▶│ Sink: websocket  │──▶ ws://...
                   │        └──────────────────┘
                   │
                   │        ┌──────────────────┐
                   └───────▶│ Sink: exec       │──▶ ./custom-sink.sh
                            └──────────────────┘
```

---

## Event Schema

```json
{"event": "run_start", "timestamp": "...", "run_id": "2026-01-15_143022"}

{"event": "context_enter", "path": "basic_http", "name": "Basic HTTP"}

{"event": "hook_start", "path": "basic_http", "hook": "before"}
{"event": "output", "stream": "stdout", "data": "Server starting...\n"}
{"event": "hook_end", "path": "basic_http", "hook": "before", "exit_code": 0}

{"event": "scenario_enter", "path": "basic_http/login", "name": "Login works"}

{"event": "hook_start", "path": "basic_http/login", "hook": "before_each", "from": "basic_http"}
{"event": "output", "stream": "stdout", "data": "..."}
{"event": "hook_end", "path": "basic_http/login", "hook": "before_each", "from": "basic_http", "exit_code": 0}

{"event": "run_start", "path": "basic_http/login"}
{"event": "output", "stream": "stdout", "data": "..."}
{"event": "run_end", "path": "basic_http/login", "exit_code": 0}

{"event": "assertion_start", "path": "basic_http/login", "index": 0, "command": "assert_equals ..."}
{"event": "output", "stream": "stdout", "data": "PASS: expected 200, got 200"}
{"event": "assertion_end", "path": "basic_http/login", "index": 0, "exit_code": 0}

{"event": "scenario_exit", "path": "basic_http/login", "status": "pass"}

{"event": "context_exit", "path": "basic_http"}

{"event": "timeout", "path": "basic_http/slow_scenario", "phase": "run", "limit": "30s"}

{"event": "run_end", "run_id": "2026-01-15_143022", "status": "fail", "passed": 12, "failed": 2}
```

---

## Output Structure: runs/

The `files` sink materializes events to disk:

```
runs/
  2026-01-15_143022/
    basic_http/
      before/
        stdout
        stderr
        exit_code
      
      user_sessions/
        before/
          stdout
          stderr
          exit_code
        
        login/
          before_each/
            basic_http/
              stdout
              stderr
              exit_code
            user_sessions/
              stdout
              stderr
              exit_code
          before/
            stdout
            stderr
            exit_code
          run/
            stdout
            stderr
            exit_code
          assertions/
            0_assert_equals/
              stdout
              stderr
              exit_code
            1_assert_contains/
              stdout
              stderr
              exit_code
          after/
            stdout
            stderr
            exit_code
        
        after/
          stdout
          stderr
          exit_code
      
      after/
        stdout
        stderr
        exit_code
```

---

## CLI Usage

```bash
# Run specs
basanos run
basanos run --spec ./my-specs

# Sinks
basanos run --sink files
basanos run --sink json-stream
basanos run --sink ws:localhost:9999
basanos run --sink exec:./my-handler.sh
basanos run --sink files --sink json-stream

# Filter
basanos run --filter basic_http/login
basanos run --filter "basic_http/*"

# Reporting
basanos report json
basanos report junit
basanos report summary

# MCP (future)
basanos serve --mcp
```

---

## Assertion Executables

Standalone binaries, distributed alongside `basanos`:

```bash
assert_equals expected.txt actual.txt
assert_contains needle.txt haystack.txt
assert_matches "pattern" target.txt
assert_gt 10 5
assert_gte 10 10
assert_lt 5 10
assert_lte 10 10
```

Each:
- Exits 0 on pass, non-zero on fail
- Outputs structured diff/comparison info to stdout
- Can be used independently of basanos

Referenced in specs:

```yaml
assertions:
  - command: assert_equals expected.fixture ${SCENARIO_OUTPUT}/stdout
    timeout: 1s
```

---

## Implementation

**Language:** Go

**Rationale:**
- Single static binary
- Trivial cross-compilation
- Fast startup
- No runtime dependencies

**Installation:**
```bash
curl -L https://github.com/.../basanos-$(uname -s)-$(uname -m) -o /usr/local/bin/basanos
chmod +x /usr/local/bin/basanos
```

---

## Future: MCP Integration

```
┌─────────────────────────────────────────────────────────────┐
│   Agent Orchestrator                                        │
│   ┌─────────────────────────────────────────────────────┐   │
│   │ MCP Client                                          │   │
│   │                                                     │   │
│   │ tools:                                              │   │
│   │   - basanos/run                                     │   │
│   │   - basanos/status                                  │   │
│   │   - basanos/failures                                │   │
│   │   - basanos/subscribe                               │   │
│   │                                                     │   │
│   │ resources:                                          │   │
│   │   - spec://basic_http/login                         │   │
│   │   - runs://latest/...                               │   │
│   │                                                     │   │
│   └─────────────────────────────────────────────────────┘   │
│                           │                                 │
│                           ▼                                 │
│                    ┌─────────────┐                          │
│                    │   Basanos   │                          │
│                    │  MCP Server │                          │
│                    └─────────────┘                          │
└─────────────────────────────────────────────────────────────┘
```

The agent can:
- Trigger test runs
- Stream events in real-time
- Query failures with full context
- Read specs and outputs
- React by delegating TDD work
- Re-run to verify fixes

---

## Inspirations

- **FitNesse** — Acceptance testing from the Clean Code tradition
- **InterUSS uss_qualifier** — Hierarchical compliance testing
- **http-spec** — Different execution contexts for different test categories
- **RSpec/speclj** — Nested contexts with lifecycle hooks
