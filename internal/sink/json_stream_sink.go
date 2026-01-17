package sink

import (
	"encoding/json"
	"fmt"
	"io"
)

type JsonStreamSink struct {
	writer io.Writer
}

func NewJsonStreamSink(writer io.Writer) Sink {
	return &JsonStreamSink{writer: writer}
}

func (sink *JsonStreamSink) Emit(incoming any) error {
	data, err := json.Marshal(incoming)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(sink.writer, "%s\n", data)
	return err
}
