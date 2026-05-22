package highlight

import (
	"strings"
	"testing"
)

func TestNewInvalidPattern(t *testing.T) {
	_, err := New("[invalid", Red)
	if err == nil {
		t.Fatal("expected error for invalid regexp, got nil")
	}
}

func TestNewEmptyPatternIsNoop(t *testing.T) {
	h, err := New("", Yellow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !h.IsNoop() {
		t.Error("expected IsNoop() == true for empty pattern")
	}
	input := "hello world"
	if got := h.Apply(input); got != input {
		t.Errorf("Apply() = %q, want %q", got, input)
	}
}

func TestApplySingleMatch(t *testing.T) {
	h, err := New("ERROR", Red)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "2024-01-01 ERROR something failed"
	got := h.Apply(input)
	if !strings.Contains(got, Red+"ERROR"+Reset) {
		t.Errorf("Apply() = %q, expected ANSI-wrapped ERROR", got)
	}
	if !strings.Contains(got, "something failed") {
		t.Errorf("Apply() dropped trailing text: %q", got)
	}
}

func TestApplyMultipleMatches(t *testing.T) {
	h, err := New("foo", Cyan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "foo bar foo baz foo"
	got := h.Apply(input)
	count := strings.Count(got, Cyan+"foo"+Reset)
	if count != 3 {
		t.Errorf("expected 3 highlighted matches, got %d in %q", count, got)
	}
}

func TestApplyNoMatch(t *testing.T) {
	h, err := New("WARN", Yellow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "2024-01-01 INFO everything is fine"
	got := h.Apply(input)
	if got != input {
		t.Errorf("Apply() = %q, want unchanged %q", got, input)
	}
}

func TestIsNoopWithPattern(t *testing.T) {
	h, err := New("something", Green)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.IsNoop() {
		t.Error("expected IsNoop() == false when pattern is set")
	}
}
