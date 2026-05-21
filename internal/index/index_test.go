package index

import (
	"strings"
	"testing"
	"time"
)

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

const sampleLogs = `2024-01-01T10:00:00Z INFO starting server
2024-01-01T10:01:00Z INFO request received
2024-01-01T10:02:00Z WARN slow response
2024-01-01T10:03:00Z ERROR connection refused
not-a-timestamp this line has no ts
2024-01-01T10:04:00Z INFO shutdown
`

func TestBuildIndex(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	idx, err := Build(r, "")
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}
	// 5 parseable lines, 1 unparseable
	if len(idx.Entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(idx.Entries))
	}
}

func TestBuildIndexEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	idx, err := Build(r, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(idx.Entries))
	}
}

func TestFindRange(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	idx, err := Build(r, "")
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	from := mustTime("2024-01-01T10:01:00Z")
	to := mustTime("2024-01-01T10:03:00Z")
	result := idx.FindRange(from, to)

	if len(result) != 3 {
		t.Errorf("expected 3 entries in range, got %d", len(result))
	}
}

func TestFindRangeNoMatch(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	idx, err := Build(r, "")
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	from := mustTime("2025-01-01T00:00:00Z")
	to := mustTime("2025-01-01T01:00:00Z")
	result := idx.FindRange(from, to)

	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestFindRangeOffsets(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	idx, err := Build(r, "")
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	if len(idx.Entries) == 0 {
		t.Fatal("no entries built")
	}
	// First entry should start at offset 0
	if idx.Entries[0].Offset != 0 {
		t.Errorf("expected first offset 0, got %d", idx.Entries[0].Offset)
	}
}
