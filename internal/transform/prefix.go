package transform

import "fmt"

// PrependField returns a Func that prepends a key=value field to each line.
// This is useful for injecting metadata such as a source filename.
func PrependField(key, value string) Func {
	if key == "" || value == "" {
		return func(line string) (string, bool) { return line, true }
	}
	prefix := fmt.Sprintf("%s=%s ", key, value)
	return func(line string) (string, bool) {
		return prefix + line, true
	}
}

// AppendField returns a Func that appends a key=value field to each line.
func AppendField(key, value string) Func {
	if key == "" || value == "" {
		return func(line string) (string, bool) { return line, true }
	}
	suffix := fmt.Sprintf(" %s=%s", key, value)
	return func(line string) (string, bool) {
		return line + suffix, true
	}
}

// AddLineNumber returns a Func that prepends a sequential line number.
// The counter starts at startAt and increments with each call.
func AddLineNumber(startAt int) Func {
	n := startAt - 1
	return func(line string) (string, bool) {
		n++
		return fmt.Sprintf("%d: %s", n, line), true
	}
}
