---
name: basanos
description: Basanos acceptance test spec format, lifecycle hooks, variables, and assertion executables
---

# Basanos Spec Format

Reference for generating Basanos acceptance test specifications.

## Directory Structure

The directory tree *is* the spec tree:

```
spec/
  context.yaml              # Root context
  authentication/
    context.yaml            # Login, logout, password reset
    expected_login.fixture
  users/
    context.yaml            # User CRUD, profiles
    registration/
      context.yaml          # Registration flows specifically
  orders/
    context.yaml            # Order lifecycle
    checkout/
      context.yaml          # Checkout flow
      expected_receipt.fixture
  inventory/
    context.yaml            # Stock management
```

Organize by **domain**, not technical layers. The structure should scream what the system does—someone glancing at the directory tree should understand the business, not the architecture.

**Good:** `authentication/`, `users/`, `orders/`, `inventory/`

**Bad:** `api/`, `ui/`, `controllers/`, `services/`

## context.yaml Schema

```yaml
# Required
name: "Human-readable name"

# Optional
description: "What this context tests"

# Environment variables - inherited and merged down the tree
# Child values override parent for same key
env:
  PORT: "8080"
  API_URL: "http://localhost:${PORT}"

# Failure handling (inherited unless overridden)
# Options: skip_children | continue | abort_run
on_failure: skip_children

# Lifecycle hooks (all optional)
before:
  run: ./start-server.sh
  timeout: 10s

after:
  run: ./stop-server.sh
  timeout: 5s

before_each:
  run: ./reset-database.sh
  timeout: 3s

after_each:
  run: ./cleanup.sh
  timeout: 2s

# Test scenarios
scenarios:
  # Leaf scenario (has 'run')
  - id: unique_snake_case_id
    name: "Human readable name"
    
    before:
      run: ./setup-test-data.sh
      timeout: 5s
    
    run:
      command: curl -s http://localhost:${PORT}/endpoint
      timeout: 30s
    
    assertions:
      - command: assert_equals expected.fixture ${RUN_OUTPUT}/stdout
        timeout: 1s
      - command: assert_equals 0 ${RUN_OUTPUT}/exit_code
        timeout: 1s
    
    after:
      run: ./cleanup-test-data.sh
      timeout: 5s

  # Group scenario (has nested 'scenarios', no 'run')
  - id: user_management
    name: "User Management"
    
    env:
      USER_API: "${API_URL}/users"
    
    before_each:
      run: ./reset-users.sh
      timeout: 3s
    
    scenarios:
      - id: create_user
        name: "Create user returns 201"
        run:
          command: curl -X POST ${USER_API}
          timeout: 10s
        assertions:
          - command: assert_contains "201" ${RUN_OUTPUT}/stdout
```

## Lifecycle Hook Execution Order

For a leaf scenario, hooks execute in this order:

1. Ancestor `before` hooks (root to leaf, run once per context)
2. Ancestor `before_each` hooks (root to leaf)
3. Scenario's own `before` hook
4. Scenario's `run` command
5. Scenario's `assertions` (in order)
6. Scenario's own `after` hook
7. Ancestor `after_each` hooks (leaf to root)
8. Ancestor `after` hooks (leaf to root, run once per context)

## Hook Applicability

| Hook | Applies To | When It Runs |
|------|------------|--------------|
| `before` | Contexts, groups, leaves | Once when entering this node |
| `after` | Contexts, groups, leaves | Once when exiting this node |
| `before_each` | Contexts, groups only | Before each descendant leaf |
| `after_each` | Contexts, groups only | After each descendant leaf |

## Variables

| Variable | Scope | Description |
|----------|-------|-------------|
| `${SPEC_ROOT}` | All | Absolute path to spec directory root |
| `${CONTEXT_OUTPUT}` | Context hooks | Output directory for current context |
| `${SCENARIO_OUTPUT}` | Scenario | Output directory for current scenario |
| `${RUN_OUTPUT}` | Scenario | Shorthand for `${SCENARIO_OUTPUT}/_run` |
| Custom `env` vars | Inherited | Merged down tree, child overrides parent |

## Assertion Executables

All assertions exit 0 on pass, non-zero on fail. They auto-detect file paths vs literal values.

### Equality
```yaml
# Compare file to file
- command: assert_equals expected.fixture ${RUN_OUTPUT}/stdout

# Compare literal to file
- command: assert_equals "200" ${RUN_OUTPUT}/stdout

# Check exit code
- command: assert_equals 0 ${RUN_OUTPUT}/exit_code
```

### Containment
```yaml
# File contains content from another file
- command: assert_contains needle.fixture ${RUN_OUTPUT}/stdout

# File contains literal string
- command: assert_contains "success" ${RUN_OUTPUT}/stdout
```

### Pattern Matching
```yaml
# File matches regex
- command: assert_matches "user_[0-9]+" ${RUN_OUTPUT}/stdout
```

### Numeric Comparisons
```yaml
# Greater than
- command: assert_gt ${RUN_OUTPUT}/stdout 0

# Greater than or equal
- command: assert_gte ${RUN_OUTPUT}/stdout 100

# Less than
- command: assert_lt ${RUN_OUTPUT}/stdout 1000

# Less than or equal
- command: assert_lte ${RUN_OUTPUT}/stdout 500
```

## Failure Modes

| Mode | Behavior |
|------|----------|
| `skip_children` | Skip remaining scenarios in this context, continue with sibling contexts |
| `continue` | Log failure, continue executing all scenarios |
| `abort_run` | Stop entire test run immediately |

## Fixture Files

For non-trivial expected outputs, create fixture files alongside specs:

```
spec/
  orders/
    checkout/
      context.yaml
      expected_receipt.fixture
      expected_confirmation_email.fixture
```

Reference in assertions:
```yaml
assertions:
  - command: assert_equals expected_receipt.fixture ${RUN_OUTPUT}/stdout
```

## Common Patterns

### Testing HTTP Endpoints
```yaml
scenarios:
  - id: list_users
    name: "GET /users returns user list"
    run:
      command: curl -s http://localhost:${PORT}/users
      timeout: 10s
    assertions:
      - command: assert_contains "alice" ${RUN_OUTPUT}/stdout
      - command: assert_equals 0 ${RUN_OUTPUT}/exit_code
```

### Testing CLI Commands
```yaml
scenarios:
  - id: help_flag
    name: "--help prints usage"
    run:
      command: myapp --help
      timeout: 5s
    assertions:
      - command: assert_contains "Usage:" ${RUN_OUTPUT}/stdout
      - command: assert_equals 0 ${RUN_OUTPUT}/exit_code
```

### Testing Error Conditions
```yaml
scenarios:
  - id: invalid_input
    name: "Invalid input returns error"
    run:
      command: myapp --invalid-flag
      timeout: 5s
    assertions:
      - command: assert_contains "unknown flag" ${RUN_OUTPUT}/stderr
      - command: assert_gt ${RUN_OUTPUT}/exit_code 0
```

### Setup/Teardown with Persistent State
```yaml
before:
  run: |
    ./start-server.sh &
    echo $! > ${CONTEXT_OUTPUT}/server.pid
  timeout: 10s

after:
  run: kill $(cat ${CONTEXT_OUTPUT}/server.pid)
  timeout: 5s
```

### Resetting State Between Tests
```yaml
before_each:
  run: ./reset-database.sh
  timeout: 3s
```

## Anti-Patterns

### Testing Implementation Details
```yaml
# BAD - peeking into database
assertions:
  - command: assert_equals "3" $(psql -c "SELECT COUNT(*) FROM users")

# GOOD - testing observable behavior
assertions:
  - command: assert_contains '"count": 3' ${RUN_OUTPUT}/stdout
```

### Overly Loose Assertions
```yaml
# BAD - only checks exit code, ignores actual output
assertions:
  - command: assert_equals 0 ${RUN_OUTPUT}/exit_code

# GOOD - verifies actual content
assertions:
  - command: assert_equals expected_response.fixture ${RUN_OUTPUT}/stdout
  - command: assert_equals 0 ${RUN_OUTPUT}/exit_code
```

### Hardcoded Environment Values
```yaml
# BAD - hardcoded port
run:
  command: curl http://localhost:8080/users

# GOOD - uses environment variable
env:
  API_PORT: "8080"
run:
  command: curl http://localhost:${API_PORT}/users
```

### Missing Timeouts
```yaml
# BAD - no timeout, could hang forever
run:
  command: curl http://localhost:${PORT}/slow-endpoint

# GOOD - explicit timeout
run:
  command: curl http://localhost:${PORT}/slow-endpoint
  timeout: 30s
```

### Overly Tight Timeouts
```yaml
# BAD - too tight, will flake
run:
  command: ./database-migration.sh
  timeout: 1s

# GOOD - reasonable buffer
run:
  command: ./database-migration.sh
  timeout: 60s
```
