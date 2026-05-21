package slicer

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/reader"
)

const sampleLogs = `2024-01-15T10:00:00Z INFO  server started
2024-01-15T10:05:00Z DEBUG request received
2024-01-15T10:10:00Z INFO  processing
2024-01-15T10:15:00Z WARN  slow query
2024-01-15T10:20:00Z INFO  done
`

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSliceBasicRange(t *testing.T) {
	lr := reader.NewLineReaderFromReader(strings.NewReader(sampleLogs))
	var out strings.Builder

	stats, err := Slice(lr, &out, Options{
		From: mustTime("2024-01-15T10:05:00Z"),
		To:   mustTime("2024-01-15T10:15:00Z"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.LinesMatched != 3 {
		t.Errorf("expected 3 matched lines, got %d", stats.LinesMatched)
	}
	if stats.LinesRead != 5 {
		t.Errorf("expected 5 lines read, got %d", stats.LinesRead)
	}
}

func TestSliceNoMatch(t *testing.T) {
	lr := reader.NewLineReaderFromReader(strings.NewReader(sampleLogs))
	var out strings.Builder

	stats, err := Slice(lr, &out, Options{
		From: mustTime("2024-01-16T00:00:00Z"),
		To:   mustTime("2024-01-16T23:59:59Z"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.LinesMatched != 0 {
		t.Errorf("expected 0 matched lines, got %d", stats.LinesMatched)
	}
	if stats.LinesSkipped != 5 {
		t.Errorf("expected 5 skipped lines, got %d", stats.LinesSkipped)
	}
}

func TestSliceUnparsedLines(t *testing.T) {
	input := "not-a-timestamp some log message\n2024-01-15T10:05:00Z INFO ok\n"
	lr := reader.NewLineReaderFromReader(strings.NewReader(input))
	var out strings.Builder

	stats, err := Slice(lr, &out, Options{
		From:            mustTime("2024-01-15T10:00:00Z"),
		To:              mustTime("2024-01-15T11:00:00Z"),
		IncludeUnparsed: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.LinesUnparsed != 1 {
		t.Errorf("expected 1 unparsed line, got %d", stats.LinesUnparsed)
	}
	if stats.LinesMatched != 1 {
		t.Errorf("expected 1 matched line, got %d", stats.LinesMatched)
	}
	if !strings.Contains(out.String(), "not-a-timestamp") {
		t.Error("expected unparsed line to be forwarded to output")
	}
}
