// Package highlight provides ANSI terminal color highlighting for matched
// substrings within log lines.
package highlight

import (
	"regexp"
	"strings"
)

// ANSI escape codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Highlighter applies color to regexp matches within a string.
type Highlighter struct {
	pattern *regexp.Regexp
	color   string
}

// New creates a Highlighter that wraps matches of pattern with the given ANSI
// color code. If pattern is empty, New returns a no-op highlighter.
func New(pattern, color string) (*Highlighter, error) {
	if pattern == "" {
		return &Highlighter{color: color}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Highlighter{pattern: re, color: color}, nil
}

// Apply returns s with every match of the pattern wrapped in the configured
// ANSI color sequence. If no pattern was set, s is returned unchanged.
func (h *Highlighter) Apply(s string) string {
	if h.pattern == nil {
		return s
	}
	var sb strings.Builder
	last := 0
	for _, loc := range h.pattern.FindAllStringIndex(s, -1) {
		sb.WriteString(s[last:loc[0]])
		sb.WriteString(h.color)
		sb.WriteString(s[loc[0]:loc[1]])
		sb.WriteString(Reset)
		last = loc[1]
	}
	sb.WriteString(s[last:])
	return sb.String()
}

// IsNoop reports whether the highlighter will leave input unchanged.
func (h *Highlighter) IsNoop() bool {
	return h.pattern == nil
}
