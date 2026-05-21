package formatter

import (
	"fmt"
	"io"
	"strings"
)

// Format represents the output format for sliced log lines.
type Format int

const (
	// FormatRaw outputs lines as-is.
	FormatRaw Format = iota
	// FormatNumbered prefixes each line with its original line number.
	FormatNumbered
	// FormatJSON wraps each line in a simple JSON envelope.
	FormatJSON
)

// ParseFormat parses a format string into a Format value.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "raw", "":
		return FormatRaw, nil
	case "numbered":
		return FormatNumbered, nil
	case "json":
		return FormatJSON, nil
	default:
		return FormatRaw, fmt.Errorf("unknown format %q: must be one of raw, numbered, json", s)
	}
}

// Writer writes log lines to an io.Writer using the configured format.
type Writer struct {
	w      io.Writer
	format Format
}

// NewWriter creates a new formatter Writer.
func NewWriter(w io.Writer, format Format) *Writer {
	return &Writer{w: w, format: format}
}

// WriteLine writes a single log line, applying the configured format.
// lineNum is the 1-based original line number from the source file.
func (fw *Writer) WriteLine(lineNum int, line string) error {
	var out string
	switch fw.format {
	case FormatNumbered:
		out = fmt.Sprintf("%d\t%s\n", lineNum, line)
	case FormatJSON:
		escaped := strings.ReplaceAll(line, `"`, `\"`)
		out = fmt.Sprintf(`{"line":%d,"text":"%s"}`, lineNum, escaped) + "\n"
	default:
		out = line + "\n"
	}
	_, err := fmt.Fprint(fw.w, out)
	return err
}

// WriteLines writes multiple lines using WriteLine.
func (fw *Writer) WriteLines(lines []string, startLineNum int) error {
	for i, line := range lines {
		if err := fw.WriteLine(startLineNum+i, line); err != nil {
			return err
		}
	}
	return nil
}
