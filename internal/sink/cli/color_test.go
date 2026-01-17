package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusChar_Pass(t *testing.T) {
	assert.Equal(t, ".", statusChar("pass"))
}

func TestStatusChar_Fail(t *testing.T) {
	assert.Equal(t, "F", statusChar("fail"))
}

func TestStatusChar_Unknown(t *testing.T) {
	assert.Equal(t, "", statusChar("unknown"))
}

func TestNoopColorizer_Green(t *testing.T) {
	noop := noopColorizer{}
	assert.Equal(t, "hello", noop.green("hello"))
}

func TestNoopColorizer_Red(t *testing.T) {
	noop := noopColorizer{}
	assert.Equal(t, "hello", noop.red("hello"))
}

func TestNoopColorizer_FormatName_Pass(t *testing.T) {
	noop := noopColorizer{}
	assert.Equal(t, "Login works .", noop.formatName("Login works", "pass"))
}

func TestNoopColorizer_FormatName_Fail(t *testing.T) {
	noop := noopColorizer{}
	assert.Equal(t, "Login fails F", noop.formatName("Login fails", "fail"))
}

func TestAnsiColorizer_Green(t *testing.T) {
	ansi := ansiColorizer{}
	assert.Equal(t, "\033[32mhello\033[0m", ansi.green("hello"))
}

func TestAnsiColorizer_Red(t *testing.T) {
	ansi := ansiColorizer{}
	assert.Equal(t, "\033[31mhello\033[0m", ansi.red("hello"))
}

func TestAnsiColorizer_FormatName_Pass(t *testing.T) {
	ansi := ansiColorizer{}
	assert.Equal(t, "\033[32mLogin works\033[0m", ansi.formatName("Login works", "pass"))
}

func TestAnsiColorizer_FormatName_Fail(t *testing.T) {
	ansi := ansiColorizer{}
	assert.Equal(t, "\033[31mLogin fails\033[0m", ansi.formatName("Login fails", "fail"))
}

func TestNewColorizer_Enabled(t *testing.T) {
	colors := newColorizer(true)
	_, ok := colors.(ansiColorizer)
	assert.True(t, ok)
}

func TestNewColorizer_Disabled(t *testing.T) {
	colors := newColorizer(false)
	_, ok := colors.(noopColorizer)
	assert.True(t, ok)
}
