package index_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/index"
)

const sampleLog = `2024-01-10T10:00:00Z INFO starting server
2024-01-10T10:01:00Z INFO listening on :8080
2024-01-10T10:02:00Z WARN high memory usage
2024-01-10T10:03:00Z ERROR connection refused
2024-01-10T10:04:00Z INFO shutdown complete
`

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestBuildIndex(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, err := index.Build(r, "")
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if len(idx.Entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(idx.Entries))
	}
}

func TestBuildIndexEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	idx, err := index.Build(r, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(idx.Entries))
	}
}

func TestFindRange(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, _ := index.Build(r, "")

	from := mustTime("2024-01-10T10:01:00Z")
	to := mustTime("2024-01-10T10:03:00Z")

	start, end, found := idx.FindRange(from, to)
	if !found {
		t.Fatal("expected range to be found")
	}
	if start >= end {
		t.Errorf("expected start < end, got start=%d end=%d", start, end)
	}
}

func TestFindRangeNoMatch(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, _ := index.Build(r, "")

	from := mustTime("2025-01-01T00:00:00Z")
	to := mustTime("2025-01-01T01:00:00Z")

	_, _, found := idx.FindRange(from, to)
	if found {
		t.Error("expected no range found")
	}
}

func TestFindRangeEmptyIndex(t *testing.T) {
	idx := &index.Index{}
	_, _, found := idx.FindRange(mustTime("2024-01-01T00:00:00Z"), mustTime("2024-01-02T00:00:00Z"))
	if found {
		t.Error("expected no range on empty index")
	}
}
