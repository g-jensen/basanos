package cli

type colorizer interface {
	green(text string) string
	red(text string) string
	formatName(name, status string) string
}

type noopColorizer struct{}

func (noop noopColorizer) green(text string) string { return text }
func (noop noopColorizer) red(text string) string   { return text }
func (noop noopColorizer) formatName(name, status string) string {
	return name + " " + statusChar(status)
}

type ansiColorizer struct{}

func (ansi ansiColorizer) green(text string) string { return "\033[32m" + text + "\033[0m" }
func (ansi ansiColorizer) red(text string) string   { return "\033[31m" + text + "\033[0m" }
func (ansi ansiColorizer) formatName(name, status string) string {
	if status == "pass" {
		return ansi.green(name)
	}
	return ansi.red(name)
}

func statusChar(status string) string {
	if status == "pass" {
		return "."
	}
	if status == "fail" {
		return "F"
	}
	return ""
}

func newColorizer(enabled bool) colorizer {
	if enabled {
		return ansiColorizer{}
	}
	return noopColorizer{}
}
