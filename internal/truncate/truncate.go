// Package truncate provides utilities for truncating long log lines
// to a configurable maximum byte length, preserving valid UTF-8 boundaries.
package truncate

import (
	"unicode/utf8"
)

// Truncator truncates lines that exceed a maximum byte length.
type Truncator struct {
	maxBytes int
	suffix   string
}

// New creates a new Truncator. maxBytes is the maximum allowed byte length
// of a line (excluding the suffix). If maxBytes <= 0, no truncation is applied.
// suffix is appended to truncated lines (e.g. "...").
func New(maxBytes int, suffix string) *Truncator {
	return &Truncator{
		maxBytes: maxBytes,
		suffix:   suffix,
	}
}

// Apply truncates line if it exceeds the configured maximum byte length.
// It ensures the result is valid UTF-8 by trimming at a rune boundary.
// If maxBytes <= 0, the original line is returned unchanged.
func (t *Truncator) Apply(line string) string {
	if t.maxBytes <= 0 || len(line) <= t.maxBytes {
		return line
	}

	cut := t.maxBytes
	// Walk back to a valid UTF-8 rune boundary.
	for cut > 0 && !utf8.RuneStart(line[cut]) {
		cut--
	}

	return line[:cut] + t.suffix
}

// IsNoop returns true when no truncation will ever be applied.
func (t *Truncator) IsNoop() bool {
	return t.maxBytes <= 0
}
