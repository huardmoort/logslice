package reader

import (
	"strings"
	"testing"
)

func TestLineReaderBasic(t *testing.T) {
	input := "line one\nline two\nline three\n"
	lr := NewLineReaderFromReader(strings.NewReader(input))

	var lines []string
	for lr.Next() {
		lines = append(lines, lr.Line())
	}
	if err := lr.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line one" {
		t.Errorf("expected 'line one', got %q", lines[0])
	}
	if lines[2] != "line three" {
		t.Errorf("expected 'line three', got %q", lines[2])
	}
}

func TestLineReaderLineNumber(t *testing.T) {
	input := "a\nb\nc"
	lr := NewLineReaderFromReader(strings.NewReader(input))

	var num int64
	for lr.Next() {
		num = lr.LineNumber()
	}
	if num != 3 {
		t.Errorf("expected final line number 3, got %d", num)
	}
}

func TestLineReaderEmpty(t *testing.T) {
	lr := NewLineReaderFromReader(strings.NewReader(""))
	if lr.Next() {
		t.Error("expected no lines from empty reader")
	}
	if err := lr.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLineReaderSingleLine(t *testing.T) {
	lr := NewLineReaderFromReader(strings.NewReader("only line"))
	if !lr.Next() {
		t.Fatal("expected one line")
	}
	if lr.Line() != "only line" {
		t.Errorf("got %q", lr.Line())
	}
	if lr.Next() {
		t.Error("expected no more lines")
	}
}
