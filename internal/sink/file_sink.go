package sink

import (
	"path/filepath"
	"strconv"

	"basanos/internal/fs"
	"basanos/internal/sinkio"
)

type FileSinkHandler interface {
	WriteToFileSink(w sinkio.FileSinkWriter) error
}

type FileSink struct {
	fs           fs.WritableFS
	runID        string
	currentPath  string
	currentPhase string
}

func NewFileSink(filesystem fs.WritableFS, runID string) *FileSink {
	return &FileSink{
		fs:    filesystem,
		runID: runID,
	}
}

func (sink *FileSink) Emit(incoming any) error {
	if handler, ok := incoming.(FileSinkHandler); ok {
		return handler.WriteToFileSink(sink)
	}
	return nil
}

func (sink *FileSink) SetCurrentPath(path string) {
	sink.currentPath = path
}

func (sink *FileSink) SetCurrentPhase(phase string) {
	sink.currentPhase = phase
}

func (sink *FileSink) WriteExitCode(path, phase string, code int) error {
	filePath := filepath.Join(sink.runID, path, phase, "exit_code")
	return sink.fs.WriteFile(filePath, []byte(strconv.Itoa(code)))
}

func (sink *FileSink) AppendOutput(stream, data string) error {
	path := filepath.Join(sink.runID, sink.currentPath, sink.currentPhase, stream)
	return sink.fs.AppendFile(path, []byte(data))
}

func (sink *FileSink) EnsureOutput(stream string) error {
	path := filepath.Join(sink.runID, sink.currentPath, sink.currentPhase, stream)
	_, err := sink.fs.ReadFile(path)
	if err != nil {
		return sink.fs.WriteFile(path, []byte{})
	}
	return nil
}
