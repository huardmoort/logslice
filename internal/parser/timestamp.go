package parser

import (
	"fmt"
	"time"
)

// CommonFormats lists timestamp formats commonly found in log files.
var CommonFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.000Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.000000",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// ParseTimestamp attempts to parse a timestamp string using a list of known
// formats. It returns the parsed time and the matched format, or an error if
// no format matches.
func ParseTimestamp(s string) (time.Time, string, error) {
	for _, layout := range CommonFormats {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, layout, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("parser: unrecognized timestamp format: %q", s)
}

// ParseTimestampWithFormat parses a timestamp using an explicit format string.
func ParseTimestampWithFormat(s, layout string) (time.Time, error) {
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("parser: cannot parse %q with format %q: %w", s, layout, err)
	}
	return t, nil
}

// InRange reports whether t falls within [start, end] (inclusive).
func InRange(t, start, end time.Time) bool {
	return !t.Before(start) && !t.After(end)
}
