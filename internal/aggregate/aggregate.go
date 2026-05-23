// Package aggregate provides line-count and field-value aggregation
// over a stream of log lines, useful for summarising log data without
// emitting every matched line.
package aggregate

import (
	"fmt"
	"sort"
	"strings"
)

// Mode controls how values are aggregated.
type Mode string

const (
	ModeCount  Mode = "count"  // count occurrences of each unique value
	ModeUnique Mode = "unique" // collect distinct values only
)

// Aggregator accumulates values extracted from log lines.
type Aggregator struct {
	mode   Mode
	counts map[string]int
}

// New returns an Aggregator for the given mode.
// Returns an error for unrecognised modes.
func New(mode Mode) (*Aggregator, error) {
	switch mode {
	case ModeCount, ModeUnique:
	default:
		return nil, fmt.Errorf("aggregate: unknown mode %q", mode)
	}
	return &Aggregator{mode: mode, counts: make(map[string]int)}, nil
}

// Add records a value extracted from a log line.
func (a *Aggregator) Add(value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	a.counts[value]++
}

// Results returns aggregated entries sorted by key.
// In ModeCount each entry is "value\tN"; in ModeUnique each entry is the
// distinct value.
func (a *Aggregator) Results() []string {
	keys := make([]string, 0, len(a.counts))
	for k := range a.counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if a.mode == ModeCount {
			out = append(out, fmt.Sprintf("%s\t%d", k, a.counts[k]))
		} else {
			out = append(out, k)
		}
	}
	return out
}

// Reset clears all accumulated data.
func (a *Aggregator) Reset() {
	a.counts = make(map[string]int)
}

// Total returns the sum of all recorded occurrences.
func (a *Aggregator) Total() int {
	n := 0
	for _, v := range a.counts {
		n += v
	}
	return n
}
