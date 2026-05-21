// Package index provides offset-based indexing for fast log file seeking.
package index

import (
	"bufio"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Entry represents a single indexed line with its byte offset and parsed time.
type Entry struct {
	Offset    int64
	Timestamp time.Time
	Line      string
}

// Index holds all entries built from a log source.
type Index struct {
	Entries []Entry
}

// Build reads all lines from r and constructs an Index by parsing timestamps.
// Lines that cannot be parsed are skipped.
func Build(r io.ReadSeeker, format string) (*Index, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var entries []Entry
	var offset int64
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		ts, err := parser.ParseTimestampWithFormat(line, format)
		if err == nil {
			entries = append(entries, Entry{
				Offset:    offset,
				Timestamp: ts,
				Line:      line,
			})
		}
		offset += int64(len(line)) + 1 // +1 for newline
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Index{Entries: entries}, nil
}

// FindRange returns the slice of entries whose timestamps fall within [from, to].
func (idx *Index) FindRange(from, to time.Time) []Entry {
	var result []Entry
	for _, e := range idx.Entries {
		if parser.InRange(e.Timestamp, from, to) {
			result = append(result, e)
		}
	}
	return result
}
