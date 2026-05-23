package truncate

import (
	"strings"
	"testing"
)

func TestNewNoopWhenZero(t *testing.T) {
	tr := New(0, "...")
	if !tr.IsNoop() {
		t.Fatal("expected IsNoop() == true for maxBytes=0")
	}
}

func TestNewNoopWhenNegative(t *testing.T) {
	tr := New(-1, "...")
	if !tr.IsNoop() {
		t.Fatal("expected IsNoop() == true for maxBytes=-1")
	}
}

func TestApplyNoTruncationNeeded(t *testing.T) {
	tr := New(100, "...")
	line := "short line"
	got := tr.Apply(line)
	if got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApplyExactLength(t *testing.T) {
	tr := New(10, "...")
	line := "1234567890"
	got := tr.Apply(line)
	if got != line {
		t.Fatalf("expected %q unchanged, got %q", line, got)
	}
}

func TestApplyTruncatesWithSuffix(t *testing.T) {
	tr := New(5, "...")
	line := "hello world"
	got := tr.Apply(line)
	if got != "hello..." {
		t.Fatalf("expected %q, got %q", "hello...", got)
	}
}

func TestApplyNoopWhenMaxBytesZero(t *testing.T) {
	tr := New(0, "...")
	line := strings.Repeat("x", 200)
	got := tr.Apply(line)
	if got != line {
		t.Fatalf("expected unchanged line, got truncated")
	}
}

func TestApplyUTF8Boundary(t *testing.T) {
	// "é" is 2 bytes (0xC3 0xA9); cutting at byte 1 would split the rune.
	tr := New(3, "")
	line := "aéb" // a=1 byte, é=2 bytes, b=1 byte => total 4 bytes
	got := tr.Apply(line)
	// maxBytes=3: bytes 0-2 = 'a' + first byte of 'é'; must walk back to byte 1.
	if !strings.HasPrefix(got, "a") {
		t.Fatalf("expected result to start with 'a', got %q", got)
	}
	if len(got) > 3 {
		t.Fatalf("result exceeds maxBytes: len=%d, got %q", len(got), got)
	}
}

func TestApplyEmptySuffix(t *testing.T) {
	tr := New(4, "")
	line := "abcdefgh"
	got := tr.Apply(line)
	if got != "abcd" {
		t.Fatalf("expected %q, got %q", "abcd", got)
	}
}
