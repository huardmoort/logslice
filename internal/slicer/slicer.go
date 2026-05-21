package slicer

import (
	"fmt"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/reader"
)

// Options controls the slicing behaviour.
type Options struct {
	From            time.Time
	To              time.Time
	TimestampFormat string // empty → auto-detect via ParseTimestamp
	IncludeUnparsed bool   // if true, lines without a timestamp are forwarded
}

// Stats holds counters collected during a slice run.
type Stats struct {
	LinesRead     int64
	LinesMatched  int64
	LinesSkipped  int64
	LinesUnparsed int64
}

// Slice reads lines from lr, writes lines whose timestamp falls within
// [opts.From, opts.To] to w, and returns collected statistics.
func Slice(lr *reader.LineReader, w io.Writer, opts Options) (Stats, error) {
	var stats Stats

	for lr.Next() {
		line := lr.Line()
		stats.LinesRead++

		var (
			ts  time.Time
			err error
		)
		if opts.TimestampFormat != "" {
			ts, err = parser.ParseTimestampWithFormat(line, opts.TimestampFormat)
		} else {
			ts, err = parser.ParseTimestamp(line)
		}

		if err != nil {
			stats.LinesUnparsed++
			if opts.IncludeUnparsed {
				if _, werr := fmt.Fprintln(w, line); werr != nil {
					return stats, werr
				}
			}
			continue
		}

		if parser.InRange(ts, opts.From, opts.To) {
			stats.LinesMatched++
			if _, werr := fmt.Fprintln(w, line); werr != nil {
				return stats, werr
			}
		} else {
			stats.LinesSkipped++
		}
	}

	return stats, lr.Err()
}
