package linecount_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/linecount"
)

func TestCountEmpty(t *testing.T) {
	c := linecount.New()
	n, err := c.Count(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
	if c.Bytes() != 0 {
		t.Errorf("expected 0 bytes, got %d", c.Bytes())
	}
}

func TestCountLines(t *testing.T) {
	input := "line one\nline two\nline three\n"
	c := linecount.New()
	n, err := c.Count(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 lines, got %d", n)
	}
	if c.Lines() != 3 {
		t.Errorf("Lines() = %d, want 3", c.Lines())
	}
}

func TestOffsetOf(t *testing.T) {
	input := "abc\ndefg\nhi\n"
	c := linecount.New()
	_, err := c.Count(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		lineNum  int64
		wantOff  int64
	}{
		{1, 0},  // "abc\n" starts at 0
		{2, 4},  // "defg\n" starts at 4
		{3, 9},  // "hi\n" starts at 9
	}
	for _, tt := range tests {
		got := c.OffsetOf(tt.lineNum)
		if got != tt.wantOff {
			t.Errorf("OffsetOf(%d) = %d, want %d", tt.lineNum, got, tt.wantOff)
		}
	}
}

func TestOffsetOfOutOfRange(t *testing.T) {
	c := linecount.New()
	_, _ = c.Count(strings.NewReader("hello\n"))

	if got := c.OffsetOf(0); got != -1 {
		t.Errorf("OffsetOf(0) = %d, want -1", got)
	}
	if got := c.OffsetOf(99); got != -1 {
		t.Errorf("OffsetOf(99) = %d, want -1", got)
	}
}

func TestReset(t *testing.T) {
	c := linecount.New()
	_, _ = c.Count(strings.NewReader("a\nb\nc\n"))
	c.Reset()
	if c.Lines() != 0 {
		t.Errorf("after Reset, Lines() = %d, want 0", c.Lines())
	}
	if c.Bytes() != 0 {
		t.Errorf("after Reset, Bytes() = %d, want 0", c.Bytes())
	}
	if got := c.OffsetOf(1); got != -1 {
		t.Errorf("after Reset, OffsetOf(1) = %d, want -1", got)
	}
}
