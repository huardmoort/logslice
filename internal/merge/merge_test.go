package merge

import (
	"testing"
	"time"
)

func makeSource(entries []Entry) <-chan Entry {
	ch := make(chan Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestMergeEmpty(t *testing.T) {
	m := New(nil)
	var got []Entry
	for e := range m.Merge() {
		got = append(got, e)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestMergeSingleSource(t *testing.T) {
	src := makeSource([]Entry{
		{Line: "a", Timestamp: mustTime("2024-01-01T00:00:01Z")},
		{Line: "b", Timestamp: mustTime("2024-01-01T00:00:02Z")},
	})
	m := New([]<-chan Entry{src})
	var got []Entry
	for e := range m.Merge() {
		got = append(got, e)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Line != "a" || got[1].Line != "b" {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestMergeMultipleSources(t *testing.T) {
	src1 := makeSource([]Entry{
		{Line: "s1-1", Timestamp: mustTime("2024-01-01T00:00:01Z")},
		{Line: "s1-3", Timestamp: mustTime("2024-01-01T00:00:03Z")},
	})
	src2 := makeSource([]Entry{
		{Line: "s2-2", Timestamp: mustTime("2024-01-01T00:00:02Z")},
		{Line: "s2-4", Timestamp: mustTime("2024-01-01T00:00:04Z")},
	})
	m := New([]<-chan Entry{src1, src2})
	var got []Entry
	for e := range m.Merge() {
		got = append(got, e)
	}
	expected := []string{"s1-1", "s2-2", "s1-3", "s2-4"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d entries, got %d", len(expected), len(got))
	}
	for i, e := range got {
		if e.Line != expected[i] {
			t.Errorf("pos %d: want %q, got %q", i, expected[i], e.Line)
		}
	}
}

func TestMergePreservesSourceIndex(t *testing.T) {
	src0 := makeSource([]Entry{{Line: "x", Timestamp: mustTime("2024-01-01T00:00:01Z")}})
	src1 := makeSource([]Entry{{Line: "y", Timestamp: mustTime("2024-01-01T00:00:02Z")}})
	m := New([]<-chan Entry{src0, src1})
	got := make([]Entry, 0)
	for e := range m.Merge() {
		got = append(got, e)
	}
	if got[0].Source != 0 {
		t.Errorf("expected source 0, got %d", got[0].Source)
	}
	if got[1].Source != 1 {
		t.Errorf("expected source 1, got %d", got[1].Source)
	}
}
