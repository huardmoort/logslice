// Package transform applies a chain of line transformations in sequence.
package transform

import "strings"

// Func is a function that transforms a single log line.
// It returns the transformed line and whether the line should be kept.
type Func func(line string) (string, bool)

// Chain holds an ordered list of transformation functions.
type Chain struct {
	funcs []Func
}

// New creates a new Chain from the provided transformation functions.
func New(funcs ...Func) *Chain {
	return &Chain{funcs: funcs}
}

// Apply runs all transformations in order on the given line.
// If any transformation drops the line (returns false), Apply returns ("", false).
func (c *Chain) Apply(line string) (string, bool) {
	for _, fn := range c.funcs {
		var keep bool
		line, keep = fn(line)
		if !keep {
			return "", false
		}
	}
	return line, true
}

// IsNoop returns true when the chain contains no transformations.
func (c *Chain) IsNoop() bool {
	return len(c.funcs) == 0
}

// TrimSpace returns a Func that trims leading and trailing whitespace.
func TrimSpace() Func {
	return func(line string) (string, bool) {
		return strings.TrimSpace(line), true
	}
}

// DropEmpty returns a Func that drops blank lines.
func DropEmpty() Func {
	return func(line string) (string, bool) {
		if strings.TrimSpace(line) == "" {
			return "", false
		}
		return line, true
	}
}

// ReplaceAll returns a Func that replaces all occurrences of old with new.
func ReplaceAll(old, replacement string) Func {
	return func(line string) (string, bool) {
		return strings.ReplaceAll(line, old, replacement), true
	}
}

// MaxLength returns a Func that truncates lines exceeding maxLen characters.
func MaxLength(maxLen int) Func {
	return func(line string) (string, bool) {
		if maxLen > 0 && len(line) > maxLen {
			return line[:maxLen], true
		}
		return line, true
	}
}
