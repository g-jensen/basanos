package assert

import "strings"

type ContainsResult struct {
	Passed   bool
	Needle   string
	Haystack string
}

func Contains(needle, haystack string) *ContainsResult {
	return &ContainsResult{
		Passed:   strings.Contains(haystack, needle),
		Needle:   needle,
		Haystack: haystack,
	}
}

func (result *ContainsResult) Format() string {
	var output strings.Builder

	if result.Passed {
		output.WriteString("PASS: needle found in haystack\n")
	} else {
		output.WriteString("FAIL: needle not found in haystack\n")
		output.WriteString("──────────────────────────────────\n")
		output.WriteString("Looking for:\n")
		output.WriteString("  " + strings.ReplaceAll(result.Needle, "\n", "\n  ") + "\n")
		output.WriteString("\nIn:\n")
		output.WriteString("  " + strings.ReplaceAll(result.Haystack, "\n", "\n  ") + "\n")
	}

	return output.String()
}
