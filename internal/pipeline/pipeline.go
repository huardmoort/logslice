// Package pipeline wires together all logslice components into a
// single end-to-end processing run.
package pipeline

import (
	"io"

	"github.com/example/logslice/internal/filter"
	"github.com/example/logslice/internal/formatter"
	"github.com/example/logslice/internal/highlight"
	"github.com/example/logslice/internal/reader"
	"github.com/example/logslice/internal/sampler"
	"github.com/example/logslice/internal/slicer"
	"github.com/example/logslice/internal/stats"
)

// Config carries all parameters for a pipeline run.
type Config struct {
	Input       io.Reader
	Output      io.Writer
	From        string
	To          string
	Format      string
	Pattern     string
	Exclude     string
	Highlight   string
	SampleMode  string // "nth" or "random"; empty disables sampling
	SampleRate  int    // ignored when SampleMode is ""
	SampleSeed  int64
}

// Run executes the full pipeline and returns collected statistics.
func Run(cfg Config) (*stats.Stats, error) {
	lr := reader.NewLineReaderFromReader(cfg.Input)

	st := stats.New()

	var smp *sampler.Sampler
	if cfg.SampleMode != "" {
		rate := cfg.SampleRate
		if rate < 1 {
			rate = 1
		}
		var err error
		smp, err = sampler.New(cfg.SampleMode, rate, cfg.SampleSeed)
		if err != nil {
			return nil, err
		}
	}

	f, err := filter.New(cfg.Pattern, cfg.Exclude)
	if err != nil {
		return nil, err
	}

	h, err := highlight.New(cfg.Highlight)
	if err != nil {
		return nil, err
	}

	fmt, err := formatter.ParseFormat(cfg.Format)
	if err != nil {
		return nil, err
	}
	w := formatter.NewWriter(cfg.Output, fmt)

	lines, err := slicer.Slice(lr, cfg.From, cfg.To)
	if err != nil {
		return nil, err
	}

	for _, l := range lines {
		st.RecordRead()
		if smp != nil && !smp.Keep() {
			continue
		}
		if !f.Match(l.Text) {
			continue
		}
		st.RecordMatched()
		l.Text = h.Apply(l.Text)
		if err := w.WriteLine(l); err != nil {
			return st, err
		}
	}
	return st, nil
}
