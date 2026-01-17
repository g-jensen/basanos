package spec

import (
	"fmt"
	"time"
)

var validOnFailure = map[string]bool{
	"":              true,
	"skip_children": true,
	"continue":      true,
	"abort_run":     true,
}

const invalidOnFailureMessage = "must be skip_children, continue, or abort_run"

type ValidationError struct {
	File    string
	Path    string
	Message string
}

type validator struct {
	file   string
	errors []ValidationError
}

func (validator *validator) addError(path, message string) {
	validator.errors = append(validator.errors, ValidationError{
		File:    validator.file,
		Path:    path,
		Message: message,
	})
}

func (validator *validator) checkOnFailure(value, path string) {
	if !validOnFailure[value] {
		validator.addError(path, invalidOnFailureMessage)
	}
}

func (validator *validator) checkTimeout(timeout, path string) {
	if timeout == "" {
		return
	}
	_, err := time.ParseDuration(timeout)
	if err != nil {
		validator.addError(path, "invalid duration")
	}
}

func (validator *validator) validateHook(hook *Hook, path string) {
	if hook == nil {
		return
	}
	if hook.Run == "" {
		validator.addError(path+".run", "required")
	}
	validator.checkTimeout(hook.Timeout, path+".timeout")
}

func (validator *validator) validateRunBlock(runBlock *RunBlock, path string) {
	if runBlock == nil {
		return
	}
	if runBlock.Command == "" {
		validator.addError(path+".command", "required")
	}
	validator.checkTimeout(runBlock.Timeout, path+".timeout")
}

func (validator *validator) validateAssertion(assertion Assertion, path string) {
	if assertion.Command == "" {
		validator.addError(path+".command", "required")
	}
	validator.checkTimeout(assertion.Timeout, path+".timeout")
}

func isLeaf(scenario Scenario) bool {
	return scenario.Run != nil && len(scenario.Scenarios) == 0
}

func isGroup(scenario Scenario) bool {
	return len(scenario.Scenarios) > 0
}

func (validator *validator) validateScenario(scenario Scenario, path string) {
	if scenario.ID == "" {
		validator.addError(path+".id", "required")
	}
	validator.checkOnFailure(scenario.OnFailure, path+".on_failure")
	if scenario.Run != nil && isGroup(scenario) {
		validator.addError(path+".run", "groups cannot have run blocks")
	}
	if isLeaf(scenario) && scenario.BeforeEach != nil {
		validator.addError(path+".before_each", "leaf scenarios cannot have before_each hooks")
	}
	if isLeaf(scenario) && scenario.AfterEach != nil {
		validator.addError(path+".after_each", "leaf scenarios cannot have after_each hooks")
	}
	validator.validateHook(scenario.Before, path+".before")
	validator.validateHook(scenario.After, path+".after")
	validator.validateRunBlock(scenario.Run, path+".run")
	for i, assertion := range scenario.Assertions {
		validator.validateAssertion(assertion, fmt.Sprintf("%s.assertions[%d]", path, i))
	}
	validator.validateScenarios(scenario.Scenarios, path+".scenarios")
}

func (validator *validator) validateScenarios(scenarios []Scenario, basePath string) {
	seenIDs := make(map[string]bool)
	for i, scenario := range scenarios {
		path := fmt.Sprintf("%s[%d]", basePath, i)
		if seenIDs[scenario.ID] {
			validator.addError(path+".id", "duplicate")
		}
		seenIDs[scenario.ID] = true
		validator.validateScenario(scenario, path)
	}
}

func Validate(ctx *Context, filePath string) []ValidationError {
	specValidator := &validator{file: filePath, errors: []ValidationError{}}
	specValidator.checkOnFailure(ctx.OnFailure, "on_failure")
	specValidator.validateHook(ctx.Before, "before")
	specValidator.validateHook(ctx.BeforeEach, "before_each")
	specValidator.validateHook(ctx.After, "after")
	specValidator.validateHook(ctx.AfterEach, "after_each")
	specValidator.validateScenarios(ctx.Scenarios, "scenarios")
	return specValidator.errors
}
