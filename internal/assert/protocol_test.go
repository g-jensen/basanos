package assert

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseProtocol_ValidBasicInput(t *testing.T) {
	input := "basanos:1\n5\nhello5\nworld"

	expected, actual, err := ParseProtocol(strings.NewReader(input))

	assert.NoError(t, err)
	assert.Equal(t, "hello", expected)
	assert.Equal(t, "world", actual)
}

func TestParseProtocol_EmptyExpected(t *testing.T) {
	input := "basanos:1\n0\n5\nworld"

	expected, actual, err := ParseProtocol(strings.NewReader(input))

	assert.NoError(t, err)
	assert.Equal(t, "", expected)
	assert.Equal(t, "world", actual)
}

func TestParseProtocol_EmptyActual(t *testing.T) {
	input := "basanos:1\n5\nhello0\n"

	expected, actual, err := ParseProtocol(strings.NewReader(input))

	assert.NoError(t, err)
	assert.Equal(t, "hello", expected)
	assert.Equal(t, "", actual)
}

func TestParseProtocol_ContentWithNewlines(t *testing.T) {
	input := "basanos:1\n11\nhello\nworld7\nfoo\nbar"

	expected, actual, err := ParseProtocol(strings.NewReader(input))

	assert.NoError(t, err)
	assert.Equal(t, "hello\nworld", expected)
	assert.Equal(t, "foo\nbar", actual)
}

func TestParseProtocol_InvalidVersionHeader(t *testing.T) {
	input := "invalid:2\n5\nhello5\nworld"

	_, _, err := ParseProtocol(strings.NewReader(input))

	assert.ErrorContains(t, err, "version")
}

func TestParseProtocol_TruncatedInput(t *testing.T) {
	input := "basanos:1\n10\nhello"

	_, _, err := ParseProtocol(strings.NewReader(input))

	assert.Error(t, err)
}

func TestBuildProtocol_BasicValues(t *testing.T) {
	result := BuildProtocol("hello", "world")

	assert.Equal(t, "basanos:1\n5\nhello5\nworld", result)
}

func TestBuildProtocol_EmptyExpected(t *testing.T) {
	result := BuildProtocol("", "world")

	assert.Equal(t, "basanos:1\n0\n5\nworld", result)
}

func TestBuildProtocol_ContentWithNewlines(t *testing.T) {
	result := BuildProtocol("hello\nworld", "foo\nbar")

	assert.Equal(t, "basanos:1\n11\nhello\nworld7\nfoo\nbar", result)
}

func TestBuildProtocol_RoundTrip(t *testing.T) {
	expected := "test\nvalue"
	actual := "other\nvalue"

	protocol := BuildProtocol(expected, actual)
	parsedExpected, parsedActual, err := ParseProtocol(strings.NewReader(protocol))

	require.NoError(t, err)
	assert.Equal(t, expected, parsedExpected)
	assert.Equal(t, actual, parsedActual)
}
