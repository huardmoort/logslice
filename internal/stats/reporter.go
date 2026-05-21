package stats

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Reporter writes a human-readable summary of Stats to an io.Writer.
type Reporter struct {
	w io.Writer
}

// NewReporter creates a Reporter that writes to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{w: w}
}

// Print writes a formatted statistics summary to the reporter's writer.
func (r *Reporter) Print(s *Stats) error {
	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)

	fmt.Fprintln(tw, "--- logslice stats ---")
	fmt.Fprintf(tw, "Duration:\t%v\n", s.Duration().Round(1000000))
	fmt.Fprintf(tw, "Total lines:\t%d\n", s.TotalLines)
	fmt.Fprintf(tw, "Matched lines:\t%d\n", s.MatchedLines)
	fmt.Fprintf(tw, "Skipped lines:\t%d\n", s.SkippedLines)
	fmt.Fprintf(tw, "Filtered lines:\t%d\n", s.FilteredLines)
	fmt.Fprintf(tw, "Match rate:\t%.1f%%\n", s.MatchRate()*100)

	return tw.Flush()
}

// PrintJSON writes a JSON-formatted statistics summary.
func (r *Reporter) PrintJSON(s *Stats) error {
	_, err := fmt.Fprintf(r.w,
		`{"duration_ms":%d,"total":%d,"matched":%d,"skipped":%d,"filtered":%d,"match_rate":%.4f}`+"\n",
		s.Duration().Milliseconds(),
		s.TotalLines,
		s.MatchedLines,
		s.SkippedLines,
		s.FilteredLines,
		s.MatchRate(),
	)
	return err
}
