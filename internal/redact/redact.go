// Package redact provides pattern-based log line redaction,
// replacing sensitive values (e.g. tokens, passwords, IPs) with
// a configurable placeholder before output.
package redact

import (
	"fmt"
	"regexp"
)

const defaultPlaceholder = "[REDACTED]"

// Redactor applies a set of compiled regex patterns to log lines,
// replacing any matched substrings with a placeholder string.
type Redactor struct {
	patterns    []*regexp.Regexp
	placeholder string
}

// New creates a Redactor from the given regex pattern strings and
// placeholder. If placeholder is empty, "[REDACTED]" is used.
// Returns an error if any pattern fails to compile.
func New(patterns []string, placeholder string) (*Redactor, error) {
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("redact: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Redactor{patterns: compiled, placeholder: placeholder}, nil
}

// Apply returns a copy of line with all pattern matches replaced by
// the placeholder. If the Redactor has no patterns, line is returned
// unchanged.
func (r *Redactor) Apply(line string) string {
	if len(r.patterns) == 0 {
		return line
	}
	for _, re := range r.patterns {
		line = re.ReplaceAllString(line, r.placeholder)
	}
	return line
}

// IsNoop reports whether the Redactor has no patterns and will never
// modify a line.
func (r *Redactor) IsNoop() bool {
	return len(r.patterns) == 0
}
