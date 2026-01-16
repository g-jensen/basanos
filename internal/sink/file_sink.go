package sink

import (
	"path/filepath"
	"strconv"

	"basanos/internal/sinkio"
)

type WritableFS interface {
	WriteFile(path string, data []byte) error
	AppendFile(path string, data []byte) error
	ReadFile(path string) ([]byte, error)
}

type FileSinkHandler interface {
	WriteToFileSink(w sinkio.FileSinkWriter) error
}

type FileSink struct {
	fs           WritableFS
	runID        string
	currentPath  string
	currentPhase string
}

func NewFileSink(fs WritableFS, runID string) *FileSink {
	return &FileSink{
		fs:    fs,
		runID: runID,
	}
}

func (s *FileSink) Emit(e any) error {
	if handler, ok := e.(FileSinkHandler); ok {
		return handler.WriteToFileSink(s)
	}
	return nil
}

func (s *FileSink) SetCurrentPath(path string) {
	s.currentPath = path
}

func (s *FileSink) SetCurrentPhase(phase string) {
	s.currentPhase = phase
}

func (s *FileSink) WriteExitCode(path, phase string, code int) error {
	filePath := filepath.Join(s.runID, path, phase, "exit_code")
	return s.fs.WriteFile(filePath, []byte(strconv.Itoa(code)))
}

func (s *FileSink) AppendOutput(stream, data string) error {
	path := filepath.Join(s.runID, s.currentPath, s.currentPhase, stream)
	return s.fs.AppendFile(path, []byte(data))
}

func (s *FileSink) EnsureOutput(stream string) error {
	path := filepath.Join(s.runID, s.currentPath, s.currentPhase, stream)
	_, err := s.fs.ReadFile(path)
	if err != nil {
		return s.fs.WriteFile(path, []byte{})
	}
	return nil
}
