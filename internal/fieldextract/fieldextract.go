// Package fieldextract provides utilities for extracting named fields
// from structured log lines (JSON or key=value logfmt style).
package fieldextract

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Format represents the detected or configured log line format.
type Format int

const (
	FormatAuto   Format = iota // detect automatically
	FormatJSON                 // JSON object per line
	FormatLogfmt               // key=value pairs
)

// Extractor extracts a named field value from a log line.
type Extractor struct {
	field  string
	format Format
}

// New creates an Extractor for the given field name and format.
func New(field string, format Format) (*Extractor, error) {
	if strings.TrimSpace(field) == "" {
		return nil, fmt.Errorf("fieldextract: field name must not be empty")
	}
	return &Extractor{field: field, format: format}, nil
}

// Extract returns the value of the configured field from line.
// Returns an empty string and false if the field is not found.
func (e *Extractor) Extract(line string) (string, bool) {
	switch e.format {
	case FormatJSON:
		return extractJSON(line, e.field)
	case FormatLogfmt:
		return extractLogfmt(line, e.field)
	default: // FormatAuto
		if len(line) > 0 && line[0] == '{' {
			return extractJSON(line, e.field)
		}
		return extractLogfmt(line, e.field)
	}
}

func extractJSON(line, field string) (string, bool) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", false
	}
	v, ok := obj[field]
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%v", v), true
}

func extractLogfmt(line, field string) (string, bool) {
	prefix := field + "="
	for _, part := range strings.Fields(line) {
		if strings.HasPrefix(part, prefix) {
			val := strings.TrimPrefix(part, prefix)
			val = strings.Trim(val, `"`)
			return val, true
		}
	}
	return "", false
}
