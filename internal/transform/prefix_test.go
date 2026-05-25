package transform

import (
	"strings"
	"testing"
)

func TestPrependField(t *testing.T) {
	c := New(PrependField("src", "app.log"))
	out, ok := c.Apply("some log line")
	if !ok {
		t.Fatal("expected line to be kept")
	}
	if !strings.HasPrefix(out, "src=app.log ") {
		t.Fatalf("expected prefix, got: %q", out)
	}
}

func TestPrependFieldEmptyKeyIsNoop(t *testing.T) {
	c := New(PrependField("", "value"))
	out, ok := c.Apply("line")
	if !ok || out != "line" {
		t.Fatalf("expected noop, got %q %v", out, ok)
	}
}

func TestPrependFieldEmptyValueIsNoop(t *testing.T) {
	c := New(PrependField("key", ""))
	out, ok := c.Apply("line")
	if !ok || out != "line" {
		t.Fatalf("expected noop, got %q %v", out, ok)
	}
}

func TestAppendField(t *testing.T) {
	c := New(AppendField("env", "prod"))
	out, ok := c.Apply("msg")
	if !ok {
		t.Fatal("expected line to be kept")
	}
	if !strings.HasSuffix(out, " env=prod") {
		t.Fatalf("expected suffix, got: %q", out)
	}
}

func TestAppendFieldEmptyIsNoop(t *testing.T) {
	c := New(AppendField("", ""))
	out, ok := c.Apply("data")
	if !ok || out != "data" {
		t.Fatalf("expected noop, got %q", out)
	}
}

func TestAddLineNumber(t *testing.T) {
	c := New(AddLineNumber(1))
	lines := []string{"alpha", "beta", "gamma"}
	for i, l := range lines {
		out, ok := c.Apply(l)
		if !ok {
			t.Fatalf("line %d dropped unexpectedly", i)
		}
		expected := strings.HasPrefix(out, "1: ") || strings.HasPrefix(out, "2: ") || strings.HasPrefix(out, "3: ")
		if !expected {
			t.Fatalf("unexpected line number format: %q", out)
		}
	}
}

func TestAddLineNumberStartAt(t *testing.T) {
	c := New(AddLineNumber(10))
	out, _ := c.Apply("first")
	if !strings.HasPrefix(out, "10: ") {
		t.Fatalf("expected start at 10, got %q", out)
	}
}
