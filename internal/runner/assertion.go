package runner

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CapturedOutput struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func resolveAssertionArgs(command string, captured CapturedOutput, env map[string]string) (first, second string, err error) {
	expanded := os.Expand(command, func(key string) string {
		return env[key]
	})

	parts := strings.Fields(expanded)
	if len(parts) < 3 {
		return "", "", fmt.Errorf("assertion command must have executable and 2 args")
	}

	first = resolveArg(parts[1], captured, env)
	second = resolveArg(parts[2], captured, env)

	return first, second, nil
}

func resolveArg(arg string, captured CapturedOutput, env map[string]string) string {
	runOutput := env["RUN_OUTPUT"]
	if runOutput == "" {
		return arg
	}

	switch arg {
	case runOutput + "/stdout":
		return captured.Stdout
	case runOutput + "/stderr":
		return captured.Stderr
	case runOutput + "/exit_code":
		return strconv.Itoa(captured.ExitCode)
	default:
		return arg
	}
}
