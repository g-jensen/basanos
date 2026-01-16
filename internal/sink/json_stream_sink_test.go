package sink

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"basanos/internal/event"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonStreamSink_EmitsNDJSON(t *testing.T) {
	buffer := &bytes.Buffer{}
	sink := NewJsonStreamSink(buffer)

	timestamp := time.Date(2026, 1, 15, 14, 30, 22, 0, time.UTC)
	sink.Emit(event.NewRunStartEvent("2026-01-15_143022", timestamp))

	var result map[string]any
	err := json.Unmarshal(buffer.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "run_start", result["event"])
	assert.Equal(t, "2026-01-15_143022", result["run_id"])
}
