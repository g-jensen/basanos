package event

import (
	"fmt"
	"strconv"

	"basanos/internal/sinkio"
)

func (e *ScenarioRunStartEvent) WriteToFileSink(w sinkio.FileSinkWriter) error {
	w.SetCurrentPath(e.Path)
	w.SetCurrentPhase("_run")
	return nil
}

func (e *HookStartEvent) WriteToFileSink(w sinkio.FileSinkWriter) error {
	w.SetCurrentPath(e.Path)
	w.SetCurrentPhase(e.Hook)
	return nil
}

func (e *HookEndEvent) WriteToFileSink(w sinkio.FileSinkWriter) error {
	return w.WriteExitCode(e.Path, e.Hook, e.ExitCode)
}

func (e *OutputEvent) WriteToFileSink(w sinkio.FileSinkWriter) error {
	return w.AppendOutput(e.Stream, e.Data)
}

func (e *AssertionStartEvent) WriteToFileSink(w sinkio.FileSinkWriter) error {
	w.SetCurrentPath(e.Path)
	w.SetCurrentPhase(fmt.Sprintf("_assertions/%d", e.Index))
	return nil
}

func (e *AssertionEndEvent) WriteToFileSink(w sinkio.FileSinkWriter) error {
	return w.WriteExitCode(e.Path, "_assertions/"+strconv.Itoa(e.Index), e.ExitCode)
}

func (e *ScenarioRunEndEvent) WriteToFileSink(w sinkio.FileSinkWriter) error {
	w.EnsureOutput("stdout")
	w.EnsureOutput("stderr")
	return w.WriteExitCode(e.Path, "_run", e.ExitCode)
}
