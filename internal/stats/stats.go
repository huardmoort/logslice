package stats

import "time"

// Stats holds processing statistics for a log slice operation.
type Stats struct {
	TotalLines    int
	MatchedLines  int
	SkippedLines  int
	FilteredLines int
	StartTime     time.Time
	EndTime       time.Time
}

// New creates a new Stats instance with the start time set to now.
func New() *Stats {
	return &Stats{
		StartTime: time.Now(),
	}
}

// Finish marks the end time of the operation.
func (s *Stats) Finish() {
	s.EndTime = time.Now()
}

// Duration returns the elapsed time between Start and Finish.
func (s *Stats) Duration() time.Duration {
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// RecordTotal increments the total line count.
func (s *Stats) RecordTotal() {
	s.TotalLines++
}

// RecordMatched increments the matched line count.
func (s *Stats) RecordMatched() {
	s.MatchedLines++
}

// RecordSkipped increments the skipped line count.
func (s *Stats) RecordSkipped() {
	s.SkippedLines++
}

// RecordFiltered increments the filtered (excluded by filter) line count.
func (s *Stats) RecordFiltered() {
	s.FilteredLines++
}

// MatchRate returns the ratio of matched to total lines, or 0 if no lines processed.
func (s *Stats) MatchRate() float64 {
	if s.TotalLines == 0 {
		return 0
	}
	return float64(s.MatchedLines) / float64(s.TotalLines)
}
