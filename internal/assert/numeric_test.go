package assert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreaterThan_Pass(t *testing.T) {
	result := GreaterThan("10", "5")

	assert.True(t, result.Passed)
}

func TestGreaterThan_Fail(t *testing.T) {
	result := GreaterThan("5", "10")

	assert.False(t, result.Passed)
}

func TestGreaterThan_Equal_Fail(t *testing.T) {
	result := GreaterThan("5", "5")

	assert.False(t, result.Passed)
}

func TestGreaterThanOrEqual_Equal_Pass(t *testing.T) {
	result := GreaterThanOrEqual("5", "5")

	assert.True(t, result.Passed)
}

func TestLessThan_Pass(t *testing.T) {
	result := LessThan("5", "10")

	assert.True(t, result.Passed)
}

func TestLessThanOrEqual_Pass(t *testing.T) {
	result := LessThanOrEqual("5", "5")

	assert.True(t, result.Passed)
}

func TestNumeric_InvalidNumber_Error(t *testing.T) {
	result := GreaterThan("abc", "5")

	assert.NotEmpty(t, result.Error)
}

func TestNumericResult_Format(t *testing.T) {
	result := NumericResult{
		Passed:   false,
		Left:     "5",
		Right:    "10",
		LeftVal:  5,
		RightVal: 10,
	}

	output := result.Format()

	assert.Contains(t, output, "FAIL")
	assert.Contains(t, output, "5")
	assert.Contains(t, output, "10")
}
