package sinkio

type FileSinkWriter interface {
	SetCurrentPath(path string)
	SetCurrentPhase(phase string)
	WriteExitCode(path, phase string, code int) error
	AppendOutput(stream, data string) error
	EnsureOutput(stream string) error
}
