// Package index provides byte-offset indexing for fast seeking within log files.
package index

import (
	"bufio"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Entry represents a single indexed line with its timestamp and byte offset.
type Entry struct {
	Offset    int64
	Timestamp time.Time
	Line      int
}

// Index holds all entries for a log file.
type Index struct {
	Entries []Entry
}

// Build reads from r and constructs an index of timestamp entries.
// Lines that cannot be parsed are skipped.
func Build(r io.ReadSeeker, format string) (*Index, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	idx := &Index{}
	scanner := bufio.NewScanner(r)
	var offset int64
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		ts, err := parser.ParseTimestampWithFormat(line, format)
		if err == nil {
			idx.Entries = append(idx.Entries, Entry{
				Offset:    offset,
				Timestamp: ts,
				Line:      lineNum,
			})
		}
		offset += int64(len(scanner.Bytes())) + 1
	}

	return idx, scanner.Err()
}

// FindRange returns the byte offsets [start, end) that cover entries
// whose timestamps fall within [from, to].
func (idx *Index) FindRange(from, to time.Time) (startOffset, endOffset int64, found bool) {
	if len(idx.Entries) == 0 {
		return 0, 0, false
	}

	start := -1
	end := len(idx.Entries) - 1

	for i, e := range idx.Entries {
		if !e.Timestamp.Before(from) && start == -1 {
			start = i
		}
		if !e.Timestamp.After(to) {
			end = i
		}
	}

	if start == -1 || start > end {
		return 0, 0, false
	}

	return idx.Entries[start].Offset, idx.Entries[end].Offset, true
}
