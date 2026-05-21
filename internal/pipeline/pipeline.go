// Package pipeline wires together reader, slicer, filter, and formatter
// into a single processing pass over a log file.
package pipeline

import (
	"io"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/slicer"
	"github.com/yourorg/logslice/internal/stats"
)

// Result holds the outcome of a pipeline run.
type Result struct {
	Stats *stats.Stats
}

// Run executes the full log-slicing pipeline using the provided config,
// reading from r and writing output to w.
func Run(cfg *config.Config, r io.Reader, w io.Writer) (*Result, error) {
	st := stats.New()
	st.Start()
	defer st.Stop()

	lr := reader.NewLineReaderFromReader(r)

	fmt, err := formatter.ParseFormat(cfg.Format)
	if err != nil {
		return nil, err
	}
	fw := formatter.NewWriter(w, fmt)

	f, err := filter.New(cfg.Pattern, cfg.Exclude)
	if err != nil {
		return nil, err
	}

	lines, sliceErr := slicer.Slice(lr, cfg.From, cfg.To, cfg.TimestampFormat)
	if sliceErr != nil {
		return nil, sliceErr
	}

	for _, line := range lines {
		st.RecordRead()
		if !f.IsEmpty() && !f.Match(line.Text) {
			st.RecordSkipped()
			continue
		}
		st.RecordMatched()
		if err := fw.WriteLine(line); err != nil {
			return nil, err
		}
	}

	return &Result{Stats: st}, nil
}
