package stats

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("expected non-nil Stats")
	}
	if s.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}
	if s.TotalLines != 0 || s.MatchedLines != 0 || s.SkippedLines != 0 || s.FilteredLines != 0 {
		t.Error("expected all counters to be zero")
	}
}

func TestRecordCounters(t *testing.T) {
	s := New()
	s.RecordTotal()
	s.RecordTotal()
	s.RecordTotal()
	s.RecordMatched()
	s.RecordMatched()
	s.RecordSkipped()
	s.RecordFiltered()

	if s.TotalLines != 3 {
		t.Errorf("expected TotalLines=3, got %d", s.TotalLines)
	}
	if s.MatchedLines != 2 {
		t.Errorf("expected MatchedLines=2, got %d", s.MatchedLines)
	}
	if s.SkippedLines != 1 {
		t.Errorf("expected SkippedLines=1, got %d", s.SkippedLines)
	}
	if s.FilteredLines != 1 {
		t.Errorf("expected FilteredLines=1, got %d", s.FilteredLines)
	}
}

func TestMatchRate(t *testing.T) {
	s := New()
	if s.MatchRate() != 0 {
		t.Error("expected MatchRate=0 when no lines")
	}
	s.RecordTotal()
	s.RecordTotal()
	s.RecordMatched()
	rate := s.MatchRate()
	if rate != 0.5 {
		t.Errorf("expected MatchRate=0.5, got %f", rate)
	}
}

func TestDuration(t *testing.T) {
	s := New()
	time.Sleep(10 * time.Millisecond)
	d := s.Duration()
	if d < 10*time.Millisecond {
		t.Errorf("expected duration >= 10ms, got %v", d)
	}

	s.Finish()
	finished := s.Duration()
	time.Sleep(10 * time.Millisecond)
	// Duration should not grow after Finish
	if s.Duration() != finished {
		t.Error("expected duration to be fixed after Finish")
	}
}
